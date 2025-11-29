package api

import "time"

type BacktestRequest struct {
	JobID     string            `json:"jobId"`
	Symbol    string            `json:"symbol"`    // "AAPL"
	Timeframe string            `json:"timeframe"` // "1D"
	StartDate time.Time         `json:"startDate"`
	EndDate   time.Time         `json:"endDate"`
	Strategy  string            `json:"strategyId"` // e.g. "SMA_CROSS" or other registered strategy ID
	Params    map[string]string `json:"params"`     // strategy-specific parameters
}

type BacktestMetrics struct {
	CAGR             float64 `json:"cagr"`
	TotalReturn      float64 `json:"totalReturn"` // strategy total return (fraction)
	MaxDrawdown      float64 `json:"maxDrawdown"` // fraction
	WinRate          float64 `json:"winRate"`     // fraction
	ProfitFactor     float64 `json:"profitFactor"`
	FinalEquity      float64 `json:"finalEquity"`
	AverageRMultiple float64 `json:"averageRMultiple"`
	BuyAndHoldReturn float64 `json:"buyAndHoldReturn"` // buy & hold return (fraction)
}

type BacktestResponse struct {
	JobID    string          `json:"jobId"`
	Symbol   string          `json:"symbol"`
	Strategy string          `json:"strategy"`
	Metrics  BacktestMetrics `json:"metrics"`
	Trades   []TradeDTO      `json:"trades"`
}

type TradeDTO struct {
	Side      string    `json:"side"`     // "BUY"/"SELL"
	Position  string    `json:"position"` // "LONG"/"SHORT"
	EntryTime time.Time `json:"entryTime"`
	Entry     float64   `json:"entryPrice"`
	ExitTime  time.Time `json:"exitTime"`
	Exit      float64   `json:"exitPrice"`
	Size      float64   `json:"size"`
	PnL       float64   `json:"pnl"`
	RMultiple *float64  `json:"rMultiple,omitempty"`
}
