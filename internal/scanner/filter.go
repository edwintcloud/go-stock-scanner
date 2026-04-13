package scanner

import (
	"strings"

	"github.com/edwintcloud/go-stock-scanner/internal/domain"
	"github.com/labstack/gommon/log"
)

func (s *Scanner) shouldSkipBar(bar domain.Bar) bool {
	if s.isBlockedSymbol(bar) ||
		s.isOutOfPriceRange(bar) ||
		s.hasLowVolume(bar) ||
		s.hasLowRelativeVolume(bar) ||
		s.isNotAGainer(bar) {
		return true
	}
	return false
}

func (s *Scanner) isBlockedSymbol(bar domain.Bar) bool {
	blockedTypes := []string{"ETF", "ETN", "REIT"}
	for _, t := range blockedTypes {
		if strings.Contains(bar.Ticker, t) {
			return true
		}
	}
	return s.blockedSymbols[bar.Ticker]
}

func (s *Scanner) isOutOfPriceRange(bar domain.Bar) bool {
	return bar.Close < s.config.ScannerOptions.MinPrice || bar.Close > s.config.ScannerOptions.MaxPrice
}

func (s *Scanner) hasLowVolume(bar domain.Bar) bool {
	return bar.TodaysVolume < s.config.ScannerOptions.MinVolume
}

func (s *Scanner) hasLowRelativeVolume(bar domain.Bar) bool {
	tickerSnapshot, ok := s.tickerSnapshotMap[bar.Ticker]
	if !ok {
		log.Debugf("no snapshot data for %s, adding to blocked symbols", bar.Ticker)
		s.blockedSymbols[bar.Ticker] = true
		return true // if we don't have snapshot data, skip
	}

	return float64(bar.TodaysVolume)/float64(tickerSnapshot.PreviousDayVolume) < s.config.ScannerOptions.MinRelativeVolume
}

func (s *Scanner) isNotAGainer(bar domain.Bar) bool {
	tickerSnapshot, ok := s.tickerSnapshotMap[bar.Ticker]
	if !ok {
		return true // if we don't have snapshot data, skip
	}
	return tickerSnapshot.DayChangePercent < s.config.ScannerOptions.MinChangePercent
}
