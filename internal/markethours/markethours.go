package markethours

import "time"

var Location, _ = time.LoadLocation("America/New_York")

func PremarketSessionStartTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 4, 0, 0, 0, Location)
}

func RegularSessionStartTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 9, 30, 0, 0, Location)
}

func RegularSessionCloseTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 16, 0, 0, 0, Location)
}

func PostmarketSessionEndTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 20, 0, 0, 0, Location)
}

func MinutesSinceMarketOpen(t time.Time) int {
	t = t.In(Location)
	marketOpen := PremarketSessionStartTime(t)
	if t.Before(marketOpen) || t.After(PostmarketSessionEndTime(t)) {
		return -1 // outside of market hours
	}
	return int(t.Sub(marketOpen).Minutes())
}
