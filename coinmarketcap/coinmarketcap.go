package coinmarketcap

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"

	"github.com/tianhai82/stock-timing/httprequester"
	"github.com/tianhai82/stock-timing/model"
)

const coinUrl = "https://api.coinmarketcap.com/data-api/v3/cryptocurrency/detail/chart?id=%d&range=1Y"

type Response struct {
	Data struct {
		Points map[string]Values `json:"points"`
	} `json:"data"`
	Status struct {
		Timestamp    time.Time `json:"timestamp"`
		ErrorCode    string    `json:"error_code"`
		ErrorMessage string    `json:"error_message"`
		Elapsed      string    `json:"elapsed"`
		CreditCount  int       `json:"credit_count"`
	} `json:"status"`
}

type Values struct {
	V []float64 `json:"v"`
}

func RetrieveHistory(instrument model.InstrumentDisplayData, period int) ([]model.Candle, error) {
	var resp Response
	err := httprequester.MakeGetRequest(
		fmt.Sprintf(coinUrl, instrument.InstrumentID-2000000),
		&resp)
	if err != nil {
		return nil, err
	}
	if resp.Status.ErrorCode != "0" {
		return nil, fmt.Errorf(resp.Status.ErrorMessage)
	}
	timestamps := make([]string, 0, len(resp.Data.Points))
	for timestamp := range resp.Data.Points {
		timestamps = append(timestamps, timestamp)
	}
	sort.Strings(timestamps)

	candles := make([]model.Candle, 0, len(timestamps))
	for i, timestamp := range timestamps {
		point := resp.Data.Points[timestamp]
		intTime, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			return nil, err
		}
		t := time.Unix(intTime, 0)
		close := point.V[0]
		if i != len(timestamps)-1 {
			close = resp.Data.Points[timestamps[i+1]].V[0]
		}
		candle := model.Candle{
			InstrumentID: instrument.InstrumentID,
			FromDate:     t,
			Open:         point.V[0],
			High:         math.Max(point.V[0], close),
			Low:          math.Min(point.V[0], close),
			Close:        close,
		}
		candles = append(candles, candle)
	}
	candles = MergeToDays(candles)
	if len(candles) > period {
		candles = candles[len(candles)-period:]
	}
	return candles, nil
}

func MergeToDays(candles []model.Candle) []model.Candle {
	outCandles := make([]model.Candle, 0, len(candles))
	var tempCandles []model.Candle
	for i, candle := range candles {
		tempCandles = append(tempCandles, candle)
		if i == 0 {
			continue
		}
		day := candle.FromDate.Day()
		prevCandleDay := candles[i-1].FromDate.Day()
		if day != prevCandleDay {
			outCandles = append(outCandles, mergeCandles(tempCandles))
			tempCandles = []model.Candle{}
		}
	}
	if len(tempCandles) > 0 {
		outCandles = append(outCandles, mergeCandles(tempCandles))
	}
	return outCandles
}

func mergeCandles(candles []model.Candle) model.Candle {
	if len(candles) == 1 {
		return candles[0]
	}
	open := candles[0].Open
	close := candles[len(candles)-1].Close
	high := highest(candles)
	low := lowest(candles)
	return model.Candle{
		InstrumentID: candles[0].InstrumentID,
		FromDate:     candles[0].FromDate,
		Open:         open,
		High:         high,
		Low:          low,
		Close:        close,
	}
}

func highest(candles []model.Candle) float64 {
	max := candles[0].High
	for _, candle := range candles {
		if candle.High > max {
			max = candle.High
		}
	}
	return max
}
func lowest(candles []model.Candle) float64 {
	min := candles[0].Low
	for _, candle := range candles {
		if candle.Low < min {
			min = candle.Low
		}
	}
	return min
}
