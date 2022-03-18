package cloudtask

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"github.com/gin-gonic/gin"
	"github.com/tianhai82/stock-timing/analyzer"
	candleStore "github.com/tianhai82/stock-timing/candle"
	"github.com/tianhai82/stock-timing/firebase"
	"github.com/tianhai82/stock-timing/mail"
	"github.com/tianhai82/stock-timing/telegram"

	// "github.com/tianhai82/stock-timing/mail"
	"github.com/tianhai82/stock-timing/model"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
)

// AddCloudTasks add cloud task handlers to router
func AddCloudTasks(router *gin.RouterGroup) {
	router.POST("/analyze-stocks", analyzeStock)
}

func analyzeStock(c *gin.Context) {
	start := time.Now()
	ctx := context.Background()
	usersSubscriptions, err := firebase.FirestoreClient.Collection("usersToEmail").Documents(ctx).GetAll()
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

		highPrices := make([]model.EmailAnalysis, 0, len(userSubs.Subscriptions))
		lowPrices := make([]model.EmailAnalysis, 0, len(userSubs.Subscriptions))
		for _, sub := range userSubs.Subscriptions {
			candles, err := candleStore.RetrieveCandles(sub.InstrumentID, sub.Period)
			if err != nil {
				continue
			}
			buyLimit := 0.55
			if sub.BuyLimit != 0.0 {
				buyLimit = sub.BuyLimit
			}
			sellLimit := 0.57
			if sub.SellLimit != 0.0 {
				sellLimit = sub.SellLimit
			}
			analysis := analyzer.AnalyzerCandles(candles, buyLimit, sellLimit)
			ua := model.EmailAnalysis{
				InstrumentDisplayName: sub.InstrumentDisplayName,
				InstrumentSymbol:      sub.Symbol,
				InstrumentID:          sub.InstrumentID,
				BuyFreq:               limitToFreq(sub.BuyLimit),
				SellFreq:              limitToFreq(sub.SellLimit),
				Period:                analysis.Period,
				CurrentPrice:          analysis.CurrentCandle.Close,
			}
			percentile := 0.5
			if math.Abs(analysis.CurrentDev) > analysis.StdDev {
				if analysis.CurrentDev > 0 {
					percentile = ((analysis.CurrentDev-analysis.StdDev)/(analysis.MaxDev-analysis.StdDev))/2 + 0.5
				} else {
					percentile = 0.5 - ((math.Abs(analysis.CurrentDev)-analysis.StdDev)/(analysis.MaxDev-analysis.StdDev))/2
				}
			}
			ua.PricePercentile = percentile

			if analysis.Signal == model.Buy {
				ua.BuyOrSell = "Buy"
			} else if analysis.Signal == model.Sell {
				ua.BuyOrSell = "Sell"
			} else {
				ua.BuyOrSell = "Hold"
			}
			if ua.PricePercentile > 0.5 {
				highPrices = append(highPrices, ua)
			} else if ua.PricePercentile < 0.5 {
				lowPrices = append(lowPrices, ua)
			}
		}
		sort.Slice(highPrices, func(i, j int) bool {
			if highPrices[i].BuyOrSell == highPrices[j].BuyOrSell {
				if highPrices[i].PricePercentile > highPrices[j].PricePercentile {
					return false
				}
				return true
			}
			if highPrices[i].BuyOrSell != "Hold" && highPrices[j].BuyOrSell == "Hold" {
				return false
			}
			return true
		})
		sort.Slice(lowPrices, func(i, j int) bool {
			if lowPrices[i].BuyOrSell == lowPrices[j].BuyOrSell {
				if lowPrices[i].PricePercentile < lowPrices[j].PricePercentile {
					return false
				}
				return true
			}
			if lowPrices[i].BuyOrSell != "Hold" && lowPrices[j].BuyOrSell == "Hold" {
				return false
			}
			return true
		})

		highMd := formMarkdownMsg(highPrices, "HIGH")
		telegram.SendMessage(highMd, telegram.CHAT_ID, telegram.TOKEN)
		lowMd := formMarkdownMsg(lowPrices, "LOW")
		telegram.SendMessage(lowMd, telegram.CHAT_ID, telegram.TOKEN)
		log.Default().Println("highMd  |" + highMd)
		log.Default().Println("lowMd  |" + lowMd)

		mailApiKey, _ := os.LookupEnv("MAIL_API_KEY")
		err = mail.Sendmail(mailApiKey, 1, gin.H{
			"highPrices": highPrices,
			"lowPrices":  lowPrices,
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

func formMarkdownMsg(analysis []model.EmailAnalysis, highLow string) string {
	md := fmt.Sprintf("*%s*\n", highLow)
	for _, p := range analysis {
		md += fmt.Sprintf(`<a href="%s">%s</a> - %s<br>Current Price: %.3f<br>Price percentile: %.3f<br>`,
			instrumentURL(p), p.InstrumentDisplayName, p.BuyOrSell, p.CurrentPrice, p.PricePercentile)
	}
	md += "<br>"
	return md
}

func instrumentURL(an model.EmailAnalysis) string {
	return fmt.Sprintf("https://stock-timing.web.app/#/%d/%d/%.0f/%.0f", an.InstrumentID, an.Period, an.BuyFreq, an.SellFreq)
}

func limitToFreq(limit float64) float64 {
	return 100 - (limit-0.25)*200
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
