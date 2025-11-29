// Yahoo Finance API client

package source

import (
	"net/http"
	"signalstack/internal/core"
	"time"

	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
)

type YahooDataSource struct {
}

func (y *YahooDataSource) Name() string {
	return "yahoo"
}

func (y *YahooDataSource) FetchOHLCV(symbol core.Symbol, timeframe core.Timeframe, start time.Time, end time.Time) ([]core.Candle, error) {
	// Set custom HTTP client with User-Agent
	finance.SetHTTPClient(&http.Client{
		Transport: &userAgentTransport{
			rt: http.DefaultTransport,
		},
		Timeout: 20 * time.Second,
	})
	params := &chart.Params{
		Symbol:   string(symbol),
		Interval: datetime.Interval(timeframe),
		Start:    datetime.New(&start),
		End:      datetime.New(&end),
	}
	iter := chart.Get(params)
	var candles []core.Candle
	for iter.Next() {
		bar := iter.Bar()
		open, _ := bar.Open.Float64()
		high, _ := bar.High.Float64()
		low, _ := bar.Low.Float64()
		close, _ := bar.Close.Float64()
		candle := core.Candle{
			Timestamp: *datetime.FromUnix(bar.Timestamp).Time(),
			Open:      open,
			High:      high,
			Low:       low,
			Close:     close,
			Volume:    bar.Volume,
		}
		candles = append(candles, candle)
	}
	return candles, nil
}

type userAgentTransport struct {
	rt http.RoundTripper
}

func (t *userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	return t.rt.RoundTrip(req)
}
