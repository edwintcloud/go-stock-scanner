package scanner

import (
	"fmt"
	"time"

	"github.com/edwintcloud/go-stock-scanner/internal/markethours"
)

// calculates the premarket gap from the 4:01am minute bar using the formula: (open - close) / close
func (s *Scanner) calculatePremarketGapPercent(ticker string) (float64, error) {
	premarketStart := markethours.PremarketSessionStartTime(time.Now())
	bars := s.tickerBars.GetBars(ticker, premarketStart, premarketStart.Add(1*time.Minute))
	if len(bars) == 0 {
		return 0, fmt.Errorf("no premarket bars for %s", ticker)
	}
	gapPct := (bars[0].Open - bars[0].Close) / bars[0].Close
	return gapPct, nil
}
