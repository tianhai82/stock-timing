package rpcs

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tianhai82/stock-timing/model"
	"github.com/tianhai82/stock-timing/tda"
)

func AddTdaRpcs(router *gin.RouterGroup) {
	router.GET("/pricehistory/:symbol/:frequencyType/:period", retrievePriceHistory)
	router.GET("/optionchain/:symbol/:contractType/:range/:fromDate/:toDate", retrieveOptions)
}
func retrieveOptions(c *gin.Context) {
	symbol := c.Param("symbol")
	contractType := c.Param("contractType")
	optionRange := c.Param("range")
	fromDate := c.Param("fromDate")
	toDate := c.Param("toDate")
	chains, err := tda.RetrieveOptionChain(symbol, contractType, optionRange, fromDate, toDate)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(200, chains)
}

func retrievePriceHistory(c *gin.Context) {
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
		candles, err := tda.RetrieveHistory(model.InstrumentDisplayData{
			SymbolFull: symbol,
		}, periodInt)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(200, candles)
	} else if frequencyType == "weekly" {
		candles, err := tda.RetrieveWeeklyHistory(model.InstrumentDisplayData{
			SymbolFull: symbol,
		}, periodInt)
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
