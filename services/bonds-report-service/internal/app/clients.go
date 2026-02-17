package app

import (
	service "bonds-report-service/internal/application"
	cbr "bonds-report-service/internal/infrastructure/cbr/client"
	cbrtransport "bonds-report-service/internal/infrastructure/cbr/transport"
	moex "bonds-report-service/internal/infrastructure/moex/client"
	moextransport "bonds-report-service/internal/infrastructure/moex/transport"
	"bonds-report-service/internal/infrastructure/sber"
	"bonds-report-service/internal/infrastructure/tinkoffApi/client/analyticsclient"
	"bonds-report-service/internal/infrastructure/tinkoffApi/client/instrumentsclient"
	"bonds-report-service/internal/infrastructure/tinkoffApi/client/portfolioclient"
	"log/slog"

	config "bonds-report-service/internal/configs"
	tinkofftransport "bonds-report-service/internal/infrastructure/tinkoffApi/transport"
)

func InitCBRClient(logger *slog.Logger, host string) *cbr.Client {
	logger.Info("initialize CBR client", slog.String("address", host))
	if host == "" {
		panic("cbr host is empty")
	}
	transport := cbrtransport.NewTransport(logger, host)
	client := cbr.NewCbrClient(logger, transport)
	return client
}

func InitTinkoffApiClient(logger *slog.Logger, host string) *service.TinkoffClients {
	logger.Info("initialize Tinkoff client", slog.String("address", host))
	if host == "" {
		panic("tinkoff host is empty")
	}
	transport := tinkofftransport.NewTransport(logger, host)
	analyticsclient := analyticsclient.NewAnalyticsTinkoffClient(logger, transport)
	instrumentsclient := instrumentsclient.NewInstrumentsTinkoffClient(logger, transport)
	portfolioclient := portfolioclient.NewPortfolioTinkoffClient(logger, transport)
	client := service.NewTinkoffClients(instrumentsclient, portfolioclient, analyticsclient)
	return client
}

func InitTiMoexClient(logger *slog.Logger, host string) *moex.Client {
	logger.Info("initialize Moex client", slog.String("address", host))
	if host == "" {
		panic("moex host is empty")
	}
	transport := moextransport.NewTransport(logger, host)
	client := moex.NewMoexClient(logger, transport)
	return client
}

func InitSberClient(logger *slog.Logger, conf *config.Config) (*sber.Client, error) {
	logger.Info("initialize Sber client", slog.String("address", conf.SberConfigPath))
	return sber.NewClient(conf.RootPath, conf.SberConfigPath)
}
