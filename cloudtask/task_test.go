package cloudtask

import (
	"fmt"
	"testing"

	"github.com/tianhai82/stock-timing/model"
	"github.com/tianhai82/stock-timing/telegram"
)

func TestSendMessage(t *testing.T) {
	highPrices := []model.EmailAnalysis{
		{
			BuyOrSell:             "Hold",
			InstrumentDisplayName: "TESLA",
			InstrumentSymbol:      "TSLA",
			InstrumentID:          -54,
			Period:                15,
			CurrentPrice:          15.6784,
			PricePercentile:       89.4444,
			BuyFreq:               16.0,
			SellFreq:              56.9,
		},
		{
			BuyOrSell:             "Hold",
			InstrumentDisplayName: "TESLA",
			InstrumentSymbol:      "TSLA",
			InstrumentID:          -54,
			Period:                15,
			CurrentPrice:          15.6784,
			PricePercentile:       89.4444,
			BuyFreq:               16.0,
			SellFreq:              56.9,
		},
	}
	highMd := formMarkdownMsg(highPrices, "HIGH")
	fmt.Println(highMd)
	telegram.SendMessage(highMd, telegram.CHAT_ID, telegram.TOKEN)
}
