This project includes a trading strategy implementation and backtest engine written in go. This is for educational purposes only.

What it contains:
- Engine
- Indicators
- Example strategies


File Structure:
signalstack/
├── cmd/
│   ├── backtest/          # CLI: run backtests
│   │   └── main.go
│   ├── data-backfill/     # CLI: backfill OHLCV data
│   │   └── main.go
│   └── analyze/           # CLI: analyze results, metrics, plots, etc.
│       └── main.go
│
├── internal/
│   ├── core/              # domain types, no I/O
│   │   ├── candle.go
│   │   ├── order.go
│   │   ├── position.go
│   │   ├── timeframe.go
│   │   └── portfolio.go
│   │
│   ├── engine/            # backtest engine
│   │   ├── engine.go
│   │   └── equity_curve.go
│   │
│   ├── indicators/        # SMA, EMA, RSI, ATR, trendlines, etc.
│   │   ├── sma.go
│   │   ├── ema.go
│   │   └── rsi.go
│   │
│   ├── strategies/        # strategy implementations
│   │   ├── interface.go   # Strategy interface
│   │   ├── ema_cross.go
│   │   └── trendline_bounce.go
│   │
│   ├── metrics/           # performance metrics
│   │   ├── metrics.go     # CAGR, max DD, win rate, etc.
│   │   └── report.go      # formatting summaries
│   │
│   ├── data/              # data sources + storage + backfill orchestration
│   │   ├── model.go       # shared types for symbol, source, etc.
│   │   ├── source/        # HTTP/data provider clients
│   │   │   ├── interface.go   # MarketDataSource interface
│   │   │   ├── yahoo.go       # example source
│   │   │   └── alpaca.go      # example source
│   │   ├── storage/       # how OHLCV is stored locally
│   │   │   ├── interface.go   # BarStore interface
│   │   │   ├── csv_store.go   # CSV-based store
│   │   │   └── sqlite_store.go
│   │   └── backfill/      # “pull from source → store” logic
│   │       └── backfill.go
│   │
│   ├── config/            # config structs + loading
│   │   ├── config.go
│   │   └── env.go
│   │
│   └── util/              # logging, small helpers
│       └── log.go
│
├── configs/
│   ├── backtest.example.yaml    # example backtest config
│   ├── datasources.example.yaml # API keys, sources, timeframes
│   └── strategies/
│       └── ema_cross.example.yaml
│
├── data/
│   ├── raw/                # raw pull (as-fetched, 1:1 with provider)
│   │   └── yahoo/
│   └── ohlcv/              # normalized OHLCV, your internal schema
│       ├── daily/
│       └── intraday/
│
├── scripts/
│   ├── download_sample_data.sh
│   └── gen_dummy_data.go   # small dev helper maybe
│
├── Makefile
├── go.mod
└── go.sum