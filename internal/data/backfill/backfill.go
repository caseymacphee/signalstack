package backfill

import (
	"signalstack/internal/core"
	"signalstack/internal/data/source"
	"signalstack/internal/data/storage"
	"time"
)


type BackfillService struct {
	Src	source.MarketDataSource
	Storage storage.BarStore
}


type BackfillRequest struct {
	Symbol core.Symbol
	Timeframe core.Timeframe
	StartDate *time.Time
	EndDate *time.Time
}


func (b *BackfillService) Backfill(request BackfillRequest) error {
	from := request.StartDate
	if from == nil {
		latest, err := b.Storage.LatestTimestamp(request.Symbol, request.Timeframe)
		if err != nil {
			return err
		}
		if latest == nil {
			// default to 1 year ago
			oneYearAgo := time.Now().AddDate(-1, 0, 0)
			from = &oneYearAgo
		} else {
			// add 1 day to the latest timestamp
			oneDayAfter := latest.AddDate(0, 0, 1)
			from = &oneDayAfter
		}
	}
	if request.EndDate == nil {
		now := time.Now()
		request.EndDate = &now
	}
	candles, err := b.Src.FetchOHLCV(request.Symbol, request.Timeframe, *from, *request.EndDate)
	if err != nil {
		return err
	}
	if len(candles) == 0 {
		return nil
	}
	return b.Storage.Store(request.Symbol, request.Timeframe, candles)
}