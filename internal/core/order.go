package core

import "time"

type Order struct {
	Symbol Symbol
	Side Side
	OrderType OrderType
	Quantity int
	FilledQuantity int
	FilledAvgPrice float64
	LimitPrice float64
	StopPrice float64
	TimeInForce TimeInForce
	SubmittedAt time.Time
	UpdatedAt time.Time
	FilledAt time.Time
	Status OrderStatus
}

type Side string

const (
	SideBuy Side = "buy"
	SideSell Side = "sell"
)

type OrderType string

//  market, limit, stop, stop_limit, trailing_stop
const (
	OrderTypeMarket OrderType = "market"
	OrderTypeLimit OrderType = "limit"
	OrderTypeStop OrderType = "stop"
	OrderTypeStopLimit OrderType = "stop_limit"
	OrderTypeTrailingStop OrderType = "trailing_stop"
)

type TimeInForce string

const (
	TimeInForceGoodTillCancelled TimeInForce = "gtc"
	TimeInForceImmediateOrCancel TimeInForce = "ioc"
	TimeInForceFillOrKill TimeInForce = "fok"
)

type OrderStatus string
// new, filled, partially_filled, canceled, expired, rejected
const (
	OrderStatusNew OrderStatus = "new"
	OrderStatusFilled OrderStatus = "filled"
	OrderStatusPartiallyFilled OrderStatus = "partially_filled"
	OrderStatusCanceled OrderStatus = "canceled"
	OrderStatusExpired OrderStatus = "expired"
	OrderStatusRejected OrderStatus = "rejected"
)


