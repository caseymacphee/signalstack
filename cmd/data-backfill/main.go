package main

import (
	"flag"
	"fmt"

	"signalstack/internal/core"
	"signalstack/internal/data/backfill"
	"signalstack/internal/data/source"
	"signalstack/internal/data/storage"
)


func main() {
	var (
		symbolStr = flag.String("symbol", "AAPL", "ticker symbol")
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
	backfillReq := backfill.BackfillRequest{
		Symbol: symbol,
		Timeframe: core.TimeframeDaily,
		StartDate: nil,
		EndDate: nil,
	}
	err := backfillSvc.Backfill(backfillReq)
	if err != nil {
		fmt.Println("Error backfilling data:", err)
		return
	}

}