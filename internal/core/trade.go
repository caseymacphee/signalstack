package core

import "time"

type Trade struct {
	PositionSide PositionSide
	EntryTime time.Time
	EntryPrice float64
	ExitTime time.Time
	ExitPrice float64
	Size float64
	PnL float64
}

type EquityPoint struct {
    Time   time.Time
    Equity float64
}