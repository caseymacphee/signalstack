package service

import (
	"fmt"

	"signalstack/internal/core"
	"signalstack/internal/data/storage"
	"signalstack/internal/engine"
	"signalstack/internal/strategy"
	"signalstack/pkg/api"
)

type BacktestService struct {
	Store            storage.BarStore
	Engine           *engine.Engine
	StrategyRegistry *strategy.Registry
}

func NewBacktestService(store storage.BarStore, registry *strategy.Registry, engine *engine.Engine) *BacktestService {
	return &BacktestService{
		Store:            store,
		Engine:           engine,
		StrategyRegistry: registry,
	}
}

func (service *BacktestService) Run(req api.BacktestRequest) (api.BacktestResponse, error) {
	// 1 load candles
	candles, err := service.Store.Fetch(core.Symbol(req.Symbol), core.Timeframe(req.Timeframe), &req.StartDate)
	if err != nil {
		return api.BacktestResponse{}, err
	}
	if len(candles) == 0 {
		return api.BacktestResponse{}, fmt.Errorf("no candles found for symbol %s and timeframe %s", req.Symbol, req.Timeframe)
	}

	// 2 create strategy
	strat, err := service.StrategyRegistry.New(req.Strategy, req.Params)
	if err != nil {
		return api.BacktestResponse{}, err
	}

	// 3 run backtest
	result := service.Engine.Run(strat, candles)

	// 4 compute metrics
	metrics := computeMetrics(result, service.Engine.Config().InitialCapital)

	// 5) convert trades to external format
	trades := make([]api.TradeDTO, 0, len(result.Trades))
	for _, t := range result.Trades {
		var rMultiple *float64
		if t.Risk > 0 {
			val := t.PnL / t.Risk
			rMultiple = &val
		}

		side := "BUY"
		if t.PositionSide == core.PositionSideShort {
			side = "SELL"
		}

		trades = append(trades, api.TradeDTO{
			Side:      side,
			Position:  string(t.PositionSide),
			EntryTime: t.EntryTime,
			Entry:     t.EntryPrice,
			ExitTime:  t.ExitTime,
			Exit:      t.ExitPrice,
			Size:      t.Size,
			PnL:       t.PnL,
			RMultiple: rMultiple,
		})
	}
	return api.BacktestResponse{
		JobID:    req.JobID,
		Symbol:   req.Symbol,
		Strategy: req.Strategy,
		Metrics:  metrics,
		Trades:   trades,
	}, nil
}
