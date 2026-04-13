package scanner

import (
	"context"
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
	restmodel "github.com/massive-com/client-go/v3/rest/gen"
)

type TickerSnapshot struct {
	DayOpen           float64
	DayClose          float64
	DayHigh           float64
	DayLow            float64
	DayVolume         uint64
	DayChangePercent  float64
	PreviousDayOpen   float64
	PreviousDayClose  float64
	PreviousDayHigh   float64
	PreviousDayLow    float64
	PreviousDayVolume uint64
}

func (s *Scanner) refreshTickerSnapshotMapOnInterval(ctx context.Context, interval time.Duration) {
	// initial load
	tickerSnapshotMap, err := s.getTickerSnapshotMap(ctx)
	if err != nil {
		log.Errorf("Error loading initial ticker snapshot map: %v", err)
	}
	s.tickerSnapshotMap = tickerSnapshotMap
	log.Debugf("Initial ticker snapshot map loaded with %d entries", len(tickerSnapshotMap))

	// refresh on interval
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			tickerSnapshotMap, err := s.getTickerSnapshotMap(ctx)
			if err != nil {
				log.Errorf("Error refreshing ticker snapshot map: %v", err)
				continue
			}
			s.tickerSnapshotMap = tickerSnapshotMap
			log.Debugf("Ticker snapshot map refreshed with %d entries", len(tickerSnapshotMap))
		}
	}
}

func (s *Scanner) getTickerSnapshotMap(ctx context.Context) (map[string]TickerSnapshot, error) {
	params := &restmodel.GetStocksSnapshotTickersParams{}
	resp, err := s.rest.GetStocksSnapshotTickersWithResponse(ctx, params)
	if err != nil {
		return nil, err
	}
	if resp.JSON200 == nil || resp.JSON200.Tickers == nil {
		return nil, fmt.Errorf("invalid response")
	}

	tickerSnapshotMap := make(map[string]TickerSnapshot)
	for _, ticker := range *resp.JSON200.Tickers {
		if ticker.Ticker == nil {
			continue
		}
		tickerSnapshotMap[*ticker.Ticker] = TickerSnapshot{
			DayOpen:           ticker.Day.O,
			DayClose:          ticker.Day.C,
			DayHigh:           ticker.Day.H,
			DayLow:            ticker.Day.L,
			DayVolume:         uint64(ticker.Day.V),
			DayChangePercent:  *ticker.TodaysChangePerc,
			PreviousDayOpen:   ticker.PrevDay.O,
			PreviousDayClose:  ticker.PrevDay.C,
			PreviousDayHigh:   ticker.PrevDay.H,
			PreviousDayLow:    ticker.PrevDay.L,
			PreviousDayVolume: uint64(ticker.PrevDay.V),
		}
	}
	return tickerSnapshotMap, nil
}
