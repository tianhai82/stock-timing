package candle

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/tianhai82/stock-timing/etoro"
	"github.com/tianhai82/stock-timing/firebase"
	"github.com/tianhai82/stock-timing/model"
	"github.com/tianhai82/stock-timing/sgx"
)

func RetrieveCandles(instrumentID int, period int) ([]model.Candle, error) {
	if instrumentID >= 0 {
		return etoro.RetrieveCandle(instrumentID, period)
	}
	symbol, err := findSymbolFromInstrumentID(instrumentID)
	if err != nil {
		fmt.Println(err)
		return nil, errors.Wrap(err, "fail to find symbol from instrument ID")
	}
	return sgx.RetrieveHistory(strings.Replace(symbol, ".SI", "", -1), period)
}

func findSymbolFromInstrumentID(instrumentID int) (string, error) {
	if firebase.StorageClient == nil {
		return "", errors.New("cannot access storage")
	}
	if firebase.Instruments == nil {
		return "", errors.New("cannot access instruments")
	}
	for _, ins := range firebase.Instruments {
		if ins.InstrumentID == instrumentID {
			return ins.SymbolFull, nil
		}
	}
	return "", errors.New("instrument ID not in global list")
}
