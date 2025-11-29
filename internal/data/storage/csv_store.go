package storage

import (
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
	"signalstack/internal/core"
	"strconv"
	"time"
)

type CSVStore struct {
	RootDir string
}

func (s *CSVStore) Path(symbol core.Symbol) string {
	return filepath.Join(s.RootDir, string(symbol)+".csv")
}

func (s *CSVStore) Store(symbol core.Symbol, timeframe core.Timeframe, candles []core.Candle) error {
	path := s.Path(symbol)
	os.MkdirAll(filepath.Dir(path), 0755)
	var alreadyExists bool = false
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		alreadyExists = true
	}
	csvFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer csvFile.Close()
	writer := csv.NewWriter(csvFile)
	// if this is a new file, write the header
	if !alreadyExists {
		writer.Write([]string{"Date", "Open", "High", "Low", "Close", "Volume"})
	}
	for _, candle := range candles {
		writer.Write([]string{
			candle.Timestamp.Format(time.RFC3339),
			strconv.FormatFloat(candle.Open, 'f', 2, 64),  // 2 decimal places
			strconv.FormatFloat(candle.High, 'f', 2, 64),  // 2 decimal places
			strconv.FormatFloat(candle.Low, 'f', 2, 64),   // 2 decimal places
			strconv.FormatFloat(candle.Close, 'f', 2, 64), // 2 decimal places
			strconv.Itoa(candle.Volume),
		})
	}
	writer.Flush()
	return nil
}

func (s *CSVStore) Fetch(symbol core.Symbol, timeframe core.Timeframe, start *time.Time) ([]core.Candle, error) {
	path := s.Path(symbol)
	csvFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()
	reader := csv.NewReader(csvFile)
	reader.Read() // skip header
	candles := []core.Candle{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		timestamp, err := time.Parse(time.RFC3339, record[0])
		if err != nil {
			return nil, err
		}
		open, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, err
		}
		high, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, err
		}
		low, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			return nil, err
		}
		close, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			return nil, err
		}
		volume, err := strconv.Atoi(record[5])
		if err != nil {
			return nil, err
		}
		candle := core.Candle{
			Timestamp: timestamp,
			Open:      open,
			High:      high,
			Low:       low,
			Close:     close,
			Volume:    volume,
		}
		if start != nil && candle.Timestamp.Before(*start) {
			continue
		}
		candles = append(candles, candle)
	}
	return candles, nil
}

func (s *CSVStore) LatestTimestamp(symbol core.Symbol, timeframe core.Timeframe) (*time.Time, error) {
	path := s.Path(symbol)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, nil
	}
	csvFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()
	reader := csv.NewReader(csvFile)
	// Read the last line
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(lines) == 0 {
		return nil, nil
	}
	lastLine := lines[len(lines)-1]
	timestamp, err := time.Parse(time.RFC3339, lastLine[0])
	if err != nil {
		return nil, err
	}
	return &timestamp, nil
}

func NewCSVStore(rootDir string) *CSVStore {
	return &CSVStore{
		RootDir: rootDir,
	}
}
