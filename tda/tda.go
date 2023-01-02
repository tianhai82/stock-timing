package tda

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/tianhai82/go-tdameritrade"
	"github.com/tianhai82/oauth2"
	"github.com/tianhai82/stock-timing/firebase"
	"github.com/tianhai82/stock-timing/model"
)

type FirebaseStore struct{}

var cachedState string

var client *tdameritrade.Client

func (s *FirebaseStore) StoreToken(token *oauth2.Token) error {
	doc := firebase.FirestoreClient.Collection("tda").Doc("token")
	doc.Set(context.Background(), token)
	return nil
}

func (s FirebaseStore) GetToken() (*oauth2.Token, error) {
	doc, err := firebase.FirestoreClient.Collection("tda").Doc("token").Get(context.Background())
	if err != nil {
		return nil, err
	}
	var token *oauth2.Token
	err = doc.DataTo(&token)
	return token, err
}

func (s FirebaseStore) StoreState(state string) error {
	cachedState = state
	return nil
}

func (s FirebaseStore) GetState() (string, error) {
	if cachedState == "" {
		return "", fmt.Errorf("state not found")
	}
	return cachedState, nil
}

// func init() {
// 	authenticator := tdameritrade.NewAuthenticator(
// 		&FirebaseStore{},
// 		oauth2.Config{
// 			ClientID: os.Getenv("CLIENT_ID"),
// 			Endpoint: oauth2.Endpoint{
// 				TokenURL: "https://api.tdameritrade.com/v1/oauth2/token",
// 				AuthURL:  "https://auth.tdameritrade.com/auth",
// 			},
// 			RedirectURL: "http://localhost:8080/callback",
// 		},
// 	)
// 	var errClient error
// 	client, errClient = authenticator.AuthenticatedClient(context.Background(), nil)
// 	if errClient != nil {
// 		fmt.Println("fail to authenticate", errClient)
// 		return
// 	}
// }

func RetrieveHistory(instrument model.InstrumentDisplayData, period int) ([]model.Candle, error) {
	noOfMonths := period / 20
	noOfMonths++
	noOfYears := 0
	if noOfMonths > 3 && noOfMonths < 6 {
		noOfMonths = 6
	} else if noOfMonths > 6 {
		noOfYears = period / 251
		noOfYears++
	}
	extendedHours := false
	endDate := time.Now().Add(24*time.Hour).Unix() * 1000
	opt := &tdameritrade.PriceHistoryOptions{
		PeriodType:            "month",
		Period:                noOfMonths,
		FrequencyType:         "daily",
		EndDate:               endDate,
		NeedExtendedHoursData: &extendedHours,
	}
	if noOfYears > 0 {
		opt = &tdameritrade.PriceHistoryOptions{
			PeriodType:            "year",
			Period:                noOfYears,
			FrequencyType:         "daily",
			EndDate:               endDate,
			NeedExtendedHoursData: &extendedHours,
		}
	}
	priceHistory, _, err := client.PriceHistory.PriceHistory(context.Background(), instrument.SymbolFull, opt)
	if err != nil {
		return nil, err
	}
	if len(priceHistory.Candles) > period {
		priceHistory.Candles = priceHistory.Candles[len(priceHistory.Candles)-period:]
	}
	return convertPriceHistoryToCandles(priceHistory, instrument.InstrumentID)
}
func RetrieveWeeklyHistory(instrument model.InstrumentDisplayData, period int) ([]model.Candle, error) {
	noOfMonths := period / 4
	noOfMonths++
	noOfYears := 0
	if noOfMonths > 3 && noOfMonths < 6 {
		noOfMonths = 6
	} else if noOfMonths > 6 {
		noOfYears = period / 52
		noOfYears++
	}
	extendedHours := false
	endDate := time.Now().Add(24*time.Hour).Unix() * 1000
	opt := &tdameritrade.PriceHistoryOptions{
		PeriodType:            "month",
		Period:                noOfMonths,
		FrequencyType:         "weekly",
		EndDate:               endDate,
		NeedExtendedHoursData: &extendedHours,
	}
	if noOfYears > 0 {
		opt = &tdameritrade.PriceHistoryOptions{
			PeriodType:            "year",
			Period:                noOfYears,
			FrequencyType:         "weekly",
			EndDate:               endDate,
			NeedExtendedHoursData: &extendedHours,
		}
	}
	priceHistory, _, err := client.PriceHistory.PriceHistory(context.Background(), instrument.SymbolFull, opt)
	if err != nil {
		return nil, err
	}
	if len(priceHistory.Candles) > period {
		priceHistory.Candles = priceHistory.Candles[len(priceHistory.Candles)-period:]
	}
	return convertPriceHistoryToCandles(priceHistory, instrument.InstrumentID)
}

func convertPriceHistoryToCandles(priceHistory *tdameritrade.PriceHistory, instrumentID int) ([]model.Candle, error) {
	if priceHistory == nil {
		return nil, fmt.Errorf("no candles found")
	}
	if priceHistory.Empty {
		return nil, fmt.Errorf("no candles found")
	}
	candles := make([]model.Candle, len(priceHistory.Candles))
	for i, tdaCandle := range priceHistory.Candles {
		c := model.Candle{
			InstrumentID: instrumentID,
			FromDate:     time.Unix(int64(tdaCandle.Datetime)/1000, 0),
			Open:         tdaCandle.Open,
			High:         tdaCandle.High,
			Low:          tdaCandle.Low,
			Close:        tdaCandle.Close,
		}
		candles[i] = c
	}
	return candles, nil
}

func RetrieveOptionChain(symbol, contractType, optionRange, fromDate, toDate string) (*tdameritrade.Chains, error) {
	urlValue := url.Values{}
	urlValue.Add("symbol", symbol)
	urlValue.Add("contractType", contractType)
	urlValue.Add("range", optionRange)
	urlValue.Add("fromDate", fromDate)
	urlValue.Add("toDate", toDate)
	urlValue.Add("includeQuotes", "TRUE")
	chains, _, err := client.Chains.GetChains(context.Background(), urlValue)
	return chains, err
}
