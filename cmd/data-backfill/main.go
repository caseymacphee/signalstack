package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
)


func fetchData(symbol string, interval string) {
	oneYearAgo := time.Now().AddDate(-1, 0, 0)
	yearAgoDatetime := datetime.New(&oneYearAgo)
	now := time.Now()
	intervalEnum := datetime.Interval(interval)
	params := &chart.Params{
		Symbol: symbol,
		Interval: intervalEnum,
		Start: yearAgoDatetime,
		End: datetime.New(&now),
	}
	iter := chart.Get(params)
	// write header to csv file
	csvFile, err := os.Create("data/raw/yahoo/" + symbol + ".csv")
	if err != nil {
		fmt.Println("Error creating csv file:", err)
		return
	}
	defer csvFile.Close()
	writer := csv.NewWriter(csvFile)
	writer.Write([]string{"Date", "Open", "High", "Low", "Close", "Volume"})
	for iter.Next() {
		// append to csv file
		bar := iter.Bar()
		newLine := []string{datetime.FromUnix(bar.Timestamp).Time().In(time.UTC).Format(time.RFC3339), bar.Open.String(), bar.High.String(), bar.Low.String(), bar.Close.String(), strconv.Itoa(bar.Volume)}
		writer.Write(newLine)
	}
	if err := iter.Err(); err != nil {
		fmt.Println("Error during iteration:", err)
	}
	writer.Flush()
}

func main() {
	symbol := os.Args[1]
	interval := os.Args[2]
	
	// Set custom HTTP client with User-Agent
	finance.SetHTTPClient(&http.Client{
		Transport: &userAgentTransport{
			rt: http.DefaultTransport,
		},
		Timeout: 20 * time.Second,
	})

	fmt.Println("Data Backfill for", symbol)
	fetchData(symbol, interval)
}

type userAgentTransport struct {
	rt http.RoundTripper
}

func (t *userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	return t.rt.RoundTrip(req)
}
