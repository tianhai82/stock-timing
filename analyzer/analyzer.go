package analyzer

import (
	"math"

	"github.com/tianhai82/stock-timing/model"
)

func AnalyzerCandles(candles []model.Candle) model.TradeAnalysis {
	periodMean := calcMean(candles, getMid)
	stdDev, maxDev := calcStdMaxDev(candles, periodMean, getMid)
	limitDev := (stdDev + maxDev) * 0.56
	period := len(candles)
	currentDev := getMid(candles[period-1]) - periodMean

	analysis := model.TradeAnalysis{
		Period:        period,
		Mean:          periodMean,
		CurrentCandle: candles[period-1],
		StdDev:        stdDev,
		MaxDev:        maxDev,
		LimitDev:      limitDev,
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
	if math.Abs(analysis.CurrentDev) >= analysis.LimitDev {
		lastCandle := candles[len(candles)-1]
		// secondLastCandle := candles[len(candles)-2]
		// thirdLastCandle := candles[len(candles)-3]
		if analysis.CurrentDev > 0 {
			if lastCandle.Close <= lastCandle.Open {
				// if /*candleVal(lastCandle) <= candleVal(secondLastCandle) &&*/ lastCandle.Close <= lastCandle.Open && secondLastCandle.Close <= secondLastCandle.Open {
				return model.Sell
			}
		} else {
			if lastCandle.Close >= lastCandle.Open {
				// if /*candleVal(lastCandle) >= candleVal(secondLastCandle) &&*/ lastCandle.Close >= lastCandle.Open && secondLastCandle.Close >= secondLastCandle.Open {
				return model.Buy
			}
		}
	}
	return model.Hold
}

func getClose(candle model.Candle) float64 {
	return candle.Close
}

func getMid(candle model.Candle) float64 {
	return (candle.Open + candle.Close) / 2
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
