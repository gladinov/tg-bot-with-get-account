package usecases

import (
	"bonds-report-service/internal/domain"
	"bonds-report-service/internal/utils/logging"
	"context"
	"time"

	"github.com/gladinov/e"
)

func (s *Service) GetUsd(ctx context.Context) (_ domain.UsdResponce, err error) {
	const op = "service.GetUsd"

	defer logging.LogOperation_Debug(ctx, s.logger, op, &err)()

	usd, err := s.Helpers.CbrGetter.GetCurrencyFromCB(ctx, "usd", time.Now())
	if err != nil {
		return domain.UsdResponce{}, e.WrapIfErr("could not get usd from CB", err)
	}
	usdResponce := domain.UsdResponce{Usd: usd}

	return usdResponce, nil
}
