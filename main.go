package main

import (
	"fmt"

	"github.com/tianhai82/stock-timing/etoro"
)

func main() {
	fmt.Println("Stock Timing starting")
	_, err := etoro.RetrieveInstruments()
	if err != nil {
		fmt.Println(err)
		return
	}
	// for _, in := range instruments {
	// 	if in.InstrumentTypeID == 5 {
	// 		fmt.Printf("Type: %d. ID: %d. Name: %s\n", in.InstrumentTypeID, in.InstrumentID, in.InstrumentDisplayName)
	// 	}
	// }
	candles, err := etoro.RetrieveCandle(2360, 60)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, candle := range candles {
		fmt.Println(candle.FromDate, candle.Close)
	}
}
