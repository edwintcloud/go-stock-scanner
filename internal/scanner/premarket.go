package scanner

import (
	"context"
	"fmt"
	"time"

	"github.com/edwintcloud/go-stock-scanner/internal/domain"
	"github.com/edwintcloud/go-stock-scanner/internal/markethours"

	massiverest "github.com/massive-com/client-go/v3/rest"
	restmodel "github.com/massive-com/client-go/v3/rest/gen"
)

// calculates the premarket gap from the 4:01am minute bar using the formula: (open - close) / close
func (s *Scanner) calculatePremarketGapPercent(ticker string) (float64, error) {
	premarketStart := markethours.PremarketSessionStartTime(time.Now())
	bars := s.tickerBars.GetBars(ticker, premarketStart, premarketStart.Add(1*time.Minute))
	if len(bars) == 0 {
		bars, err := s.fetchBars(ticker, premarketStart, premarketStart.Add(1*time.Minute))
		if err != nil {
			return 0, fmt.Errorf("error fetching bars for %s: %w", ticker, err)
		}
		if len(bars) == 0 {
			return 0, fmt.Errorf("no bars found for %s in premarket", ticker)
		}
		s.tickerBars.AddBar(bars[0])
	}
	gapPct := (bars[0].Open - bars[0].Close) / bars[0].Close
	return gapPct, nil
}

func (s *Scanner) fetchBars(ticker string, start, end time.Time) ([]domain.Bar, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	params := &restmodel.GetStocksAggregatesParams{
		Limit: massiverest.Ptr(50000),
	}
	resp, err := s.rest.GetStocksAggregatesWithResponse(
		ctx,
		ticker,
		1,
		restmodel.Minute,
		fmt.Sprintf("%d", start.In(markethours.Location).UnixMilli()),
		fmt.Sprintf("%d", end.In(markethours.Location).UnixMilli()),
		params,
	)
	if err != nil {
		return nil, err
	}
	if resp.JSON200 == nil || resp.JSON200.Results == nil {
		return nil, fmt.Errorf("no bars in response for %s", ticker)
	}
	bars := make([]domain.Bar, len(*resp.JSON200.Results))
	for i, b := range *resp.JSON200.Results {
		bars[i] = domain.Bar{
			Ticker:         ticker,
			Open:           b.O,
			Close:          b.C,
			High:           b.H,
			Low:            b.L,
			Volume:         uint64(b.V),
			CloseTimestamp: time.UnixMilli(int64(b.Timestamp)).In(markethours.Location).Add(1 * time.Minute), // bars come back with 4:00am timestamp, but we want to treat them as 4:01am bars
		}
	}
	return bars, nil
}
