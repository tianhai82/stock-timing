package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/tianhai82/stock-timing/httprequester"
)

// const period = "OneYearAgo"
const period = "SixMonthsAgo"
const positiveMonths = 9

const url = "https://www.etoro.com/sapi/rankings/rankings/?blocked=false&bonusonly=false&copyblock=false&copytradespctmax=20&dailyddmin=-5&gainmin=10&istestaccount=false&lastactivitymax=30&maxmonthlyriskscoremax=6&maxmonthlyriskscoremin=1&optin=true&page=1&pagesize=1000&period=%s&profitablemonthspctmin=100&sort=-gain&tradesmin=5&verified=true&weeklyddmin=-15"

type TraderInfo struct {
	CustomerID             int     `json:"CustomerId"`
	UserName               string  `json:"UserName"`
	HasAvatar              bool    `json:"HasAvatar"`
	IsSocialConnected      bool    `json:"IsSocialConnected"`
	IsTestAccount          bool    `json:"IsTestAccount"`
	DisplayFullName        bool    `json:"DisplayFullName"`
	BonusOnly              bool    `json:"BonusOnly"`
	Blocked                bool    `json:"Blocked"`
	Verified               bool    `json:"Verified"`
	PopularInvestor        bool    `json:"PopularInvestor"`
	CopyBlock              bool    `json:"CopyBlock"`
	IsFund                 bool    `json:"IsFund"`
	IsBronze               bool    `json:"IsBronze"`
	FundType               int     `json:"FundType"`
	Gain                   float64 `json:"Gain"`
	DailyGain              float64 `json:"DailyGain"`
	ThisWeekGain           float64 `json:"ThisWeekGain"`
	RiskScore              int     `json:"RiskScore"`
	MaxDailyRiskScore      int     `json:"MaxDailyRiskScore"`
	MaxMonthlyRiskScore    int     `json:"MaxMonthlyRiskScore"`
	Copiers                int     `json:"Copiers"`
	CopiedTrades           int     `json:"CopiedTrades"`
	CopyTradesPct          float64 `json:"CopyTradesPct"`
	CopyInvestmentPct      float64 `json:"CopyInvestmentPct"`
	BaseLineCopiers        int     `json:"BaseLineCopiers"`
	CopiersGain            float64 `json:"CopiersGain"`
	AUMTier                int     `json:"AUMTier"`
	AUMTierV2              int     `json:"AUMTierV2"`
	VirtualCopiers         int     `json:"VirtualCopiers"`
	Trades                 int     `json:"Trades"`
	WinRatio               float64 `json:"WinRatio"`
	DailyDD                float64 `json:"DailyDD"`
	WeeklyDD               float64 `json:"WeeklyDD"`
	ProfitableWeeksPct     float64 `json:"ProfitableWeeksPct"`
	ProfitableMonthsPct    float64 `json:"ProfitableMonthsPct"`
	Velocity               float64 `json:"Velocity"`
	Exposure               float64 `json:"Exposure"`
	AvgPosSize             float64 `json:"AvgPosSize"`
	OptimalCopyPosSize     float64 `json:"OptimalCopyPosSize"`
	HighLeveragePct        float64 `json:"HighLeveragePct"`
	MediumLeveragePct      float64 `json:"MediumLeveragePct"`
	LowLeveragePct         float64 `json:"LowLeveragePct"`
	PeakToValley           float64 `json:"PeakToValley"`
	PeakToValleyStart      string  `json:"PeakToValleyStart"`
	PeakToValleyEnd        string  `json:"PeakToValleyEnd"`
	LongPosPct             float64 `json:"LongPosPct"`
	TopTradedInstrumentID  int     `json:"TopTradedInstrumentId"`
	TopTradedAssetClassID  int     `json:"TopTradedAssetClassId"`
	ActiveWeeks            int     `json:"ActiveWeeks"`
	FirstActivity          string  `json:"FirstActivity"`
	LastActivity           string  `json:"LastActivity"`
	ActiveWeeksPct         float64 `json:"ActiveWeeksPct"`
	WeeksSinceRegistration int     `json:"WeeksSinceRegistration"`
	Country                string  `json:"Country"`
	AffiliateID            int     `json:"AffiliateId"`
	InstrumentPct          float64 `json:"InstrumentPct"`
	FullName               string  `json:"FullName,omitempty"`
}

