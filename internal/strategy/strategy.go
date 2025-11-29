package strategy

import "signalstack/internal/core"

type Decision struct {
	EnterLong  bool
	ExitLong   bool
	StopLoss   *float64
	TakeProfit *float64
}

type Context struct {
	Index    int
	Candle   core.Candle
	Position *core.Position
	Equity   float64
}

type Strategy interface {
	OnBar(ctx Context) Decision
	Name() string
}
