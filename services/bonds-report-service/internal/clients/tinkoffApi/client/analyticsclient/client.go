package analyticsclient

import (
	"bonds-report-service/internal/clients/tinkoffApi/transport"
	"bonds-report-service/internal/models/domain"
	"context"
	"log/slog"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=AnalyticsClient
type AnalyticsClient interface {
	GetLastPriceInPersentageToNominal(ctx context.Context, instrumentUid string) (domain.LastPrice, error)
	GetAllAssetUids(ctx context.Context) (map[string]string, error)
	GetBondsActions(ctx context.Context, instrumentUid string) (domain.BondIdentIdentifiers, error)
}

type AnalyticsTinkoffClient struct {
	logger    *slog.Logger
	transport transport.TransportClient
}

func NewAnalyticsTinkoffClient(logger *slog.Logger, transport transport.TransportClient) *AnalyticsTinkoffClient {
	return &AnalyticsTinkoffClient{
		logger:    logger,
		transport: transport,
	}
}
