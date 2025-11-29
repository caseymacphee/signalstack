package service

import (
	"math"

	"signalstack/internal/core"
	"signalstack/internal/engine"
	"signalstack/pkg/api"
)

func computeMetrics(result engine.BacktestResult, initialCapital float64) api.BacktestMetrics {
	final := result.FinalEquity
	cagr := computeCAGR(initialCapital, final, result.EquityCurve)
	winRate, profitFactor, avgRMultiple := computeTradeStats(result.Trades)

	var totalReturn float64
	if initialCapital > 0 {
		totalReturn = (final - initialCapital) / initialCapital
	}

	return api.BacktestMetrics{
		CAGR:             cagr,
		TotalReturn:      totalReturn,
		MaxDrawdown:      result.MaxDrawdown,
		WinRate:          winRate,
		ProfitFactor:     profitFactor,
		FinalEquity:      final,
		AverageRMultiple: avgRMultiple,
		BuyAndHoldReturn: result.BuyAndHoldReturn,
	}
}

func computeCAGR(initialCapital, finalEquity float64, equityCurve []core.EquityPoint) float64 {
	if initialCapital <= 0 || finalEquity <= 0 || len(equityCurve) < 2 {
		return 0
	}

	start := equityCurve[0].Time
	end := equityCurve[len(equityCurve)-1].Time
	years := end.Sub(start).Hours() / 24 / 365
	if years <= 0 {
		return 0
	}
	return math.Pow(finalEquity/initialCapital, 1/years) - 1
}

func computeTradeStats(trades []core.Trade) (float64, float64, float64) {
	if len(trades) == 0 {
		return 0, 0, 0
	}
	var wins, total int
	var profit, grossLoss, totalR float64
	var rCount int

	for _, trade := range trades {
		total++
		if trade.PnL > 0 {
			wins++
			profit += trade.PnL
		} else {
			grossLoss += -trade.PnL
		}

		if trade.Risk > 0 {
			totalR += trade.PnL / trade.Risk
			rCount++
		}
	}

	winRate := float64(wins) / float64(total)

	var profitFactor float64
	if grossLoss > 0 {
		profitFactor = profit / grossLoss
	}

	var avgRMultiple float64
	if rCount > 0 {
		avgRMultiple = totalR / float64(rCount)
	}

	return winRate, profitFactor, avgRMultiple
}
