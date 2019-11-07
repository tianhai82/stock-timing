package cron

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/tianhai82/stock-timing/cloudtask"
	"github.com/tianhai82/stock-timing/model"
	"github.com/tianhai82/stock-timing/rpcs"
)

// AddCronJobs adds cron job handler to router group
func AddCronJobs(router *gin.RouterGroup) {
	router.GET("/stockSubscriptions", stockSubscriptions)
}

func stockSubscriptions(c *gin.Context) {
	ctx := context.Background()
	iter := rpcs.FirestoreClient.Collection("subscription").Documents(ctx)
	docs, err := iter.GetAll()
	if err != nil {
		fmt.Println("unable to retrieve subscriptions", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ids := distinctInstrumentID(docs)
	instrumentCol := rpcs.FirestoreClient.Collection("instrumentsToAnalyse")
	for _, id := range ids {
		_, errSet := instrumentCol.Doc(strconv.Itoa(id)).Set(ctx, map[string]int{"id": id})
		if errSet != nil {
			fmt.Println(errSet, id)
		}
	}
	userSubscriptions := groupUsers(docs)
	usersToEmailCol := rpcs.FirestoreClient.Collection("usersToEmail")
	for _, subs := range userSubscriptions {
		_, errSet := usersToEmailCol.Doc(subs.UserID).Set(ctx, subs)
		if errSet != nil {
			fmt.Println(errSet, subs)
		}
	}
	_, err = cloudtask.CreateTask("stock-timing", "asia-south1", "analyze-stocks", nil)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	// create cloud task for Analyse Stock
}

func groupUsers(docs []*firestore.DocumentSnapshot) []model.UserSubscription {
	userSub := make(map[string][]model.Stock)
	for _, doc := range docs {
		var stockSub model.StockSubscription
		err := doc.DataTo(&stockSub)
		if err != nil {
			fmt.Println(err)
			continue
		}

		sub, found := userSub[stockSub.UserID]
		if !found {
			userSub[stockSub.UserID] = []model.Stock{stockSub.Stock}
		} else {
			sub = append(sub, stockSub.Stock)
			userSub[stockSub.UserID] = sub
		}
	}
	userSubscriptions := make([]model.UserSubscription, 0, len(userSub))
	for userID, subs := range userSub {
		userSubscription := model.UserSubscription{
			UserID:      userID,
			Instruments: subs,
		}
		userSubscriptions = append(userSubscriptions, userSubscription)
	}
	return userSubscriptions
}

func distinctInstrumentID(docs []*firestore.DocumentSnapshot) []int {
	ids := make([]int, 0, len(docs))
	for _, doc := range docs {
		id, err := doc.DataAt("InstrumentID")
		if err == nil {
			if idInt, ok := id.(int64); ok {
				ids = append(ids, int(idInt))
			} else {
				fmt.Println("id not int", id)
			}
		} else {
			fmt.Println(err)
		}
	}
	return ids
}