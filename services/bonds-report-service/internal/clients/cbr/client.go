package cbr

import (
	"bonds-report-service/internal/models"
	"bonds-report-service/internal/utils/logging"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/gladinov/e"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=CbrClient
type CbrClient interface {
	GetAllCurrencies(ctx context.Context, date time.Time) (res models.CurrenciesResponce, err error)
}

type Client struct {
	logger    *slog.Logger
	transport TransportClient
}

func NewCbrClient(logger *slog.Logger, transport TransportClient) *Client {
	return &Client{
		logger:    logger,
		transport: transport,
	}
}

func (c *Client) GetAllCurrencies(ctx context.Context, date time.Time) (res models.CurrenciesResponce, err error) {
	const op = "cbr.GetAllCurrencies"
	logg := c.logger.With()
	defer logging.LogOperation_Debug(ctx, logg, op, &err)()

	request := models.NewCurrencyRequest(date)
	Path := path.Join("cbr", "currencies")
	params := url.Values{}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return models.CurrenciesResponce{}, e.WrapIfErr("failed json.Marshal", err)
	}
	formatRequestBody := bytes.NewBuffer(requestBody)

	httpResponse, err := c.transport.DoRequest(ctx, Path, params, formatRequestBody)
	if err != nil {
		return models.CurrenciesResponce{}, e.WrapIfErr("failed transport DoRequest", err)
	}

	switch httpResponse.StatusCode {
	case http.StatusBadRequest:
		return models.CurrenciesResponce{}, errors.New("Do request status is bad request")
	case http.StatusInternalServerError:
		return models.CurrenciesResponce{}, errors.New("Do request status is internal server error")
	}

	err = json.Unmarshal(httpResponse.Body, &res)
	if err != nil {
		return models.CurrenciesResponce{}, e.WrapIfErr("failed to unmarshal response", err)
	}

	return res, nil
}
