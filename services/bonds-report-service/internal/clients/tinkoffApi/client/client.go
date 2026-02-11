package tinkoffApi

import (
	"bonds-report-service/internal/clients/tinkoffApi/client/analyticsclient"
	"bonds-report-service/internal/clients/tinkoffApi/client/instrumentsclient"
	"bonds-report-service/internal/clients/tinkoffApi/client/portfolioclient"
	"log/slog"
)

type Client struct {
	logger                   *slog.Logger
	instrumentsTinkoffClient instrumentsclient.InstrumentsClient
	portfolioTinkoffClient   portfolioclient.PortfolioClient
	analyticsTinkoffClient   analyticsclient.AnalyticsClient
}

func NewTinkoffClient(logger *slog.Logger,
	instrumentsTinkoffClient instrumentsclient.InstrumentsClient,
	portfolioTinkoffClient portfolioclient.PortfolioClient,
	analyticsTinkoffClient analyticsclient.AnalyticsClient,
) *Client {
	return &Client{
		logger:                   logger,
		instrumentsTinkoffClient: instrumentsTinkoffClient,
		portfolioTinkoffClient:   portfolioTinkoffClient,
		analyticsTinkoffClient:   analyticsTinkoffClient,
	}
}
