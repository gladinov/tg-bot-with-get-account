package moex

import (
	"context"
	"encoding/json"
	"log/slog"
	"moex/internal/models"
	"moex/internal/utils/logging"
	"net/url"
	"path"
	"time"

	"github.com/gladinov/e"
)

const (
	layout = "2006-01-02"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=MoexClient
type MoexClient interface {
	GetSpecifications(ctx context.Context, ticker string, date time.Time) (_ models.SpecificationsResponce, err error)
}

type Client struct {
	logger    *slog.Logger
	transport TransportClient
}

func NewMoexClient(logger *slog.Logger, transport TransportClient) *Client {
	return &Client{
		logger:    logger,
		transport: transport,
	}
}

func (c *Client) GetSpecifications(ctx context.Context, ticker string, date time.Time) (_ models.SpecificationsResponce, err error) {
	const op = "cbr.GetSpecifications"
	logg := c.logger.With()
	defer logging.LogOperation_Debug(ctx, logg, op, &err)()

	formatDate := date.Format(layout)
	path := path.Join("iss", "history", "engines", "stock", "markets", "bonds", "sessions", "3", "securities", ticker+".json")
	params := url.Values{}
	params.Add("limit", "1")
	params.Add("iss.meta", "off")
	params.Add("history.columns", "TRADEDATE,MATDATE,OFFERDATE,BUYBACKDATE,YIELDCLOSE,YIELDTOOFFER,FACEVALUE,FACEUNIT,DURATION, SHORTNAME")
	params.Add("limit", "1")
	params.Add("from", formatDate)
	params.Add("to", formatDate)

	body, err := c.transport.DoRequest(ctx, path, params)
	if err != nil {
		return models.SpecificationsResponce{}, e.WrapIfErr("failed DoRequest", err)
	}
	var data models.SpecificationsResponce
	err = json.Unmarshal(body, &data)
	if err != nil {
		return models.SpecificationsResponce{}, e.WrapIfErr("failed unmarshall json", err)
	}
	return data, nil
}
