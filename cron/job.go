package cron

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/tianhai82/stock-timing/model"
	"github.com/tianhai82/stock-timing/rpcs"
)

func AddCronJobs(router *gin.RouterGroup) {
	router.GET("/stockSubscriptions", stockSubscriptions)
}

func stockSubscriptions(c *gin.Context) {
	ctx := context.Background()
	iter := rpcs.FirestoreClient.Collection("subscription").Documents(ctx)
	docs, err := iter.GetAll()
	if err != nil {
		fmt.Println(err)
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
	// cloudtask.CreateTask()
	// create cloud task for Analyse Stock
}

func groupUsers(docs []*firestore.DocumentSnapshot) []model.UserSubscription {
	userSub := make(map[string][]int)
	for _, doc := range docs {
		id, err := doc.DataAt("InstrumentID")
		if err != nil {
			fmt.Println(err)
			continue
		}
		idInt, ok := id.(int64)
		if !ok {
			fmt.Println("id not int64", id)
			continue
		}
		userID, err := doc.DataAt("UserID")
		if err != nil {
			continue
		}
		userIDStr, ok := userID.(string)
		if !ok {
			continue
		}
		sub, found := userSub[userIDStr]
		if !found {
			userSub[userIDStr] = []int{int(idInt)}
		} else {
			sub = append(sub, int(idInt))
			userSub[userIDStr] = sub
		}
	}
	userSubscriptions := make([]model.UserSubscription, 0, len(userSub))
	for userID, subs := range userSub {
		userSubscription := model.UserSubscription{
			UserID:        userID,
			InstrumentIDs: subs,
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
