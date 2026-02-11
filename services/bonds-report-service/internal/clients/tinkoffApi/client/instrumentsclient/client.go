package instrumentsclient

import (
	"bonds-report-service/internal/clients/tinkoffApi/transport"
	"bonds-report-service/internal/models/domain"
	"context"
	"log/slog"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=InstrumentsClient
type InstrumentsClient interface {
	FindBy(ctx context.Context, query string) ([]domain.InstrumentShort, error)
	GetBondByUid(ctx context.Context, uid string) (domain.Bond, error)
	GetCurrencyBy(ctx context.Context, figi string) (domain.Currency, error)
	GetFutureBy(ctx context.Context, figi string) (domain.Future, error)
	GetShareCurrencyBy(ctx context.Context, figi string) (domain.ShareCurrency, error)
}

type InstrumentsTinkoffClient struct {
	logger    *slog.Logger
	transport transport.TransportClient
}

func NewInstrumentsTinkoffClient(logger *slog.Logger, transport transport.TransportClient) *InstrumentsTinkoffClient {
	return &InstrumentsTinkoffClient{
		logger:    logger,
		transport: transport,
	}
}
