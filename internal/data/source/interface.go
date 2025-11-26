package source

import (
	"signalstack/internal/core"
	"time"
)

type MarketDataSource interface {
	Name() string
	FetchOHLCV(
		symbol core.Symbol,
		timeframe core.Timeframe,
		start time.Time,
		end time.Time,
	) ([]core.Candle, error)
}