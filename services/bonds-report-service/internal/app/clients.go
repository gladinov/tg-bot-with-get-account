package app

import (
	"bonds-report-service/internal/clients/cbr"
	"bonds-report-service/internal/clients/moex"
	"bonds-report-service/internal/clients/sber"
	"bonds-report-service/internal/clients/tinkoffApi"
	config "bonds-report-service/internal/configs"
	"log/slog"
)

func InitCBRClient(logger *slog.Logger, host string) *cbr.Client {
	logger.Info("initialize CBR client", slog.String("addres", host))
	transport := cbr.NewTransport(logger, host)
	client := cbr.NewCbrClient(logger, transport)
	return client
}

func InitTinkoffApiClient(logger *slog.Logger, host string) *tinkoffApi.Client {
	logger.Info("initialize Tinkoff client", slog.String("addres", host))
	client := tinkoffApi.NewClient(logger, host)
	return client
}

func InitTiMoexClient(logger *slog.Logger, host string) *moex.Client {
	logger.Info("initialize Moex client", slog.String("addres", host))
	client := moex.NewClient(logger, host)
	return client
}

func InitSberClient(logger *slog.Logger, conf *config.Config) (*sber.Client, error) {
	logger.Info("initialize Sber client", slog.String("addres", conf.SberConfigPath))
	return sber.NewClient(conf.RootPath, conf.SberConfigPath)
}
