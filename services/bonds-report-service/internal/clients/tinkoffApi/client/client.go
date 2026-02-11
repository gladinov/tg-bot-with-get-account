package tinkoffApi

import (
	"bonds-report-service/internal/clients/tinkoffApi/client/analyticsclient"
	"bonds-report-service/internal/clients/tinkoffApi/client/instrumentsclient"
	"bonds-report-service/internal/clients/tinkoffApi/client/portfolioclient"
	"log/slog"
)

type Client struct {
	logger                   *slog.Logger
	InstrumentsTinkoffClient instrumentsclient.InstrumentsClient
	PortfolioTinkoffClient   portfolioclient.PortfolioClient
	AnalyticsTinkoffClient   analyticsclient.AnalyticsClient
}

func NewTinkoffClient(logger *slog.Logger,
	instrumentsTinkoffClient instrumentsclient.InstrumentsClient,
	portfolioTinkoffClient portfolioclient.PortfolioClient,
	analyticsTinkoffClient analyticsclient.AnalyticsClient,
) *Client {
	return &Client{
		logger:                   logger,
		InstrumentsTinkoffClient: instrumentsTinkoffClient,
		PortfolioTinkoffClient:   portfolioTinkoffClient,
		AnalyticsTinkoffClient:   analyticsTinkoffClient,
	}
}
