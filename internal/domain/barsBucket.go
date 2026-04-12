package domain

import (
	"time"

	"github.com/edwintcloud/go-stock-scanner/internal/markethours"
)

// BarsBucket represents a collection of bars for a specific ticker
type BarsBucket []Bar

func NewBarsBucket() BarsBucket {
	return make(BarsBucket, 960) // 16 hours of trading * 60 minutes
}

func (b *BarsBucket) AddBar(bar Bar) {
	i := b.idx(bar.CloseTimestamp)
	if i == -1 {
		return // out of range, ignore
	}
	(*b)[i] = bar
}

func (b *BarsBucket) GetBars(start, end time.Time) []Bar {
	startIndex := b.idx(start)
	endIndex := b.idx(end)
	if startIndex == -1 || endIndex == -1 || startIndex > endIndex {
		return nil // out of range or invalid range
	}
	return (*b)[startIndex : endIndex+1]
}

func (b *BarsBucket) idx(timestamp time.Time) int {
	n := len(*b)
	minutesSinceOpen := markethours.MinutesSinceMarketOpen(timestamp) - 1 // -1 because the first bar (4:01) should be at index 0
	if minutesSinceOpen < 0 || minutesSinceOpen >= n {
		return -1 // out of range
	}
	return minutesSinceOpen
}
