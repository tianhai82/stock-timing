package cloudtask

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"github.com/gin-gonic/gin"
	"github.com/tianhai82/stock-timing/analyzer"
	"github.com/tianhai82/stock-timing/etoro"
	"github.com/tianhai82/stock-timing/mail"

	// "github.com/tianhai82/stock-timing/mail"
	"github.com/tianhai82/stock-timing/model"
	"github.com/tianhai82/stock-timing/rpcs"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
)

// AddCloudTasks add cloud task handlers to router
func AddCloudTasks(router *gin.RouterGroup) {
	router.POST("/analyze-stocks", analyzeStock)
}

func analyzeStock(c *gin.Context) {
	start := time.Now()
	ctx := context.Background()
	usersSubscriptions, err := rpcs.FirestoreClient.Collection("usersToEmail").Documents(ctx).GetAll()
	if err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}
	for u, userSubscriptions := range usersSubscriptions {
		var userSubs model.UserSubscription
		err = userSubscriptions.DataTo(&userSubs)
		if err != nil {
			fmt.Println("error convert to UserSubscription", userSubscriptions.Ref.ID)
			_, _ = userSubscriptions.Ref.Delete(ctx)
			continue
		}

		userAnalysises := make([]model.EmailAnalysis, len(userSubs.Subscriptions))
		for i, sub := range userSubs.Subscriptions {
			candles, err := etoro.RetrieveCandle(sub.InstrumentID, sub.Period)
			if err != nil {
				continue
			}
			analysis := analyzer.AnalyzerCandles(candles)
			userAnalysises[i] = model.EmailAnalysis{
				InstrumentDisplayName: sub.InstrumentDisplayName,
				InstrumentSymbol:      sub.Symbol,
				Period:                analysis.Period,
				CurrentPrice:          analysis.CurrentCandle.Close,
			}
			percentile := 0.0
			if math.Abs(analysis.CurrentDev) > analysis.StdDev {
				if analysis.CurrentDev > 0 {
					percentile = ((analysis.CurrentDev-analysis.StdDev)/(analysis.MaxDev-analysis.StdDev))/2 + 0.5
				} else {
					percentile = 0.5 - ((math.Abs(analysis.CurrentDev)-analysis.StdDev)/(analysis.MaxDev-analysis.StdDev))/2
				}
			}
			userAnalysises[i].PricePercentile = percentile

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
		}, []mail.Email{{Email: userSubs.UserID, Name: userSubs.UserID}})
		if err != nil {
			fmt.Println("error sending mail", err)
		}
		_, _ = userSubscriptions.Ref.Delete(ctx)
		if u != len(usersSubscriptions)-1 {
			duration := time.Since(start)
			if duration.Minutes() > 8.0 {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}
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
