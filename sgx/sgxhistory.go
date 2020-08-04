package sgx

import (
	"errors"
	"fmt"
	"time"

	"github.com/tianhai82/stock-timing/httprequester"
	"github.com/tianhai82/stock-timing/model"
)

const quoteUrl = "https://api.sgx.com/securities/v1.1/%s/code/%s?params=nc,adjusted-vwap,b,bv,p,c,change_vs_pc,change_vs_pc_percentage,cx,cn,dp,dpc,du,ed,fn,h,iiv,iopv,lt,l,o,p_,pv,ptd,s,sv,trading_time,v_,v,vl,vwap,vwap-currency"
const historyUrl = "https://api.sgx.com/securities/v1.1/charts/historic/%s/code/%s/1y"

type sgxPrice struct {
	High        float64 `json:"h"`
	Low         float64 `json:"l"`
	Close       float64 `json:"lt"`
	Name        string  `json:"n"`
	Symbol      string  `json:"nc"`
	Open        float64 `json:"o"`
	TradingTime string  `json:"trading_time"`
	Type        string  `json:"type"`
}
type meta struct {
	Code           string `json:"code"`
	Message        string `json:"message"`
	ProcessedTime  int    `json:"processedTime"`
	ProcessingTime string `json:"processingTime"`
}
type sgxResp struct {
	Data struct {
		Prices []sgxPrice `json:"prices"`
	} `json:"data"`
	Meta meta `json:"meta"`
}
type sgxHistoryResp struct {
	Data struct {
		Historic []sgxPrice `json:"historic"`
	} `json:"data"`
	Meta meta `json:"meta"`
}

func RetrieveHistory(symbol model.InstrumentDisplayData, period int) ([]model.Candle, error) {
	historyQuoteURL := fmt.Sprintf(historyUrl, symbol.Type, symbol.SymbolFull)
	var historyResp sgxHistoryResp
	err := httprequester.MakeGetRequest(historyQuoteURL, &historyResp)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if historyResp.Meta.Code != "200" {
		fmt.Println(historyResp.Meta.Message)
		return nil, errors.New(historyResp.Meta.Message)
	}

	if len(historyResp.Data.Historic) == 0 {
		fmt.Println("empty history returned from SGX")
		return nil, errors.New("empty history returned from SGX")
	}
	prices := historyResp.Data.Historic

	t, err := time.Parse("20060102_150405", historyResp.Data.Historic[len(historyResp.Data.Historic)-1].TradingTime)
	now := time.Now()
	if now.Weekday() == time.Saturday {
		now = now.AddDate(0, 0, -1)
	} else if now.Weekday() == time.Sunday {
		now = now.AddDate(0, 0, -2)
	}
	if t.Day() != now.Day() {
		currentQuoteURL := fmt.Sprintf(quoteUrl, symbol.Type, symbol.SymbolFull)
		var currentQuoteResp sgxResp
		err = httprequester.MakeGetRequest(currentQuoteURL, &currentQuoteResp)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		if currentQuoteResp.Meta.Code != "200" {
			fmt.Println(currentQuoteResp.Meta.Message)
			return nil, errors.New(currentQuoteResp.Meta.Message)
		}
		prices = append(prices, currentQuoteResp.Data.Prices...)
	}
	if len(prices) > period {
		prices = prices[len(prices)-period:]
	}
	candles := convertPricesToCandles(prices)
	lastCandle := candles[len(candles)-1]
	if lastCandle.Open != 0.0 {
		if lastCandle.High == 0.0 {
			lastCandle.High = lastCandle.Open
		}
		if lastCandle.Low == 0.0 {
			lastCandle.Low = lastCandle.Open
		}
		if lastCandle.Close == 0.0 {
			lastCandle.Close = lastCandle.Open
		}
		candles[len(candles)-1] = lastCandle
	} else if lastCandle.Close != 0.0 {
		if lastCandle.High == 0.0 {
			lastCandle.High = lastCandle.Close
		}
		if lastCandle.Low == 0.0 {
			lastCandle.Low = lastCandle.Close
		}
		if lastCandle.Open == 0.0 {
			lastCandle.Open = lastCandle.Close
		}
		candles[len(candles)-1] = lastCandle
	}
	return candles, nil
}

func convertPricesToCandles(prices []sgxPrice) []model.Candle {
	candles := make([]model.Candle, 0, len(prices))
	for _, price := range prices {
		t, err := time.Parse("20060102_150405", price.TradingTime)
		if err != nil {
			fmt.Println(err)
			continue
		}
		candle := model.Candle{
			FromDate: t,
			Open:     price.Open,
			High:     price.High,
			Low:      price.Low,
			Close:    price.Close,
		}
		candles = append(candles, candle)
	}
	return candles
}
