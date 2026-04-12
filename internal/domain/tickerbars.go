package domain

import "time"

type TickerBars map[string]BarsBucket

func (tb TickerBars) AddBar(bar Bar) {
	bucket, exists := tb[bar.Ticker]
	if !exists {
		tb[bar.Ticker] = NewBarsBucket()
	}
	bucket.AddBar(bar)
}

func (tb TickerBars) GetBars(ticker string, start, end time.Time) []Bar {
	bucket, exists := tb[ticker]
	if !exists {
		return nil
	}
	return bucket.GetBars(start, end)
}
