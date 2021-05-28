package aastocks

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tianhai82/stock-timing/httprequester"
	"github.com/tianhai82/stock-timing/model"
)

func RetrieveHistory(instrument model.InstrumentDisplayData, period int) ([]model.Candle, error) {
	out := ""
	err := httprequester.MakeGetStringRequest(
		fmt.Sprintf("http://chartdata1.internet.aastocks.com/servlet/iDataServlet/getdaily?id=%s&type=24&market=1&level=1&period=56&encoding=utf8", instrument.SymbolFull),
		&out)
	if err != nil {
		return nil, err
	}
	a := strings.Split(out, "|")
	a = a[2 : len(a)-1]

	candles := make([]model.Candle, 0, len(a))
	for _, i := range a {
		if strings.HasPrefix(i, "!") {
			break
		}
		split := strings.Split(i, ";")
		open, err := parseInt(split[1])
		if err != nil {
			return nil, err
		}
		high, err := parseInt(split[2])
		if err != nil {
			return nil, err
		}
		low, err := parseInt(split[3])
		if err != nil {
			return nil, err
		}
		close, err := parseInt(split[4])
		if err != nil {
			return nil, err
		}
		fromDate, err := parseDate(split[0])
		if err != nil {
			return nil, err
		}
		candle := model.Candle{
			InstrumentID: instrument.InstrumentID,
			Open:         open,
			High:         high,
			Low:          low,
			Close:        close,
			FromDate:     fromDate,
		}
		candles = append(candles, candle)
	}
	if len(candles) > period {
		candles = candles[len(candles)-period:]
	}
	return candles, nil
}
func parseDate(s string) (time.Time, error) {
	return time.Parse("01/02/2006", s)
}

func parseInt(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}
