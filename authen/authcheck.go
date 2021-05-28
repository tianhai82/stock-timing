package authen

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tianhai82/stock-timing/firebase"
	"github.com/tianhai82/stock-timing/model"
)

// AuthCheck is a middleware that extracts and verify the idToken and put the user in gin context
func AuthCheck(c *gin.Context) {
	if firebase.AuthClient == nil {
		return
	}
	idToken := c.GetHeader("Authorization")
	if idToken == "" {
		return
	}
	ctx := context.Background()
	token, err := firebase.AuthClient.VerifyIDTokenAndCheckRevoked(ctx, idToken[len("Bearer "):])
	if err != nil {
		fmt.Println(err)
		return
	}
	email, found := token.Claims["email"]
	if !found {
		return
	}
	emailStr, isString := email.(string)
	if !isString || emailStr == "" {
		return
	}
	emailVerified := token.Claims["email_verified"]
	emailVerifiedBool := emailVerified.(bool)

	user := model.UserAccount{
		Email:         emailStr,
		EmailVerified: emailVerifiedBool,
	}
	c.Set("loginUser", user)
}

func TdaAuth(c *gin.Context) {
	key := os.Getenv("TDA_KEY")
	if key == "" {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	now := time.Now().UTC().Unix()

	idToken := c.GetHeader("Authorization")
	if idToken == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	idToken = idToken[len("Bearer "):]
	segment := strings.Split(idToken, ".")
	if len(segment) != 2 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	payload, err := base64.RawURLEncoding.DecodeString(segment[0])
	if err != nil {
		fmt.Println("fail to decode payload", segment[0], err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	reqTime, err := strconv.ParseInt(string(payload), 10, 64)
	if err != nil {
		fmt.Println("fail to convert payload to int64", segment[0], err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if (now-120 > reqTime) || (now+120 < reqTime) {
		fmt.Println("payload is out of acceptable tme range", reqTime, now)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	signature, err := base64.RawURLEncoding.DecodeString(segment[1])
	if err != nil {
		fmt.Println("fail to decode signature", segment[1], err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if !bytes.Equal(signature, hmac256hash(payload, []byte(key))) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}
func hmac256hash(msg []byte, key []byte) []byte {
	sig := hmac.New(sha256.New, key)
	sig.Write([]byte(msg))
	return []byte(hex.EncodeToString(sig.Sum(nil)))
}
