package storage

import (
	"time"

	"signalstack/internal/core"
)

type BarStore interface {
	Store(symbol core.Symbol, timeframe core.Timeframe, candles []core.Candle) error
	Fetch(symbol core.Symbol, timeframe core.Timeframe, start *time.Time) ([]core.Candle, error)
	LatestTimestamp(symbol core.Symbol, timeframe core.Timeframe) (*time.Time, error)
}
