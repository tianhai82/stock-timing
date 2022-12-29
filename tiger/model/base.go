package model

type BaseParams struct {
	Version string `json:"version"`
	Lang    string `json:"lang"`
}

type MultipleQuoteParams struct {
	BaseParams

	Symbols            []string `json:"symbols"`
	Symbol             string   `json:"symbol"`
	IncludeHourTrading *bool    `json:"includeHourTrading"`
	IncludeAskBid      bool     `json:"includeAskBid"`
	Right              string   `json:"right"`
	Period             string   `json:"period"`
	BeginTime          string   `json:"beginTime"`
	EndTime            string   `json:"endTime"`
	Limit              int      `json:"limit"`
	BeginIndex         int      `json:"beginIndex"`
	EndIndex           int      `json:"endIndex"`
	Date               string   `json:"date"`
	PageToken          string   `json:"pageToken"`
	TradeSession       string   `json:"tradeSession"`
}
