package yahoo

import (
	"fmt"
	"github.com/tianhai82/stock-timing/httprequester"
	"time"
)

const yahooUrl = "https://query1.finance.yahoo.com/v8/finance/chart/CSV?symbol=%s&period1=%d&period2=%d&interval=1d"

//period is a 10 digit number (millisecond since epoch)
type YahooResp struct {
	Chart struct {
		Result []struct {
			Meta struct {
				Currency             string  `json:"currency"`
				Symbol               string  `json:"symbol"`
				ExchangeName         string  `json:"exchangeName"`
				InstrumentType       string  `json:"instrumentType"`
				FirstTradeDate       int     `json:"firstTradeDate"`
				RegularMarketTime    int     `json:"regularMarketTime"`
				Gmtoffset            int     `json:"gmtoffset"`
				Timezone             string  `json:"timezone"`
				ExchangeTimezoneName string  `json:"exchangeTimezoneName"`
				RegularMarketPrice   float64 `json:"regularMarketPrice"`
				ChartPreviousClose   float64 `json:"chartPreviousClose"`
				PriceHint            int     `json:"priceHint"`
				CurrentTradingPeriod struct {
					Pre struct {
						Timezone  string `json:"timezone"`
						Start     int    `json:"start"`
						End       int    `json:"end"`
						Gmtoffset int    `json:"gmtoffset"`
					} `json:"pre"`
					Regular struct {
						Timezone  string `json:"timezone"`
						Start     int    `json:"start"`
						End       int    `json:"end"`
						Gmtoffset int    `json:"gmtoffset"`
					} `json:"regular"`
					Post struct {
						Timezone  string `json:"timezone"`
						Start     int    `json:"start"`
						End       int    `json:"end"`
						Gmtoffset int    `json:"gmtoffset"`
					} `json:"post"`
				} `json:"currentTradingPeriod"`
				DataGranularity string   `json:"dataGranularity"`
				Range           string   `json:"range"`
				ValidRanges     []string `json:"validRanges"`
			} `json:"meta"`
			Timestamp  []int `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Close  []float64 `json:"close"`
					High   []float64 `json:"high"`
					Low    []float64 `json:"low"`
					Volume []int     `json:"volume"`
					Open   []float64 `json:"open"`
				} `json:"quote"`
				Adjclose []struct {
					Adjclose []float64 `json:"adjclose"`
				} `json:"adjclose"`
			} `json:"indicators"`
		} `json:"result"`
		Error interface{} `json:"error"`
	} `json:"chart"`
}

func RetrieveHistory(symbol string, days int) {
	now := time.Now()
	period2 := now.UnixNano() / 1000000000
	daysAgo := days * 7 / 5
	from := now.AddDate(0, 0, -daysAgo)
	fmt.Println(from, now, daysAgo)
	period1 := from.UnixNano() / 1000000000

	url := fmt.Sprintf(yahooUrl, symbol, period1, period2)
	var resp YahooResp
	err := httprequester.MakeGetRequest(url, &resp)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.Chart.Error)
	fmt.Println(resp.Chart.Result[0])
}
