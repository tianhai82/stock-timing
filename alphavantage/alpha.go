package alpha

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/tianhai82/stock-timing/httprequester"
	"github.com/tianhai82/stock-timing/model"
)

const (
	token = "1ZE9CMTPW8D9YIED"

	DailyURL  = "https://www.alphavantage.co/query?function=TIME_SERIES_DAILY_ADJUSTED&symbol=%s&apikey=%s"
	WeeklyURL = "https://www.alphavantage.co/query?function=TIME_SERIES_WEEKLY_ADJUSTED&symbol=%s&apikey=%s"
)

type DailyData struct {
	MetaData        MetaData         `json:"Meta Data"`
	TimeSeriesDaily map[string]Quote `json:"Time Series (Daily)"`
}
type WeeklyData struct {
	MetaData         MetaData         `json:"Meta Data"`
	TimeSeriesWeekly map[string]Quote `json:"Weekly Adjusted Time Series"`
}
type MetaData struct {
	Symbol        string `json:"2. Symbol"`
	LastRefreshed string `json:"3. Last Refreshed"`
}
type Quote struct {
	Open             string `json:"1. open"`
	High             string `json:"2. high"`
	Low              string `json:"3. low"`
	Close            string `json:"4. close"`
	AdjustedClose    string `json:"5. adjusted close"`
	Volume           string `json:"6. volume"`
	DividendAmount   string `json:"7. dividend amount"`
	SplitCoefficient string `json:"8. split coefficient"`
}

func RetrieveHistory(ticker string, days int) ([]model.Candle, error) {
	if days > 100 {
		return nil, fmt.Errorf("days greater than 100 is not supported")
	}
	var out DailyData
	err := httprequester.MakeGetRequest(fmt.Sprintf(DailyURL, ticker, token), &out)
	if err != nil {
		return nil, fmt.Errorf("RetrieveHistory: %w", err)
	}
	b, err := json.Marshal(out)
	if err != nil {
		panic(err)
	}
	log.Printf(string(b))
	return nil, nil
}

func RetrieveWeeklyHistory(ticker string, weeks int) ([]model.Candle, error) {
	if weeks > 100 {
		return nil, fmt.Errorf("weeks greater than 100 is not supported")
	}
	var out WeeklyData
	err := httprequester.MakeGetRequest(fmt.Sprintf(WeeklyURL, ticker, token), &out)
	if err != nil {
		return nil, fmt.Errorf("RetrieveWeeklyHistory: %w", err)
	}
	b, err := json.Marshal(out)
	if err != nil {
		panic(err)
	}
	log.Printf(string(b))
	return nil, nil
}
