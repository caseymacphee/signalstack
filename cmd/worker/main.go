package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"signalstack/internal/data/storage"
	"signalstack/internal/engine"
	"signalstack/internal/service"
	"signalstack/internal/strategy"
	"signalstack/pkg/api"
)

func main() {
	// Wire dependencies
	store := &storage.CSVStore{
		RootDir: "./data/raw/yahoo",
	}

	reg := strategy.NewRegistry()
	strategy.RegisterBuiltins(reg)

	eng := engine.New(engine.EngineConfig{
		InitialCapital:     10000,
		SlippageBps:        5,
		CommissionPerTrade: 0,
		AllInOnEntry:       true,
	})

	svc := service.NewBacktestService(store, reg, eng)

	// For now: read a BacktestRequest from stdin as JSON
	var req api.BacktestRequest
	if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
		log.Fatalf("decode request: %v", err)
	}

	// If you want, default dates if zero
	if req.StartDate.IsZero() {
		req.StartDate = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	if req.EndDate.IsZero() {
		req.EndDate = time.Now().UTC()
	}

	resp, err := svc.Run(req)
	if err != nil {
		log.Fatalf("backtest failed: %v", err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(resp); err != nil {
		log.Fatalf("encode response: %v", err)
	}
}
