package rpcs

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tianhai82/stock-timing/model"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func AddSubscriptionRpcs(router *gin.RouterGroup) {
	router.POST("/subscribe/period/:period", subscribe)
	router.GET("/subscriptions", retrieveSubscription)
	router.DELETE("/subscriptions/:instrumentID", deleteSubscription)
}

func deleteSubscription(c *gin.Context) {
	user, exist := c.Get("loginUser")
	if !exist {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	loginUser := user.(model.UserAccount)
	instrumentIDStr := c.Param("instrumentID")
	instrumentID, err := strconv.Atoi(instrumentIDStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid period")
		return
	}
	docID := fmt.Sprintf("%s|%d", loginUser.Email, instrumentID)
	_, err = FirestoreClient.Collection("subscription").Doc(docID).Delete(context.Background())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "fail to remove subscription")
		return
	}
}

func retrieveSubscription(c *gin.Context) {
	user, exist := c.Get("loginUser")
	if !exist {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	loginUser := user.(model.UserAccount)
	subscriptions, err := FirestoreClient.Collection("subscription").Where("UserID", "==", loginUser.Email).Documents(context.Background()).GetAll()
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	stockSubs := make([]model.StockSubscription, 0, len(subscriptions))
	for _, s := range subscriptions {
		var stockSub model.StockSubscription
		err := s.DataTo(&stockSub)
		if err != nil {
			fmt.Println(err)
			continue
		}
		stockSubs = append(stockSubs, stockSub)
	}
	c.JSON(200, stockSubs)
}

func subscribe(c *gin.Context) {
	user, exist := c.Get("loginUser")
	if !exist {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	loginUser := user.(model.UserAccount)

	periodStr := c.Param("period")
	period, err := strconv.Atoi(periodStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid period")
		return
	}

	var stock model.Stock
	err = c.BindJSON(&stock)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	subscription := model.StockSubscription{
		Stock:         stock,
		UserID:        loginUser.Email,
		Period:        period,
		LastUpdatedAt: time.Now(),
	}
	key := strings.ToLower(loginUser.Email) + "|" + strconv.Itoa(stock.InstrumentID)
	_, err = FirestoreClient.Collection("subscription").Doc(key).Set(context.Background(), subscription)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}
