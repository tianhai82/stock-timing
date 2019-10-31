package rpcs

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"github.com/tianhai82/stock-timing/etoro"
)

var config = &firebase.Config{
	StorageBucket: "stock-timing.appspot.com",
}

// AddRpcs adds API handlers to the gin router
func AddRpcs(router *gin.RouterGroup) {
	router.POST("/instruments", retrieveInstruments)
	router.POST("/candles/:instrumentID", retrieveCandles)
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
	candles, err := etoro.RetrieveCandle(id, 60)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, "retrieval failed")
		return
	}
	c.JSON(200, candles)
}

func retrieveInstruments(c *gin.Context) {
	app, err := firebase.NewApp(context.Background(), config)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	client, err := app.Storage(context.Background())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	bucket, err := client.DefaultBucket()
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
