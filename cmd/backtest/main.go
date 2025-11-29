package main

import (
	"flag"
	"fmt"
	"log"
	"signalstack/internal/core"
	"signalstack/internal/data/storage"
	"signalstack/internal/engine"
	"signalstack/internal/strategy"
	"time"
)

func main() {
	var (
		symbol    = flag.String("symbol", "AAPL", "ticker symbol")
		timeframe = flag.String("timeframe", "1d", "timeframe")
		startDate = flag.String("start-date", "", "start date")
	)
	flag.Parse()
	store := storage.NewCSVStore("data/raw/yahoo")
	var startDateTimestamp *time.Time
	if *startDate != "" {
		parsed, err := time.Parse("2006-01-02", *startDate)
		if err != nil {
			log.Fatalf("Error parsing start date: %v", err)
		}
		startDateTimestamp = &parsed
	}
	candles, err := store.Fetch(core.Symbol(*symbol), core.Timeframe(*timeframe), startDateTimestamp)
	if err != nil {
		log.Fatalf("Error fetching candles: %v", err)
	}
	fmt.Println("Fetched", len(candles), "candles")
	strat := strategy.NewSMACross(20, 50)
	engine := engine.New(engine.EngineConfig{
		InitialCapital:     100000,
		SlippageBps:        10,
		CommissionPerTrade: 10,
		AllInOnEntry:       true,
	})
	result := engine.Run(strat, candles)
	fmt.Println("Result:", result)
}
