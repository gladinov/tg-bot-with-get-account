package service

import (
	"bonds-report-service/internal/infrastructure/sber"
	"log/slog"
	"time"
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
	logger                     *slog.Logger
	Tinkoff                    *TinkoffClients
	External                   *ExternalApis
	Storage                    Storage
	UidProvider                UidProvider
	ReportProcessor            ReportProcessor
	BondReportProcessor        BondReportProcessor
	GeneralBondReportProcessor GeneralBondReportProcessor
	now                        func() time.Time
}

func NewService(
	logger *slog.Logger,
	tinkoffClients *TinkoffClients,
	externalApis *ExternalApis,
	storage Storage,
	uidProvider UidProvider,
	reportProcessor ReportProcessor,
	bondReportProcessor BondReportProcessor,
	generalBondReportProcessor GeneralBondReportProcessor,
) *Service {
	return &Service{
		logger:                     logger,
		Tinkoff:                    tinkoffClients,
		External:                   externalApis,
		Storage:                    storage,
		UidProvider:                uidProvider,
		ReportProcessor:            reportProcessor,
		BondReportProcessor:        bondReportProcessor,
		GeneralBondReportProcessor: generalBondReportProcessor,
		now:                        time.Now,
	}
}
