package engine

import (
	"math"
	"signalstack/internal/core"
	"signalstack/internal/strategy"
)

type EngineConfig struct {
	InitialCapital     float64
	SlippageBps        float64
	CommissionPerTrade float64
	AllInOnEntry       bool
}

type BacktestResult struct {
	Trades      []core.Trade
	EquityCurve []core.EquityPoint

	FinalEquity      float64
	MaxDrawdown      float64
	BuyAndHoldReturn float64
}

type Engine struct {
	cfg EngineConfig
}

func New(cfg EngineConfig) *Engine {
	return &Engine{
		cfg: cfg,
	}
}

func (e *Engine) Run(
	strat strategy.Strategy,
	candles []core.Candle,
) BacktestResult {
	if len(candles) == 0 {
		return BacktestResult{}
	}
	cash := e.cfg.InitialCapital
	var pos *core.Position
	trades := make([]core.Trade, 0, 16)
	equityCurve := make([]core.EquityPoint, 0, len(candles))

	for i, c := range candles {
		equity := cash
		if pos != nil {
			equity += pos.Quantity * c.Close
		}
		equityCurve = append(equityCurve, core.EquityPoint{
			Time:   c.Timestamp,
			Equity: equity,
		})
		if pos != nil {
			exited, trade := e.checkExitOnBar(pos, c)
			if exited {
				cash += trade.PnL + pos.Quantity*pos.AvgEntryPrice
				trades = append(trades, trade)
				pos = nil
				equityCurve = append(equityCurve, core.EquityPoint{
					Time:   c.Timestamp,
					Equity: equity,
				})
				equity = cash
			}
		}

		ctx := strategy.Context{
			Index:    i,
			Candle:   c,
			Position: pos,
			Equity:   cash,
		}
		decision := strat.OnBar(ctx)
		if decision.EnterLong {
			if pos == nil {
				newPos := e.openLongFromDecision(decision, c, equity)
				if newPos != nil {
					cash -= newPos.CostBasis
					pos = newPos
				}
			}
		}
		if decision.ExitLong {
			if pos != nil {
				trade := e.exitPosition(pos, c)
				cash += pos.Quantity*trade.ExitPrice - e.cfg.CommissionPerTrade
				trades = append(trades, trade)
				pos = nil
			}
		}
	}
	// Force close at last bar, if still open
	if len(candles) > 0 && pos != nil {
		last := candles[len(candles)-1]
		trade := e.exitPosition(pos, last)
		cash += trade.PnL + pos.Quantity*pos.AvgEntryPrice
		trades = append(trades, trade)
		pos = nil
		equityCurve = append(equityCurve, core.EquityPoint{
			Time:   last.Timestamp,
			Equity: cash,
		})
	}
	finalEquity := cash
	maxDD := computeMaxDrawdown(equityCurve)

	var buyAndHold float64
	if len(candles) > 0 {
		first := candles[0].Close
		last := candles[len(candles)-1].Close
		if first > 0 {
			buyAndHold = (last - first) / first
		}
	}

	return BacktestResult{
		Trades:           trades,
		EquityCurve:      equityCurve,
		FinalEquity:      finalEquity,
		MaxDrawdown:      maxDD,
		BuyAndHoldReturn: buyAndHold,
	}
}

// computeMaxDrawdown computes max peak-to-trough drawdown as a fraction (0.0–1.0)
// from an equity curve. If curve is empty or constant, returns 0.
func computeMaxDrawdown(equityCurve []core.EquityPoint) float64 {
	if len(equityCurve) == 0 {
		return 0.0
	}

	peak := equityCurve[0].Equity
	maxDD := 0.0
	for _, pt := range equityCurve {
		if pt.Equity > peak {
			peak = pt.Equity
			continue
		}
		if peak <= 0 {
			continue
		}
		dd := (peak - pt.Equity) / peak
		if dd > maxDD {
			maxDD = dd
		}
	}
	return maxDD
}

// openLongFromDecision opens a new long postition based on a buy order
// - equity is current total equity (cash + any open positions marked to market)
// - if order size == 0 we do all in sizing
func (e *Engine) openLongFromDecision(decision strategy.Decision, c core.Candle, equity float64) *core.Position {
	if !decision.EnterLong {
		return nil
	}
	if equity <= 0 {
		return nil
	}
	price := c.Close
	if e.cfg.SlippageBps != 0 {
		price *= (1 + e.cfg.SlippageBps/10000)
	}
	if price <= 0 {
		return nil
	}
	var size float64
	switch {
	case e.cfg.AllInOnEntry:
		size = math.Floor(equity / price)
	default:
		panic("invalid position sizing config")
	}
	if size <= 0 {
		return nil
	}
	return &core.Position{
		PositionSide:  core.PositionSideLong,
		Quantity:      size,
		AvgEntryPrice: price,
		StopLoss:      decision.StopLoss,
		TakeProfit:    decision.TakeProfit,
		EntryTime:     c.Timestamp,
		CostBasis:     price * size,
		CurrentPrice:  price,
	}
}

func (e *Engine) Config() EngineConfig {
	return e.cfg
}

// simple intraday exit check
func (e *Engine) checkExitOnBar(pos *core.Position, c core.Candle) (bool, core.Trade) {
	if pos == nil {
		return false, core.Trade{}
	}
	if pos.StopLoss == nil && pos.TakeProfit == nil {
		return false, core.Trade{}
	}
	if pos.StopLoss != nil && c.Low < *pos.StopLoss {
		trade := e.exitPosition(pos, c)
		return true, trade
	}
	if pos.TakeProfit != nil && c.High > *pos.TakeProfit {
		trade := e.exitPosition(pos, c)
		return true, trade
	}
	return false, core.Trade{}
}

func (e *Engine) exitPosition(
	pos *core.Position,
	c core.Candle,
) core.Trade {
	exitPrice := c.Close
	if e.cfg.SlippageBps != 0 {
		exitPrice *= 1.0 - e.cfg.SlippageBps/10000.0
	}

	grossPnL := (exitPrice - pos.AvgEntryPrice) * pos.Quantity
	netPnL := grossPnL - e.cfg.CommissionPerTrade

	var risk float64
	if pos.StopLoss != nil {
		risk = math.Abs(pos.AvgEntryPrice-*pos.StopLoss) * pos.Quantity
	}

	trade := core.Trade{
		Symbol:       pos.Symbol,
		PositionSide: pos.PositionSide,
		EntryTime:    pos.EntryTime,
		EntryPrice:   pos.AvgEntryPrice,
		ExitTime:     c.Timestamp,
		ExitPrice:    exitPrice,
		Size:         pos.Quantity,
		PnL:          netPnL,
		Risk:         risk,
	}
	return trade
}
