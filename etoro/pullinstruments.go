// +build ignore

package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/tianhai82/stock-timing/etoro"
	"github.com/tianhai82/stock-timing/model"
	"github.com/tianhai82/stock-timing/yahoo"
)

type Empty struct{}

func main() {

	instruments, err := retrieveEtoroSymbols()
	if err != nil {
		fmt.Println(err)
		return
	}
	// sgxInstruments, err := retrieveYahooSymbols()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Printf("etoro: %d, yahoo: %d\n", len(instruments), len(sgxInstruments))
	// instruments = append(instruments, sgxInstruments...)
	b, err := json.Marshal(instruments)
	newFile, err := os.Create("out2.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	newFile.Write(b)
	newFile.Close()
	//testYahooHistory()

}

func retrieveYahooSymbols() ([]model.InstrumentDisplayData, error) {
	f, err := os.Open("myData.csv")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	line := -1
	instruments := make([]model.InstrumentDisplayData, 0, 800)
	for {
		line++
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if line == 0 {
			continue
		}
		typ := record[0]
		name := record[1]
		symbol := record[2]
		id := record[3]
		idInt, errConv := strconv.Atoi(id)
		if errConv != nil {
			return nil, errConv
		}
		ins := model.InstrumentDisplayData{
			Type:                  typ,
			InstrumentID:          idInt,
			InstrumentDisplayName: name,
			SymbolFull:            symbol,
		}
		instruments = append(instruments, ins)
	}
	return instruments, nil
}

func testYahooHistory() {
	candles, err := yahoo.RetrieveHistory("C06.SI", 10)
	if err != nil {
		fmt.Println(err)
	}
	b, _ := json.Marshal(candles)
	fmt.Println(string(b))
}

func retrieveEtoroSymbols() ([]model.InstrumentDisplayData, error) {

	fmt.Println("pull etoro instruments")
	instruments, err := etoro.RetrieveInstruments()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	filteredInstruments := make([]model.InstrumentDisplayData, 0, len(instruments))
	for _, ins := range instruments {
		if ins.InstrumentTypeID == 5 || ins.InstrumentTypeID == 6 || ins.InstrumentTypeID == 1 {
			filteredInstruments = append(filteredInstruments, ins)
		}
	}
	return filteredInstruments, nil
}
