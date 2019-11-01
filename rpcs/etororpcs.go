package rpcs

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	firebase "firebase.google.com/go"
	auth "firebase.google.com/go/auth"
	storage "firebase.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"github.com/tianhai82/stock-timing/analyzer"
	"github.com/tianhai82/stock-timing/etoro"
	"github.com/tianhai82/stock-timing/model"
)

var config = &firebase.Config{
	StorageBucket: "stock-timing.appspot.com",
}
var app *firebase.App
var AuthClient *auth.Client
var StorageClient *storage.Client

func init() {
	var err error
	ctx := context.Background()
	app, err = firebase.NewApp(ctx, config)
	if err != nil {
		fmt.Println(err)
		return
	}
	AuthClient, err = app.Auth(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	StorageClient, err = app.Storage(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
}

const period = 24
const candlePeriod = 250

// AddRpcs adds API handlers to the gin router
func AddEtoroRpcs(router *gin.RouterGroup) {
	router.GET("/instruments", retrieveInstruments)
	router.GET("/candles/:instrumentID", retrieveCandles)
	router.GET("/signals/:instrumentID", analyseInstrument)
}

func analyseInstrument(c *gin.Context) {
	instrumentID := c.Param("instrumentID")
	if instrumentID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid instrument ID")
		return
	}
	id, err := strconv.Atoi(instrumentID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid instrument ID")
		return
	}
	candles, err := etoro.RetrieveCandle(id, candlePeriod+period)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, "retrieval failed")
		return
	}
	advices := make([]model.TradeAdvice, 0)
	for i := 0; i < len(candles)-period; i++ {
		analysis := analyzer.AnalyzerCandles(candles[i : i+period])
		if analysis.Signal == model.Buy || analysis.Signal == model.Sell {
			advice := model.TradeAdvice{
				Date:   analysis.CurrentCandle.FromDate,
				Price:  analysis.CurrentCandle.Close,
				Signal: analysis.Signal,
			}
			advices = append(advices, advice)
		}
	}
	c.JSON(200, advices)
}

func retrieveCandles(c *gin.Context) {
	instrumentID := c.Param("instrumentID")
	if instrumentID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid instrument ID")
		return
	}
	id, err := strconv.Atoi(instrumentID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid instrument ID")
		return
	}
	candles, err := etoro.RetrieveCandle(id, candlePeriod)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, "retrieval failed")
		return
	}
	c.JSON(200, candles)
}

func retrieveInstruments(c *gin.Context) {
	if StorageClient == nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	bucket, err := StorageClient.DefaultBucket()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	reader, err := bucket.Object("etoro_stocks.json").NewReader(context.Background())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.DataFromReader(200, reader.Attrs.Size, "application/json", reader, make(map[string]string))
}
