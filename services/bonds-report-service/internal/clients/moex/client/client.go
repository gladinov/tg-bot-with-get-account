package moex

import (
	"bonds-report-service/internal/clients/moex/transport"
	"bonds-report-service/internal/models/domain"
	dtoMoex "bonds-report-service/internal/models/dto/moex"
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

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=MoexClient
type MoexClient interface {
	GetSpecifications(ctx context.Context, ticker string, date time.Time) (data domain.Values, err error)
}

type Client struct {
	logger    *slog.Logger
	transport transport.TransportClient
}

func NewMoexClient(logger *slog.Logger, transport transport.TransportClient) *Client {
	return &Client{
		logger:    logger,
		transport: transport,
	}
}

func (c *Client) GetSpecifications(ctx context.Context, ticker string, date time.Time) (data domain.Values, err error) {
	const op = "moex.GetSpecifications"
	logg := c.logger.With()
	defer logging.LogOperation_Debug(ctx, logg, op, &err)()

	request := dtoMoex.NewSpecificationsRequest(ticker, date)
	Path := path.Join("moex", "specifications")
	params := url.Values{}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return domain.Values{}, e.WrapIfErr("failed json.Marshal", err)
	}
	formatRequestBody := bytes.NewBuffer(requestBody)

	httpResponse, err := c.transport.DoRequest(ctx, Path, params, formatRequestBody)
	if err != nil {
		return domain.Values{}, e.WrapIfErr("failed transport DoRequest", err)
	}

	switch httpResponse.StatusCode {
	case http.StatusBadRequest:
		return domain.Values{}, errors.New("Do request status is bad request")
	case http.StatusInternalServerError:
		return domain.Values{}, errors.New("Do request status is internal server error")
	}

	var res dtoMoex.Values
	err = json.Unmarshal(httpResponse.Body, &res)
	if err != nil {
		return domain.Values{}, e.WrapIfErr("failed to unmarshal response", err)
	}

	domainRes := MapValueFromDTOToDomain(res)

	return domainRes, nil
}
