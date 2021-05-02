package candle

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/tianhai82/stock-timing/firebase"
	"github.com/tianhai82/stock-timing/model"
	"github.com/tianhai82/stock-timing/sgx"
	"github.com/tianhai82/stock-timing/tda"
)

func RetrieveCandles(instrumentID int, period int) ([]model.Candle, error) {
	if instrumentID >= 0 {
		symbol, err := findSymbolFromInstrumentID(instrumentID)
		if err != nil {
			fmt.Println(err)
			return nil, errors.Wrap(err, "fail to find symbol from instrument ID")
		}
		return tda.RetrieveHistory(symbol, period)
	}
	symbol, err := findSymbolFromInstrumentID(instrumentID)
	if err != nil {
		fmt.Println(err)
		return nil, errors.Wrap(err, "fail to find symbol from instrument ID")
	}
	return sgx.RetrieveHistory(symbol, period)
}

func findSymbolFromInstrumentID(instrumentID int) (model.InstrumentDisplayData, error) {
	if firebase.StorageClient == nil {
		return model.InstrumentDisplayData{}, errors.New("cannot access storage")
	}
	if firebase.Instruments == nil {
		return model.InstrumentDisplayData{}, errors.New("cannot access instruments")
	}
	for _, ins := range firebase.Instruments {
		if ins.InstrumentID == instrumentID {
			return ins, nil
		}
	}
	return model.InstrumentDisplayData{}, errors.New("instrument ID not in global list")
}
