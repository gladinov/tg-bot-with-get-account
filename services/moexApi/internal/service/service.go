package service

import (
	"context"
	"errors"
	"log/slog"
	"moex/internal/clients/moex"
	"moex/internal/models"
	"moex/internal/utils/logging"
	"time"

	"github.com/gladinov/e"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=ServiceClient
type ServiceClient interface {
	GetSpecifications(ctx context.Context, req models.SpecificationsRequest) (values models.Values, err error)
}

type Service struct {
	logger *slog.Logger
	client moex.MoexClient
}

func NewServiceClient(logger *slog.Logger, client moex.MoexClient) *Service {
	return &Service{
		logger: logger,
		client: client,
	}
}

func (c *Service) GetSpecifications(ctx context.Context, req models.SpecificationsRequest) (values models.Values, err error) {
	const op = "service.GetSpecifications"

	logg := c.logger.With()
	defer logging.LogOperation_Debug(ctx, logg, op, &err)()

	ticker := req.Ticker
	date := req.Date

	now := time.Now()
	date = clampDate(req.Date, now)

	var data models.SpecificationsResponce
	daysMax := 14
	var dayToSubstract int
	for dayToSubstract = 1; dayToSubstract <= daysMax; dayToSubstract++ {
		data, err = c.client.GetSpecifications(ctx, ticker, date)
		if err != nil {
			return models.Values{}, e.WrapIfErr("could not get specification from moexClient", err)
		}
		if data.History != nil {
			if len(data.History.Data) != 0 {
				break
			}
		}
		date = date.AddDate(0, 0, -1)
	}
	if dayToSubstract > daysMax {
		return models.Values{}, errors.New("could not find data in MOEX")
	}

	resp := data.History.Data[0]
	return resp, nil
}

func clampDate(date, now time.Time) time.Time {
	if date.After(now) {
		return now
	}
	return date
}
