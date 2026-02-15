package bondreport

import (
	"bonds-report-service/internal/models/domain"
	"context"
	"log/slog"
	"time"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=MoexClient
type MoexClient interface {
	GetSpecifications(ctx context.Context, ticker string, date time.Time) (data domain.ValuesMoex, err error)
}

type BondReporter struct {
	logger *slog.Logger
	Moex   MoexClient
	now    func() time.Time
}

func NewBondReporter(logger *slog.Logger, moex MoexClient) *BondReporter {
	return &BondReporter{
		logger: logger,
		Moex:   moex,
	}
}
