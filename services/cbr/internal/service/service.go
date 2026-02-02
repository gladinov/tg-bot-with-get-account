package service

import (
	"cbr/internal/clients/cbr"
	"cbr/internal/models"
	"cbr/internal/utils"
	"context"
	"log/slog"
	"time"

	"github.com/gladinov/e"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=CurrencyService
type CurrencyService interface {
	GetAllCurrencies(ctx context.Context, date time.Time) (models.CurrenciesResponce, error)
}

type Service struct {
	Logger       *slog.Logger
	Client       cbr.HTTPClient
	TimeLocation *time.Location
}

func NewService(logger *slog.Logger, client cbr.HTTPClient, timeLocation *time.Location) *Service {
	return &Service{
		Logger:       logger,
		Client:       client,
		TimeLocation: timeLocation,
	}
}

func (s *Service) GetAllCurrencies(ctx context.Context, date time.Time) (models.CurrenciesResponce, error) {
	const op = "service.GetAllCurrencies"
	location := s.TimeLocation
	now := time.Now().In(location)
	startDate := utils.GetStartSingleExchangeRateRubble(location)

	formatDate := utils.NormalizeDate(date, now, startDate)

	currResp, err := s.Client.GetAllCurrencies(ctx, formatDate)
	if err != nil {
		return models.CurrenciesResponce{}, e.WrapIfErr("failed to get all currencies from client", err)
	}

	return currResp, nil
}
