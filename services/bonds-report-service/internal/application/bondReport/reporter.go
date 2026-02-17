package bondreport

import (
	"log/slog"
	"time"
)

type BondReporter struct {
	logger *slog.Logger
	now    func() time.Time
}

func NewBondReporter(logger *slog.Logger) *BondReporter {
	return &BondReporter{
		logger: logger,
		now:    time.Now,
	}
}
