package authen

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/tianhai82/stock-timing/model"
	"github.com/tianhai82/stock-timing/rpcs"
)

func AuthCheck(c *gin.Context) {
	if rpcs.AuthClient == nil {
		return
	}
	idToken := c.GetHeader("idToken")
	if idToken == "" {
		return
	}
	ctx := context.Background()
	token, err := rpcs.AuthClient.VerifyIDTokenAndCheckRevoked(ctx, idToken)
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
