package main

import (
	"encoding/json"
	"fmt"

	"github.com/tianhai82/stock-timing/analyzer"
	"github.com/tianhai82/stock-timing/etoro"
)

func main() {
	fmt.Println("Stock Timing starting")
	// instruments, err := etoro.RetrieveInstruments()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// json.NewEncoder(w io.Writer)

	candles, err := etoro.RetrieveCandle(2360, 60)
	if err != nil {
		fmt.Println(err)
		return
	}
	candles[len(candles)-2].Close = 14.49
	candles[len(candles)-3].Close = 14.50
	analysis := analyzer.AnalyzerCandles(candles)
	b, _ := json.Marshal(analysis)
	fmt.Println(string(b))
}
