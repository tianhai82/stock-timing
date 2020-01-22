package cron

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/tianhai82/stock-timing/cloudtask"
	"github.com/tianhai82/stock-timing/firebase"
	"github.com/tianhai82/stock-timing/model"
)

// AddCronJobs adds cron job handler to router group
func AddCronJobs(router *gin.RouterGroup) {
	router.GET("/stockSubscriptions", stockSubscriptions)
}

func stockSubscriptions(c *gin.Context) {
	ctx := context.Background()
	iter := firebase.FirestoreClient.Collection("subscription").Documents(ctx)
	docs, err := iter.GetAll()
	if err != nil {
		fmt.Println("unable to retrieve subscriptions", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	userSubscriptions := groupUsers(docs)
	usersToEmailCol := firebase.FirestoreClient.Collection("usersToEmail")
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
}

func groupUsers(docs []*firestore.DocumentSnapshot) []model.UserSubscription {
	userSub := make(map[string][]model.StockSubscription)
	for _, doc := range docs {
		var stockSub model.StockSubscription
		err := doc.DataTo(&stockSub)
		if err != nil {
			fmt.Println(err)
			continue
		}

		sub, found := userSub[stockSub.UserID]
		if !found {
			userSub[stockSub.UserID] = []model.StockSubscription{stockSub}
		} else {
			sub = append(sub, stockSub)
			userSub[stockSub.UserID] = sub
		}
	}
	userSubscriptions := make([]model.UserSubscription, 0, len(userSub))
	for userID, subs := range userSub {
		userSubscription := model.UserSubscription{
			UserID:        userID,
			Subscriptions: subs,
		}
		userSubscriptions = append(userSubscriptions, userSubscription)
	}
	return userSubscriptions
}
