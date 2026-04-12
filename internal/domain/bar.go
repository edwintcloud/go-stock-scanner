package domain

import (
	"time"

	"github.com/edwintcloud/go-stock-scanner/internal/markethours"
	wsmodels "github.com/massive-com/client-go/v3/websocket/models"
)

type Bar struct {
	Ticker         string
	Open           float64
	Close          float64
	High           float64
	Low            float64
	Volume         uint64
	CloseTimestamp time.Time
	TodaysVolume   uint64
	TodaysVWAP     float64
}

func BarFromEquityAgg(equityAgg wsmodels.EquityAgg) Bar {
	return Bar{
		Ticker:         equityAgg.Symbol,
		Open:           equityAgg.Open,
		Close:          equityAgg.Close,
		High:           equityAgg.High,
		Low:            equityAgg.Low,
		Volume:         uint64(equityAgg.Volume),
		CloseTimestamp: time.UnixMilli(equityAgg.EndTimestamp).In(markethours.Location),
		TodaysVolume:   uint64(equityAgg.AccumulatedVolume),
		TodaysVWAP:     equityAgg.VWAP,
	}
}
