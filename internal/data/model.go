package data

import "signalstack/internal/core"

type Symbol struct {
	Symbol core.Symbol
	Source Source
}

type Source string

const (
	SourceYahoo Source = "yahoo"
)
