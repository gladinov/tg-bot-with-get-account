package bondreport

import (
	"bonds-report-service/internal/models/domain"
	"bonds-report-service/internal/utils/logging"
	"context"
	"time"

	"github.com/gladinov/e"
)

func (s *BondReporter) GetSpecificationsFromMoex(ctx context.Context, ticker string, date time.Time) (data domain.ValuesMoex, err error) {
	const op = "service.GetSpecificationsFromMoex"
	defer logging.LogOperation_Debug(ctx, s.logger, op, &err)()
	if ticker == "" {
		return domain.ValuesMoex{}, domain.ErrEmptyTicker
	}
	data, err = s.Moex.GetSpecifications(ctx, ticker, date)
	if err != nil {
		return domain.ValuesMoex{}, e.WrapIfErr("failed to get specification from MOEX", err)
	}
	return data, nil
}
