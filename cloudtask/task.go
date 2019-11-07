package cloudtask

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"github.com/gin-gonic/gin"
	"github.com/tianhai82/stock-timing/analyzer"
	"github.com/tianhai82/stock-timing/etoro"
	"github.com/tianhai82/stock-timing/mail"
	"github.com/tianhai82/stock-timing/model"
	"github.com/tianhai82/stock-timing/rpcs"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
)

// AddCloudTasks add cloud task handlers to router
func AddCloudTasks(router *gin.RouterGroup) {
	router.POST("/analyze-stocks", analyzeStock)
	router.POST("/email-subscribers", emailSubscribers)
}

func emailSubscribers(c *gin.Context) {
	start := time.Now()
	ctx := context.Background()
	users, err := rpcs.FirestoreClient.Collection("usersToEmail").Documents(ctx).GetAll()
	if err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}
	for u, user := range users {
		var userSub model.UserSubscription
		errSub := user.DataTo(&userSub)
		idStr := user.Ref.ID
		if errSub != nil {
			fmt.Println("fail to convert to userSub", errSub)
			_, _ = rpcs.FirestoreClient.Collection("usersToEmail").Doc(idStr).Delete(ctx)
			continue
		}
		userAnalysises := make([]model.EmailAnalysis, len(userSub.Instruments))
		for i, instrument := range userSub.Instruments {
			doc, errGet := rpcs.FirestoreClient.Collection("dailyAnalysis").Doc(strconv.Itoa(instrument.InstrumentID)).Get(ctx)
			if errGet != nil {
				fmt.Println("fail to get daily analysis", instrument.Symbol, errGet)
				continue
			}
			var analysis model.TradeAnalysis
			errTo := doc.DataTo(&analysis)
			if errTo != nil {
				fmt.Println("fail to convert data to TradeAnalysis", instrument.Symbol, errTo)
				continue
			}
			userAnalysises[i] = model.EmailAnalysis{
				InstrumentDisplayName: instrument.InstrumentDisplayName,
				InstrumentSymbol:      instrument.Symbol,
				Period:                analysis.Period,
				Mean:                  analysis.Mean,
				StdDev:                analysis.StdDev,
				MaxDev:                analysis.MaxDev,
				LimitDev:              analysis.LimitDev,
				CurrentDev:            analysis.CurrentDev,
				CurrentPrice:          analysis.CurrentCandle.Close,
			}
			if analysis.Signal == model.Buy {
				userAnalysises[i].BuyOrSell = "Buy"
			} else if analysis.Signal == model.Sell {
				userAnalysises[i].BuyOrSell = "Sell"
			} else {
				userAnalysises[i].BuyOrSell = "Hold"
			}
		}
		mailApiKey, _ := os.LookupEnv("MAIL_API_KEY")
		err = mail.Sendmail(mailApiKey, 1, gin.H{
			"analysises": userAnalysises,
		}, []mail.Email{{Email: userSub.UserID, Name: userSub.UserID}})
		if err != nil {
			fmt.Println("error sending mail", err)
		}
		if u != len(users)-1 {
			duration := time.Since(start)
			if duration.Minutes() > 8.0 {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}
	}
	c.Status(200)
}

func analyzeStock(c *gin.Context) {
	start := time.Now()
	ctx := context.Background()
	instruments, err := rpcs.FirestoreClient.Collection("instrumentsToAnalyse").Documents(ctx).GetAll()
	if err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}
	for i, instrument := range instruments {
		idStr := instrument.Ref.ID
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Println("error convert id to string", idStr)
			_, _ = rpcs.FirestoreClient.Collection("instrumentsToAnalyse").Doc(idStr).Delete(ctx)
			continue
		}
		candles, err := etoro.RetrieveCandle(id, rpcs.Period)
		if err != nil {
			fmt.Println("error retrieving candles", id, err)
			_, _ = rpcs.FirestoreClient.Collection("instrumentsToAnalyse").Doc(idStr).Delete(ctx)
			continue
		}
		analysis := analyzer.AnalyzerCandles(candles)
		_, err = rpcs.FirestoreClient.Collection("dailyAnalysis").Doc(strconv.Itoa(analysis.CurrentCandle.InstrumentID)).Set(ctx, analysis)
		if err != nil {
			fmt.Println("fail to save analysis", id, err)
		}

		_, _ = rpcs.FirestoreClient.Collection("instrumentsToAnalyse").Doc(idStr).Delete(ctx)
		if i != len(instruments)-1 {
			duration := time.Since(start)
			if duration.Minutes() > 8.0 {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}
	}
	_, err = CreateTask("stock-timing", "asia-south1", "email-subscribers", nil)
	if err != nil {
		fmt.Println("error creating email subscribers task", err)
	}
	c.Status(200)
}

// CreateTask creates a new task in your App Engine queue.
func CreateTask(projectID, locationID, queueID string, message gin.H) (*taskspb.Task, error) {
	// Create a new Cloud Tasks client instance.
	// See https://godoc.org/cloud.google.com/go/cloudtasks/apiv2
	ctx := context.Background()
	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewClient: %v", err)
	}

	// Build the Task queue path.
	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", projectID, locationID, queueID)

	// Build the Task payload.
	// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2#CreateTaskRequest
	req := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2#AppEngineHttpRequest
			MessageType: &taskspb.Task_AppEngineHttpRequest{
				AppEngineHttpRequest: &taskspb.AppEngineHttpRequest{
					HttpMethod:  taskspb.HttpMethod_POST,
					RelativeUri: "/task/" + queueID,
				},
			},
		},
	}

	if message != nil {
		// Add a payload message if one is present.
		b, err := json.Marshal(message)
		if err == nil {
			req.Task.GetAppEngineHttpRequest().Body = b
		}
	}

	createdTask, err := client.CreateTask(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cloudtasks.CreateTask: %v", err)
	}

	return createdTask, nil
}
