package service

import (
	timezone "cbr/lib/timeZone"
	"fmt"
	"time"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=CurrencyService
type CurrencyService interface {
	GetAllCurrencies(date time.Time) (CurrenciesResponce, error)
}

type Service struct {
	Client HTTPClient
}

func NewService(client HTTPClient) *Service {
	return &Service{Client: client}
}

func (s *Service) GetAllCurrencies(date time.Time) (CurrenciesResponce, error) {
	const op = "service.GetAllCurrencies"
	location, err := timezone.GetMoscowLocation()
	if err != nil {
		return CurrenciesResponce{}, fmt.Errorf("op: %s, error: failed to load Moscow location", op)
	}
	now := time.Now().In(location)
	startDate := timezone.GetStartSingleExchangeRateRubble(location)

	formatDate := normalizeDate(date, now, startDate)

	return s.Client.GetAllCurrencies(formatDate)
}

func normalizeDate(date, now, startDate time.Time) string {
	switch {
	case date.After(now):
		return now.Format(layout)
	case date.Before(startDate):
		return startDate.Format(layout)
	default:
		return date.Format(layout)
	}
}
