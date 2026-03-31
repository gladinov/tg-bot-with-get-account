package cbr

import (
	"cbr/internal/models"
	"cbr/internal/utils/logging"
	"context"
	"log/slog"
	"net/url"
	"path"

	"github.com/gladinov/e"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=CbrClient
type CbrClient interface {
	GetAllCurrencies(ctx context.Context, formatDate string) (models.CurrenciesResponce, error)
}

type Client struct {
	transport TransportClient
	logger    *slog.Logger
}

func NewClient(logger *slog.Logger, transport TransportClient) *Client {
	return &Client{
		logger:    logger,
		transport: transport,
	}
}

func (c *Client) GetAllCurrencies(ctx context.Context, formatDate string) (_ models.CurrenciesResponce, err error) {
	const op = "cbr.GetAllCurrencies"
	logg := c.logger.With()
	defer logging.LogOperation_Debug(ctx, logg, op, &err)()

	Path := path.Join("scripts", "XML_daily.asp")

	params := url.Values{}
	params.Add("date_req", formatDate)

	body, err := c.transport.DoRequest(ctx, Path, params)
	if err != nil {
		return models.CurrenciesResponce{}, e.WrapIfErr("could not do request", err)
	}

	currResp, err := parseCurrencies(ctx, logg, body)
	if err != nil {
		return models.CurrenciesResponce{}, e.WrapIfErr("could not parse currencies", err)
	}

	return currResp, nil
}
