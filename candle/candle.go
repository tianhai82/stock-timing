package candle

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/tianhai82/stock-timing/aastocks"
	"github.com/tianhai82/stock-timing/coinmarketcap"
	"github.com/tianhai82/stock-timing/firebase"
	"github.com/tianhai82/stock-timing/model"
	"github.com/tianhai82/stock-timing/sgx"
	tigerClient "github.com/tianhai82/stock-timing/tiger/client"
)

func RetrieveCandles(instrumentID int, period int) ([]model.Candle, error) {
	if instrumentID >= 0 && instrumentID < 1000001 {
		symbol, err := findSymbolFromInstrumentID(instrumentID)
		if err != nil {
			fmt.Println(err)
			return nil, errors.Wrap(err, "fail to find symbol from instrument ID")
		}
		return tigerClient.RetrieveHistory(symbol, period, "day")
	}

	if instrumentID >= 1000001 && instrumentID < 2000001 {
		symbol, err := findSymbolFromInstrumentID(instrumentID)
		if err != nil {
			fmt.Println(err)
			return nil, errors.Wrap(err, "fail to find symbol from instrument ID")
		}
		return aastocks.RetrieveHistory(symbol, period)
	}

	if instrumentID >= 2000001 {
		symbol, err := findSymbolFromInstrumentID(instrumentID)
		if err != nil {
			fmt.Println(err)
			return nil, errors.Wrap(err, "fail to find symbol from instrument ID")
		}
		return coinmarketcap.RetrieveHistory(symbol, period)
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