type filterResponse struct {
	Status    string       `json:"Status"`
	TotalRows int          `json:"TotalRows"`
	Items     []TraderInfo `json:"Items"`
}

type TraderGain struct {
	CustomerID int `json:"customerId"`
	Monthly    []struct {
		Start        time.Time `json:"start"`
		Gain         float64   `json:"gain"`
		IsSimulation bool      `json:"isSimulation"`
	} `json:"monthly"`
	Yearly []struct {
		Start        time.Time `json:"start"`
		Gain         float64   `json:"gain"`
		IsSimulation bool      `json:"isSimulation"`
	} `json:"yearly"`
}

func GetEligibleTraders() ([]TraderInfo, error) {
	var filterResp filterResponse
	urlPeriod := fmt.Sprintf(url, period)
	err := httprequester.MakeGetRequest(urlPeriod, &filterResp)
	if err != nil {
		err = httprequester.MakeGetRequest(urlPeriod, &filterResp)
		if err != nil {
			return nil, fmt.Errorf("fail to query etoro. %v", err)
		}
	}
	return filterResp.Items, nil
}

func GetTradersFromFile() ([]TraderInfo, error) {
	f, err := os.Open("eligible_traders.json")
	if err != nil {
		return nil, err
	}
	var resp filterResponse
	err = json.NewDecoder(f).Decode(&resp)
	return resp.Items, err
}

func main() {
	fmt.Println("test")
	// traderInfos, err := GetEligibleTraders()
	traderInfos, err := GetTradersFromFile()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(len(traderInfos))
	okTraders := make([]TraderInfo, 0, len(traderInfos))

	for i := 0; i < len(traderInfos); i++ {
		trader := traderInfos[i]
		custId := trader.CustomerID
		var traderGain TraderGain
		errGain := httprequester.MakeGetRequest(fmt.Sprintf("https://www.etoro.com/sapi/userstats/gain/cid/%d/history?IncludeSimulation=false", custId), &traderGain)
		if errGain != nil {
			fmt.Println(errGain)
			fmt.Println("pause 15 min")
			time.Sleep(15 * time.Minute)
			i = i - 1
			continue
		}
		fmt.Printf("Processing trader: %s\n", trader.UserName)
		ok := true
		if len(traderGain.Monthly) < (positiveMonths + 1) {
			fmt.Printf("  --  %s is too new\n", trader.UserName)
			time.Sleep(1000 * time.Millisecond)
			continue
		}
		if trader.CopyInvestmentPct > 20 {
			fmt.Printf("  --  %s is is also a copier\n", trader.UserName)
			time.Sleep(1000 * time.Millisecond)
			continue
		}
		totalGain := 0.0
		totalProcessed := 0
		for i := len(traderGain.Monthly) - 1; totalProcessed < positiveMonths; i-- {
			gain := traderGain.Monthly[i]
			if gain.Gain > 0 {
				totalGain += gain.Gain
				totalProcessed++
				continue
			} else {
				fmt.Printf("  --  %s month: %s is not positive\n", trader.UserName, gain.Start.Month().String())
				ok = false
				break
			}
		}
		if ok {
			totalProcessed = 0
			for i := len(traderGain.Monthly) - 1; totalProcessed < positiveMonths; i-- {
				gain := traderGain.Monthly[i]
				if gain.Gain > (totalGain / (float64(positiveMonths) / 3.5)) {
					ok = false
					fmt.Printf("  --  %s month %s gain is too high at %.1f. Total %d month gain is %.1f\n",
						trader.UserName, gain.Start.Month().String(), gain.Gain, positiveMonths, totalGain)
					break
				}
				totalProcessed++
			}
		}
		if ok {
			okTraders = append(okTraders, trader)
			fmt.Printf("  ::ADDING TRADER: %s. Total Count: %d\n", trader.UserName, len(okTraders))
		}
		if len(okTraders) >= 3 {
			break
		}
		time.Sleep(5000 * time.Millisecond)
	}
	fmt.Println("total ok trader", len(okTraders))
	timeStr := time.Now().Format("2006-01-02-15-04")
	f, err := os.Create(fmt.Sprintf("ok_trader_%s.json", timeStr))
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(f).Encode(okTraders)
}
