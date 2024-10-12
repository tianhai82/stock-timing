package rpcs

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/tianhai82/stock-timing/analyzer"
	candleStore "github.com/tianhai82/stock-timing/candle"
	"github.com/tianhai82/stock-timing/firebase"
	"github.com/tianhai82/stock-timing/model"
)

const CandlePeriod = 350

// AddRpcs adds API handlers to the gin router
func AddEtoroRpcs(router *gin.RouterGroup) {
	router.GET("/instruments", retrieveInstruments)
	router.GET("/candles/:instrumentID", retrieveCandles)
	router.GET("/signals/:instrumentID/period/:period/buyLimit/:buyLimit/sellLimit/:sellLimit", analyseInstrument)
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
	periodStr := c.Param("period")
	period, err := strconv.Atoi(periodStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid period")
		return
	}
	buyLimitStr := c.Param("buyLimit")
	buyLimit, err := strconv.ParseFloat(buyLimitStr, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid buy limit")
		return
	}
	sellLimitStr := c.Param("sellLimit")
	sellLimit, err := strconv.ParseFloat(sellLimitStr, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid sell limit")
		return
	}
	candles, err := candleStore.RetrieveCandles(id, CandlePeriod+period)
	if err != nil {
		log.Printf("[analyseInstrument] id:%d, totalPeriod:%d, candleStore.RetrieveCandles failed: %v", id, CandlePeriod+period, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, "retrieval failed")
		return
	}
	advices := make([]model.TradeAdvice, 0)
	for i := 0; i < len(candles)-period+1; i++ {
		candleToAnalyze := candles[i : i+period]
		analysis := analyzer.AnalyzerCandles(candleToAnalyze, buyLimit, sellLimit)
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
	candles, err := candleStore.RetrieveCandles(id, CandlePeriod)
	if err != nil {
		log.Printf("[retrieveCandles] id:%d period:%d candleStore.RetrieveCandles failed:%v", id, CandlePeriod, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, "retrieval failed")
		return
	}
	c.JSON(200, candles)
}

func retrieveInstruments(c *gin.Context) {
	if firebase.StorageClient == nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	bucket, err := firebase.StorageClient.DefaultBucket()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	reader, err := bucket.Object("etoro_tda.json").NewReader(context.Background())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.DataFromReader(200, reader.Attrs.Size, "application/json", reader, make(map[string]string))
}
