package yahoo

import (
	"fmt"
	"time"

	"github.com/tianhai82/stock-timing/httprequester"
	"github.com/tianhai82/stock-timing/model"
)

const yahooUrl = "https://query1.finance.yahoo.com/v8/finance/chart/CSV?symbol=%s&period1=%d&period2=%d&interval=1d"

//period is a 10 digit number (second since epoch)

type yahooResp struct {
	Chart struct {
		Result []yahooResult `json:"result"`
		Error  interface{}   `json:"error"`
	} `json:"chart"`
}
type yahooResult struct {
	Meta struct {
		Currency             string  `json:"currency"`
		Symbol               string  `json:"symbol"`
		ExchangeName         string  `json:"exchangeName"`
		InstrumentType       string  `json:"instrumentType"`
		FirstTradeDate       int64   `json:"firstTradeDate"`
		RegularMarketTime    int64   `json:"regularMarketTime"`
		Gmtoffset            int64   `json:"gmtoffset"`
		Timezone             string  `json:"timezone"`
		ExchangeTimezoneName string  `json:"exchangeTimezoneName"`
		RegularMarketPrice   float64 `json:"regularMarketPrice"`
		ChartPreviousClose   float64 `json:"chartPreviousClose"`
		PriceHint            int     `json:"priceHint"`
		CurrentTradingPeriod struct {
			Pre struct {
				Timezone  string `json:"timezone"`
				Start     int64  `json:"start"`
				End       int64  `json:"end"`
				Gmtoffset int64  `json:"gmtoffset"`
			} `json:"pre"`
			Regular struct {
				Timezone  string `json:"timezone"`
				Start     int64  `json:"start"`
				End       int64  `json:"end"`
				Gmtoffset int64  `json:"gmtoffset"`
			} `json:"regular"`
			Post struct {
				Timezone  string `json:"timezone"`
				Start     int64  `json:"start"`
				End       int64  `json:"end"`
				Gmtoffset int64  `json:"gmtoffset"`
			} `json:"post"`
		} `json:"currentTradingPeriod"`
		DataGranularity string   `json:"dataGranularity"`
		Range           string   `json:"range"`
		ValidRanges     []string `json:"validRanges"`
	} `json:"meta"`
	Timestamp  []int64 `json:"timestamp"`
	Indicators struct {
		Quote []struct {
			Close  []float64 `json:"close"`
			High   []float64 `json:"high"`
			Low    []float64 `json:"low"`
			Volume []int64   `json:"volume"`
			Open   []float64 `json:"open"`
		} `json:"quote"`
		Adjclose []struct {
			Adjclose []float64 `json:"adjclose"`
		} `json:"adjclose"`
	} `json:"indicators"`
}

func RetrieveHistory(symbol string, days int) ([]model.Candle, error) {
	now := time.Now()
	period2 := now.UnixNano() / 1000000000
	daysAgo := days * 7 / 5
	from := now.AddDate(0, 0, -daysAgo)
	period1 := from.UnixNano() / 1000000000

	url := fmt.Sprintf(yahooUrl, symbol, period1, period2)
	var resp yahooResp
	err := httprequester.MakeGetRequest(url, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Chart.Error != nil {
		return nil, fmt.Errorf("Error in query response detected")
	}
	return convertToCandles(resp.Chart.Result[0])
}
func convertToCandles(res yahooResult) ([]model.Candle, error) {
	if len(res.Indicators.Quote) != 1 || len(res.Indicators.Adjclose) != 1 {
		return nil, fmt.Errorf("len of quote/adjclose not equal to 1")
	}
	if len(res.Timestamp) != len(res.Indicators.Quote[0].Close) ||
		len(res.Timestamp) != len(res.Indicators.Adjclose[0].Adjclose) ||
		len(res.Timestamp) != len(res.Indicators.Quote[0].High) ||
		len(res.Timestamp) != len(res.Indicators.Quote[0].Low) ||
		len(res.Timestamp) != len(res.Indicators.Quote[0].Open) ||
		len(res.Timestamp) != len(res.Indicators.Quote[0].Volume) {
		return nil, fmt.Errorf("invalid yahoo finance query results: len of timestamps, quotes and adjclose not equal")
	}
	candles := make([]model.Candle, 0, len(res.Timestamp))
	for i, ts := range res.Timestamp {
		if res.Indicators.Quote[0].Volume[i] == 0 &&
			res.Indicators.Quote[0].Open[i] == 0 &&
			res.Indicators.Quote[0].Close[i] == 0 {
			continue
		}
		candles = append(candles, model.Candle{
			FromDate: time.Unix(ts+res.Meta.Gmtoffset, 0),
			Open:     res.Indicators.Quote[0].Open[i],
			Close:    res.Indicators.Quote[0].Close[i],
			High:     res.Indicators.Quote[0].High[i],
			Low:      res.Indicators.Quote[0].Low[i],
		})
	}
	return candles, nil
}
