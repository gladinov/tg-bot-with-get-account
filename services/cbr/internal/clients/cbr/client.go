package cbr

import (
	"bytes"
	"cbr/internal/models"
	"cbr/internal/utils/logging"
	"context"
	"encoding/xml"
	"io"
	"log/slog"
	"net/url"
	"path"

	"github.com/gladinov/e"
	"golang.org/x/text/encoding/charmap"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=HTTPClient
type HTTPClient interface {
	GetAllCurrencies(ctx context.Context, formatDate string) (models.CurrenciesResponce, error)
}

type Client struct {
	transport HTTPTransport
	logger    *slog.Logger
}

func NewClient(logger *slog.Logger, transport HTTPTransport) *Client {
	return &Client{
		logger:    logger,
		transport: transport,
	}
}

func (c *Client) GetAllCurrencies(ctx context.Context, formatDate string) (_ models.CurrenciesResponce, err error) {
	const op = "service.GetAllCurrencies"
	logg := c.logger.With()
	defer logging.LogOperation_Debug(ctx, logg, op, &err)()

	Path := path.Join("scripts", "XML_daily.asp")

	params := url.Values{}
	params.Add("date_req", formatDate)

	body, err := c.transport.DoRequest(ctx, Path, params)
	if err != nil {
		return models.CurrenciesResponce{}, e.WrapIfErr("could not do request", err)
	}

	currResp, err := c.parseCurrencies(ctx, body)
	if err != nil {
		return models.CurrenciesResponce{}, e.WrapIfErr("could not parse currencies", err)
	}

	return currResp, nil
}

func (c *Client) parseCurrencies(ctx context.Context, data []byte) (_ models.CurrenciesResponce, err error) {
	const op = "service.parseCurrencies"
	logg := c.logger.With()
	defer logging.LogOperation_Debug(ctx, logg, op, &err)()

	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.CharsetReader = func(label string, input io.Reader) (io.Reader, error) {
		if label == "windows-1251" {
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		}
		return input, nil
	}
	var curr models.CurrenciesResponce
	err = decoder.Decode(&curr)
	if err != nil {
		return models.CurrenciesResponce{}, e.WrapIfErr("could not decode Xml file", err)
	}

	return curr, nil
}
