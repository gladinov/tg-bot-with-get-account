package service

import (
	"bonds-report-service/internal/models/domain"
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/gladinov/e"
)

func (c *Client) GetCurrencyFromCB(ctx context.Context, charCode string, date time.Time) (vunit_rate float64, err error) {
	const op = "service.GetCurrencyFromCB"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("can't get currency from CB", err)
	}()

	vunit_rate, err = c.Storage.GetCurrency(context.Background(), charCode, date)
	if err == nil {
		return vunit_rate, nil
	}
	if err != nil && !errors.Is(err, domain.ErrNoCurrency) {
		return vunit_rate, err
	}

	currencies, err := c.External.Cbr.GetAllCurrencies(ctx, date)
	if err != nil {
		return vunit_rate, err
	}
	charCode = strings.ToLower(charCode)
	vunit_rate = currencies.CurrenciesMap[charCode].VunitRate

	err = c.Storage.SaveCurrency(context.Background(), currencies, date)
	if err != nil {
		return vunit_rate, err
	}
	return vunit_rate, nil
}
