package core

import (
	"time"
)


type Position struct {
	Symbol Symbol
	PositionSide PositionSide
	Quantity float64
	AvgEntryPrice float64
	CostBasis float64
	CurrentPrice float64
	UnrealizedProfit float64
	UnrealizedProfitPercent float64
	EntryTime time.Time
	LastExitTime time.Time
	StopLoss *float64
	TakeProfit *float64
}

type PositionStatus string

const (
	PositionStatusOpen PositionStatus = "open"
	PositionStatusClosed PositionStatus = "closed"
)

