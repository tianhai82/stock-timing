package analyzer

import (
	"math"

	"github.com/tianhai82/stock-timing/model"
)

func AnalyzerCandles(candles []model.Candle, buyLimit, sellLimit float64) model.TradeAnalysis {
	periodMean := calcMean(candles, getNearClose)
	stdDev, maxDev := calcStdMaxDev(candles, periodMean, getNearClose)
	period := len(candles)
	currentDev := getClose(candles[period-1]) - periodMean

	analysis := model.TradeAnalysis{
		Period:        period,
		Mean:          periodMean,
		CurrentCandle: candles[period-1],
		StdDev:        stdDev,
		MaxDev:        maxDev,
		BuyLimitDev:   -buyLimit * (stdDev + maxDev),
		SellLimitDev:  sellLimit * (stdDev + maxDev),
		CurrentDev:    currentDev,
	}
	signal := getSignal(candles, analysis)
	analysis.Signal = signal
	return analysis
}

func getSignal(candles []model.Candle, analysis model.TradeAnalysis) model.TradeSignal {
	if len(candles) < 3 {
		return model.Hold
	}
	lastCandle := candles[len(candles)-1]
	if analysis.CurrentDev <= analysis.BuyLimitDev {
		if lastCandle.Close > lastCandle.Open {
			// if /*candleVal(lastCandle) >= candleVal(secondLastCandle) &&*/ lastCandle.Close >= lastCandle.Open && secondLastCandle.Close >= secondLastCandle.Open {
			return model.Buy
		}
	} else if analysis.CurrentDev >= analysis.SellLimitDev {
		if lastCandle.Close < lastCandle.Open {
			// if /*candleVal(lastCandle) <= candleVal(secondLastCandle) &&*/ lastCandle.Close <= lastCandle.Open && secondLastCandle.Close <= secondLastCandle.Open {
			return model.Sell
		}
	}
	return model.Hold
}

func getClose(candle model.Candle) float64 {
	return candle.Close
}

func getNearClose(candle model.Candle) float64 {
	return candle.Open*0.2 + candle.Close*0.8
}

func calcStdMaxDev(candles []model.Candle, mean float64, candleVal func(model.Candle) float64) (float64, float64) {
	maxDev := 0.0
	varianceSum := 0.0
	for _, candle := range candles {
		close := candleVal(candle)
		diff := math.Abs(close - mean)
		if diff > maxDev {
			maxDev = diff
		}
		varianceSum += math.Pow(diff, 2.0)
	}
	variance := varianceSum / float64(len(candles))
	stdDev := math.Sqrt(variance)
	return stdDev, maxDev

}

func calcMean(candles []model.Candle, candleVal func(model.Candle) float64) float64 {
	total := 0.0
	for _, c := range candles {
		total += candleVal(c)
	}
	return total / float64(len(candles))
}
