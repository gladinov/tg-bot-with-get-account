package service

import (
	"bonds-report-service/internal/domain"
	"bonds-report-service/internal/utils/logging"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

func (s *Service) GetCurrencyFromCB(ctx context.Context, charCode string, date time.Time) (vunit_rate float64, err error) {
	const op = "service.GetCurrencyFromCB"

	defer logging.LogOperation_Debug(ctx, s.logger, op, &err)()

	vunit_rate, err = s.Storage.GetCurrency(ctx, charCode, date)
	if err == nil {
		return vunit_rate, nil
	}
	if err != nil && !errors.Is(err, domain.ErrNoCurrency) {
		return vunit_rate, err
	}

	currencies, err := s.External.Cbr.GetAllCurrencies(ctx, date)
	if err != nil {
		return vunit_rate, err
	}
	charCode = strings.ToLower(charCode)
	if val, ok := currencies.CurrenciesMap[charCode]; ok {
		vunit_rate = val.VunitRate
	} else {
		return 0, fmt.Errorf("currency %s not found", charCode)
	}

	err = s.Storage.SaveCurrency(ctx, currencies, date)
	if err != nil {
		s.logger.Warn("failed to save currency", "error", err)
	}
	return vunit_rate, nil
}
