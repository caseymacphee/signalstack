package core

import (
	"time"
)


type Position struct {
	Symbol Symbol
	Side Side
	Quantity int
	AvgEntryPrice float64
	CostBasis float64
	CurrentPrice float64
	UnrealizedProfit float64
	UnrealizedProfitPercent float64
	EntryTime time.Time
	LastExitTime time.Time
}

type PositionStatus string

const (
	PositionStatusOpen PositionStatus = "open"
	PositionStatusClosed PositionStatus = "closed"
)

