package portfolioclient

import (
	"bonds-report-service/internal/clients/tinkoffApi/transport"
	"bonds-report-service/internal/models/domain"
	"context"
	"log/slog"
	"time"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=PortfolioClient
type PortfolioClient interface {
	GetAccounts(ctx context.Context) (_ map[string]domain.Account, err error)
	GetPortfolio(ctx context.Context, accountID string, accountStatus int64) (domain.Portfolio, error)
	GetOperations(ctx context.Context, accountId string, date time.Time) (_ []domain.Operation, err error)
}

type PortfolioTinkoffClient struct {
	logger    *slog.Logger
	transport transport.TransportClient
}

func NewPortfolioTinkoffClient(logger *slog.Logger, transport transport.TransportClient) *PortfolioTinkoffClient {
	return &PortfolioTinkoffClient{
		logger:    logger,
		transport: transport,
	}
}
