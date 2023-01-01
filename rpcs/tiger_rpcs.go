package rpcs

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tianhai82/stock-timing/model"
	tigerClient "github.com/tianhai82/stock-timing/tiger/client"
)

func AddTigerRpcs(router *gin.RouterGroup) {
	router.GET("/pricehistory/:symbol/:frequencyType/:period", retrieveTigerPriceHistory)
	router.GET("/optionchain/:symbol/:contractType/:range/:fromDate/:toDate", retrieveTigerOptions)
}

func retrieveTigerOptions(c *gin.Context) {
	symbol := c.Param("symbol")
	contractType := c.Param("contractType")
	optionRange := c.Param("range")
	fromDate := c.Param("fromDate")
	toDate := c.Param("toDate")
	log.Printf("[retrieveTigerOptions] symbol:%s contractType:%s optionRange:%s fromDate:%s toDate:%s", symbol, contractType, optionRange, fromDate, toDate)
	c.AbortWithStatusJSON(500, "not implemented yet")
}

func retrieveTigerPriceHistory(c *gin.Context) {
	symbol := c.Param("symbol")
	frequencyType := c.Param("frequencyType")
	period := c.Param("period")
	periodInt, err := strconv.Atoi(period)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if frequencyType == "daily" {
		candles, err := tigerClient.RetrieveHistory(model.InstrumentDisplayData{
			SymbolFull: symbol,
		}, periodInt, "day")
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(200, candles)
	} else if frequencyType == "weekly" {
		candles, err := tigerClient.RetrieveHistory(model.InstrumentDisplayData{
			SymbolFull: symbol,
		}, periodInt, "week")
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(200, candles)
	} else {
		fmt.Println("frequencyType invalid", frequencyType)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
}
