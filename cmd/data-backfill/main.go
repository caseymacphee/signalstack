package main

import (
	"flag"
	"fmt"
	"time"

	"signalstack/internal/core"
	"signalstack/internal/data/backfill"
	"signalstack/internal/data/source"
	"signalstack/internal/data/storage"
)


func main() {
	var (
		symbolStr = flag.String("symbol", "AAPL", "ticker symbol")
		timeframeStr = flag.String("timeframe", "1d", "timeframe")
		startDateStr = flag.String("start-date", "", "start date")
		endDateStr = flag.String("end-date", "", "end date")
	)

	flag.Parse()
	symbol := core.Symbol(*symbolStr)
	fmt.Println("Data Backfill for", symbol)

	client := &source.YahooDataSource{}
	store := &storage.CSVStore{
		RootDir: "data/raw/yahoo",
	}
	backfillSvc := backfill.BackfillService{
		Src: client,
		Storage: store,
	}
	timeframe := core.Timeframe(*timeframeStr)
	var startDate *time.Time
	var endDate *time.Time
	var err error

	if *startDateStr != "" {
		parsed, err := time.Parse("2006-01-02", *startDateStr)
		if err != nil {
			fmt.Printf("Error parsing start date: %v\n", err)
			return
		}
		startDate = &parsed
	} else {
		oneYearAgo := time.Now().AddDate(-1, 0, 0)
		startDate = &oneYearAgo
	}

	if *endDateStr != "" {
		parsed, err := time.Parse("2006-01-02", *endDateStr)
		if err != nil {
			fmt.Printf("Error parsing end date: %v\n", err)
			return
		}
		endDate = &parsed
	} else {
		now := time.Now()
		endDate = &now
	}
	backfillReq := backfill.BackfillRequest{
		Symbol: symbol,
		Timeframe: timeframe,
		StartDate: startDate,
		EndDate: endDate,
	}
	err = backfillSvc.Backfill(backfillReq)
	if err != nil {
		fmt.Println("Error backfilling data:", err)
		return
	}

}