package model

import "time"

type EtoroInstruments struct {
	InstrumentDisplayDatas []InstrumentDisplayData `json:"InstrumentDisplayDatas"`
}
type InstrumentDisplayData struct {
	InstrumentID                int    `json:"InstrumentID"`
	InstrumentDisplayName       string `json:"InstrumentDisplayName"`
	InstrumentTypeID            int    `json:"InstrumentTypeID"`
	ExchangeID                  int    `json:"ExchangeID"`
	SymbolFull                  string `json:"SymbolFull"`
	StocksIndustryID            int    `json:"StocksIndustryID,omitempty"`
	InstrumentTypeSubCategoryID int    `json:"InstrumentTypeSubCategoryID,omitempty"`
}
type EtoroCandle struct {
	Interval string        `json:"Interval"`
	Candles  []CandleOuter `json:"Candles"`
}
type Candle struct {
	InstrumentID int       `json:"InstrumentID"`
	FromDate     time.Time `json:"FromDate"`
	Open         float64   `json:"Open"`
	High         float64   `json:"High"`
	Low          float64   `json:"Low"`
	Close        float64   `json:"Close"`
}
type CandleOuter struct {
	InstrumentID int      `json:"InstrumentId"`
	Candles      []Candle `json:"Candles"`
	RangeOpen    float64  `json:"RangeOpen"`
	RangeClose   float64  `json:"RangeClose"`
	RangeHigh    float64  `json:"RangeHigh"`
	RangeLow     float64  `json:"RangeLow"`
}
type TradeSignal int

const (
	Hold = iota
	Buy
	Sell
)

type TradeAnalysis struct {
	Period        int
	Mean          float64
	StdDev        float64
	MaxDev        float64
	LimitDev      float64
	CurrentDev    float64
	CurrentCandle Candle
	Signal        TradeSignal
}

type TradeAdvice struct {
	Date   time.Time
	Price  float64
	Signal TradeSignal
}

type UserAccount struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}
type Stock struct {
	Symbol                string `json:"symbol"`
	InstrumentID          int    `json:"instrumentID"`
	InstrumentDisplayName string `json:"instrumentDisplayName"`
}
type StockSubscription struct {
	Stock         `json:",inline"`
	UserID        string    `json:"email"`
	LastUpdatedAt time.Time `json:"lastUpdatedAt"`
}
type UserSubscription struct {
	UserID        string
	InstrumentIDs []int
}
