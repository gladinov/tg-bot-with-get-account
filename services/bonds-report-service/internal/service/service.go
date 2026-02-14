package service

import (
	"bonds-report-service/internal/clients/sber"
	"bonds-report-service/internal/models/domain"
	"bonds-report-service/internal/models/domain/report"
	"context"
	"log/slog"
	"time"

	service_storage "bonds-report-service/internal/repository"
)

const (
	layoutTime = "2006-01-02_15-04-05"
	reportPath = "service/visualization/tables/"
)

const (
	bond     = "bond"
	share    = "share"
	futures  = "futures"
	etf      = "etf"
	currency = "currency"
)

const (
	rub       = "rub"
	cny       = "cny"
	usd       = "usd"
	eur       = "eur"
	hkd       = "hkd"
	futuresPt = "pt."
)

const (
	commodityType = "TYPE_COMMODITY"
	currencyType  = "TYPE_CURRENCY"
	securityType  = "TYPE_SECURITY"
	indexType     = "TYPE_INDEX"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=ReportProcessor
type ReportProcessor interface {
	ProcessOperations(ctx context.Context, reportLine *domain.ReportLine) (_ *report.ReportPositions, err error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=UidProvider
type UidProvider interface {
	GetUid(ctx context.Context, instrumentUid string) (string, error)
	UpdateAndGetUid(ctx context.Context, instrumentUid string) (string, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=CbrClient
type CbrClient interface {
	GetAllCurrencies(ctx context.Context, date time.Time) (res domain.CurrenciesCBR, err error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=MoexClient
type MoexClient interface {
	GetSpecifications(ctx context.Context, ticker string, date time.Time) (data domain.ValuesMoex, err error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=TinkoffInstrumentsClient
type TinkoffInstrumentsClient interface {
	FindBy(ctx context.Context, query string) (domain.InstrumentShortList, error)
	GetBondByUid(ctx context.Context, uid string) (domain.Bond, error)
	GetCurrencyBy(ctx context.Context, figi string) (domain.Currency, error)
	GetFutureBy(ctx context.Context, figi string) (domain.Future, error)
	GetShareCurrencyBy(ctx context.Context, figi string) (domain.ShareCurrency, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=TinkoffPortfolioClient
type TinkoffPortfolioClient interface {
	GetAccounts(ctx context.Context) (_ map[string]domain.Account, err error)
	GetPortfolio(ctx context.Context, accountID string, accountStatus int64) (domain.Portfolio, error)
	GetOperations(ctx context.Context, accountId string, date time.Time) (_ []domain.Operation, err error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=TinkoffAnalyticsClient
type TinkoffAnalyticsClient interface {
	GetLastPriceInPersentageToNominal(ctx context.Context, instrumentUid string) (domain.LastPrice, error)
	GetAllAssetUids(ctx context.Context) (map[string]string, error)
	GetBondsActions(ctx context.Context, instrumentUid string) (domain.BondIdentIdentifiers, error)
}

type TinkoffClients struct {
	Instruments TinkoffInstrumentsClient
	Portfolio   TinkoffPortfolioClient
	Analytics   TinkoffAnalyticsClient
}

func NewTinkoffClients(
	instruments TinkoffInstrumentsClient,
	portfolio TinkoffPortfolioClient,
	analytics TinkoffAnalyticsClient,
) *TinkoffClients {
	return &TinkoffClients{
		Instruments: instruments,
		Portfolio:   portfolio,
		Analytics:   analytics,
	}
}

type ExternalApis struct {
	Moex MoexClient
	Cbr  CbrClient
	Sber *sber.Client
}

func NewExternalApis(
	moex MoexClient,
	cbr CbrClient,
	sber *sber.Client,
) *ExternalApis {
	return &ExternalApis{
		Moex: moex,
		Cbr:  cbr,
		Sber: sber,
	}
}

type Service struct {
	logger          *slog.Logger
	Tinkoff         *TinkoffClients
	External        *ExternalApis
	Storage         service_storage.Storage
	UidProvider     UidProvider
	ReportProcessor ReportProcessor
	now             func() time.Time
}

func NewService(
	logger *slog.Logger,
	tinkoffClients *TinkoffClients,
	externalApis *ExternalApis,
	storage service_storage.Storage,
	reportProcessor ReportProcessor,
	uidProvider UidProvider,
) *Service {
	return &Service{
		logger:          logger,
		Tinkoff:         tinkoffClients,
		External:        externalApis,
		Storage:         storage,
		UidProvider:     uidProvider,
		ReportProcessor: reportProcessor,
		now:             time.Now,
	}
}
