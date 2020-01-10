// +build ignore

package main

import (
	"fmt"

	"github.com/tianhai82/stock-timing/etoro"
	"github.com/tianhai82/stock-timing/model"
	"github.com/tianhai82/stock-timing/yahoo"
)

type Empty struct{}

func main() {

	// retrieveEtoroSymbols()
	testYahooHistory()

}

func testYahooHistory(){
	yahoo.RetrieveHistory("C06.SI", 10)
}

func retrieveEtoroSymbols(){
	
	fmt.Println("pull etoro instruments")
	instruments, err := etoro.RetrieveInstruments()
	if err!=nil{
		fmt.Println(err)
		return
	}
	filteredInstruments := make([]model.InstrumentDisplayData,0,len(instruments))
	for _, ins := range instruments {
		if ins.InstrumentTypeID == 5 ||ins.InstrumentTypeID == 6 {
			filteredInstruments = append(filteredInstruments, ins)
		}
	}
	fmt.Println("total", len(instruments), len(filteredInstruments))
}