package etoro

import (
	"errors"
	"fmt"
	"time"

	"github.com/tianhai82/stock-timing/httprequester"
	"github.com/tianhai82/stock-timing/model"
)

func RetrieveInstruments() ([]model.InstrumentDisplayData, error) {
	var etoroInstruments model.EtoroInstruments
	err := httprequester.MakeGetRequest("https://api.etorostatic.com/sapi/instrumentsmetadata/V1.1/instruments", &etoroInstruments)
	if err != nil {
		return nil, err
	}
	return etoroInstruments.InstrumentDisplayDatas, nil
}

func RetrieveCandle(instrumentID int, period int) ([]model.Candle, error) {
	url := fmt.Sprintf("https://www.etoro.com/sapi/candles/candles/desc.json/OneDay/%d/%d", period, instrumentID)
	var etoroCandles model.EtoroCandle
	err := httprequester.MakeGetRequest(url, &etoroCandles)
	if err != nil {
		return nil, err
	}
	if len(etoroCandles.Candles) != 1 {
		return nil, errors.New("outer candle count must be 1")
	}
	latestWeekday := getLatestWeekday()
	if etoroCandles.Candles[0].Candles[0].FromDate.Format("2006-01-02") != latestWeekday.Format("2006-01-02") {
		var etoroCandles2 model.EtoroCandle
		err2 := httprequester.MakeGetRequest(url, &etoroCandles2)
		if err2 != nil {
			return reverseCandles(etoroCandles.Candles[0].Candles), nil
		}
		if len(etoroCandles2.Candles) != 1 {
			return reverseCandles(etoroCandles.Candles[0].Candles), nil
		}
		if etoroCandles2.Candles[0].Candles[0].FromDate.After(etoroCandles.Candles[0].Candles[0].FromDate) {
			return reverseCandles(etoroCandles2.Candles[0].Candles), nil
		}
	}
	return reverseCandles(etoroCandles.Candles[0].Candles), nil
}

func reverseCandles(candles []model.Candle) []model.Candle {
	size := len(candles)
	reversed := make([]model.Candle, size)
	for i := range candles {
		reversed[size-i-1] = candles[i]
	}
	return reversed
}

func getLatestWeekday() time.Time {
	now := time.Now()
	for now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		now = now.AddDate(0, 0, -1)
	}
	return now
}
