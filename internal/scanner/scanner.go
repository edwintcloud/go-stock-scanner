package scanner

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/edwintcloud/go-stock-scanner/internal/config"
	"github.com/edwintcloud/go-stock-scanner/internal/domain"
	"github.com/labstack/gommon/log"
	massiverest "github.com/massive-com/client-go/v3/rest"
	massivews "github.com/massive-com/client-go/v3/websocket"
	wsmodels "github.com/massive-com/client-go/v3/websocket/models"
)

type Scanner struct {
	config            *config.Config
	ws                *massivews.Client
	rest              *massiverest.Client
	tickerBars        domain.TickerBars
	tickerSnapshotMap map[string]TickerSnapshot
	blockedSymbols    map[string]bool
	currentCanidates  []domain.Bar
}

func NewScanner(config *config.Config) (*Scanner, error) {
	ws, err := massivews.New(massivews.Config{
		APIKey: config.MassiveAPIKey,
		Feed:   massivews.RealTime,
		Market: massivews.Stocks,
	})
	if err != nil {
		return nil, err
	}
	return &Scanner{
		config:            config,
		ws:                ws,
		rest:              massiverest.NewWithOptions(config.MassiveAPIKey, massiverest.WithPagination(true)),
		tickerBars:        make(domain.TickerBars),
		tickerSnapshotMap: make(map[string]TickerSnapshot),
		blockedSymbols:    make(map[string]bool),
		currentCanidates:  make([]domain.Bar, 0),
	}, nil
}

func (s *Scanner) Start(ctx context.Context) error {
	defer s.ws.Close()

	go s.refreshTickerSnapshotMapOnInterval(ctx, 15*time.Minute)

	err := s.ws.Subscribe(massivews.StocksMinAggs, "*")
	if err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}

	if err := s.ws.Connect(); err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-s.ws.Error():
			return fmt.Errorf("websocket error: %w", err)
		case out, more := <-s.ws.Output():
			if !more {
				return fmt.Errorf("output channel closed")
			}
			switch v := out.(type) {
			case wsmodels.EquityAgg:
				bar := domain.BarFromEquityAgg(v)
				s.tickerBars.AddBar(bar)
				if s.shouldSkipBar(bar) {
					continue
				}
				s.currentCanidates = append(s.currentCanidates, bar)
			default:
				log.Debugf("unknown type: %T\n", out)
				continue
			}
		}
		s.emitCanidates()
	}
}

func (s *Scanner) emitCanidates() {
	if len(s.currentCanidates) < 10 {
		return
	}
	slices.SortFunc(s.currentCanidates, func(a, b domain.Bar) int {
		if b.TodaysVolume > a.TodaysVolume {
			return 1
		}
		return -1
	})
	fmt.Println("Top 10 canidates:")
	for i := range 10 {
		c := s.currentCanidates[i]
		fmt.Printf("%s: Close=%.2f, Volume=%d\n", c.Ticker, c.Close, c.TodaysVolume)
	}
	fmt.Println()
	s.currentCanidates = make([]domain.Bar, 0)
}
