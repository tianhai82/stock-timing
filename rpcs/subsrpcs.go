package rpcs

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tianhai82/stock-timing/model"
)

func AddSubscriptionRpcs(router *gin.RouterGroup) {
	router.POST("/subscribe", subscribe)
}

func subscribe(c *gin.Context) {
	user, exist := c.Get("loginUser")
	if !exist {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	loginUser := user.(model.UserAccount)
	fmt.Println(loginUser.Email)
	var stock model.Stock
	err := c.BindJSON(&stock)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	subscription := model.StockSubscription{
		Stock:         stock,
		UserID:        loginUser.Email,
		LastUpdatedAt: time.Now(),
	}
	key := strings.ToLower(loginUser.Email) + "|" + strconv.Itoa(stock.InstrumentID)
	_, err = FirestoreClient.Collection("subscription").Doc(key).Set(context.Background(), subscription)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}
