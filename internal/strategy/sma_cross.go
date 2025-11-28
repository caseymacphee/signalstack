package strategy

type SMACross struct {
	shortWindow int
	longWindow int
	closes []float64
	shortSum float64
	longSum float64

}

func NewSMACross(shortWindow int, longWindow int) *SMACross {
	if shortWindow <= 0 || longWindow <= 0 || shortWindow >= longWindow {
		panic("invalid window sizes")
	}
	return &SMACross{
		shortWindow: shortWindow,
		longWindow: longWindow,
		closes: make([]float64, 0, longWindow*2),
	}
}

func (s *SMACross) Name() string {
	return "SMA_CROSS"
}


func (s *SMACross) OnBar(ctx Context) Decision {
	c := ctx.Candle
	price := c.Close

	// update internal state
	s.closes = append(s.closes, price)

	// update rolling sums
	n := len(s.closes)
	s.shortSum = 0
	s.longSum = 0
	if n > s.shortWindow {
		for i := n - s.shortWindow; i < n; i ++ {
			s.shortSum += s.closes[i]
		}
	}
	if n > s.longWindow {
		for i := n - s.longWindow; i < n; i ++ {
			s.longSum += s.closes[i]
		}
	}
    if n < s.longWindow {
        return Decision{}
    }

	// compute SMAs
	shortSMA := s.shortSum / float64(s.shortWindow)
	longSMA := s.longSum / float64(s.longWindow)

	if ctx.Position == nil && shortSMA > longSMA {
		return Decision{
			EnterLong: true,
			StopLoss: nil,
			TakeProfit: nil,
		}
	}
	if ctx.Position != nil && shortSMA < longSMA {
		return Decision{
			ExitLong: true,
			StopLoss: nil,
			TakeProfit: nil,
		}
	}
	return Decision{}
}