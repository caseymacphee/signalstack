package core

import "time"

type Candle struct {
	Timestamp time.Time
	Open float64
	High float64
	Low float64
	Close float64
	Volume int
}

type Symbol string

type Timeframe string

const (
	TimeframeDaily Timeframe = "1d"
	TimeframeHourly Timeframe = "1h"
)