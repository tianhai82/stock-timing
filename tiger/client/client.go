package client

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tianhai82/stock-timing/model"
	"github.com/tianhai82/stock-timing/tiger"
	tigerModel "github.com/tianhai82/stock-timing/tiger/model"
	"github.com/tianhai82/stock-timing/tiger/openapi"
)

var (
	privateKey = ""
	tigerID    = ""
	account    = ""
)

func init() {
	privateKey = os.Getenv("TIGER_KEY")
	if privateKey == "" {
		panic("cannot retrieve TIGER_KEY")
	}
	tigerID = os.Getenv("TIGER_ID")
	if tigerID == "" {
		panic("cannot retrieve TIGER_ID")
	}
	account = os.Getenv("TIGER_ACCOUNT")
	if account == "" {
		panic("cannot retrieve TIGER_ACCOUNT")
	}
}

func RetrieveHistory(instrument model.InstrumentDisplayData, period int, freq string) ([]model.Candle, error) {
	client := tiger.NewTigerOpenClient(tiger.NewTigerOpenClientConfig(false, false, privateKey, tigerID, account), log.Default())
	req := &openapi.OpenApiRequest{
		Method: tigerModel.GET_QUOTE_PERMISSION,
	}
	_, err := client.Execute(req, "")
	if err != nil {
		return nil, fmt.Errorf("tigerClient.RetrieveHistory GET_QUOTE_PERMISSION failed: %w", err)
	}

	req = &openapi.OpenApiRequest{
		Method: tigerModel.KLINE,
	}
	params := map[string]interface{}{}
	params["symbols"] = []string{instrument.SymbolFull}
	params["period"] = freq
	params["begin_time"] = -1
	params["end_time"] = -1
	params["right"] = "br"
	params["limit"] = period
	params["lang"] = tiger.LANGUAGE
	req.BizModel = params
	resp, err := client.Execute(req, "")
	if err != nil {
		return nil, fmt.Errorf("tigerClient.RetrieveHistory KLINE failed: %w", err)
	}
	data := resp["data"].([]interface{})[0]
	respPeriod := data.(map[string]interface{})["period"].(string)
	respSymbol := data.(map[string]interface{})["symbol"].(string)
	if respPeriod != freq {
		return nil, fmt.Errorf("wrong freq returned. requested:%s, gotten:%s", freq, respPeriod)
	}
	if respSymbol != instrument.SymbolFull {
		return nil, fmt.Errorf("wrong symbol returned. requested:%s, gotten:%s", instrument.SymbolFull, respSymbol)
	}
	items := data.(map[string]interface{})["items"].([]interface{})
	candles := make([]model.Candle, 0, len(items))
	for _, item := range items {
		mapItem := item.(map[string]interface{})
		close := mapItem["close"].(float64)
		high := mapItem["high"].(float64)
		low := mapItem["low"].(float64)
		open := mapItem["open"].(float64)
		candleTime := mapItem["time"].(float64)

		candle := model.Candle{
			InstrumentID: instrument.InstrumentID,
			Open:         open,
			High:         high,
			Low:          low,
			Close:        close,
			FromDate:     time.UnixMilli(int64(candleTime)),
		}
		candles = append(candles, candle)
	}
	return candles, nil
}
