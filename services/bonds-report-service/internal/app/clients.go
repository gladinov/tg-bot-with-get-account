package app

import (
	cbr "bonds-report-service/internal/clients/cbr/client"
	cbrtransport "bonds-report-service/internal/clients/cbr/transport"
	moex "bonds-report-service/internal/clients/moex/client"
	moextransport "bonds-report-service/internal/clients/moex/transport"
	"bonds-report-service/internal/clients/sber"
	tinkoffApi "bonds-report-service/internal/clients/tinkoffApi/client"
	"bonds-report-service/internal/clients/tinkoffApi/client/analyticsclient"
	"bonds-report-service/internal/clients/tinkoffApi/client/instrumentsclient"
	"bonds-report-service/internal/clients/tinkoffApi/client/portfolioclient"
	tinkofftransport "bonds-report-service/internal/clients/tinkoffApi/transport"
	config "bonds-report-service/internal/configs"
	"log/slog"
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

func InitTinkoffApiClient(logger *slog.Logger, host string) *tinkoffApi.Client {
	logger.Info("initialize Tinkoff client", slog.String("address", host))
	if host == "" {
		panic("tinkoff host is empty")
	}
	transport := tinkofftransport.NewTransport(logger, host)
	analyticsclient := analyticsclient.NewAnalyticsTinkoffClient(logger, transport)
	instrumentsclient := instrumentsclient.NewInstrumentsTinkoffClient(logger, transport)
	portfolioclient := portfolioclient.NewPortfolioTinkoffClient(logger, transport)
	client := tinkoffApi.NewTinkoffClient(
		logger,
		instrumentsclient,
		portfolioclient,
		analyticsclient)
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
