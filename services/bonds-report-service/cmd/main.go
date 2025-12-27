package main

import (
	"bonds-report-service/clients/cbr"
	"bonds-report-service/clients/moex"
	"bonds-report-service/clients/sber"
	"bonds-report-service/clients/tinkoffApi"
	config "bonds-report-service/internal/configs"
	"bonds-report-service/internal/handlers"
	"bonds-report-service/internal/repository"
	"bonds-report-service/internal/service"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	sl "github.com/gladinov/mylogger"
)

func main() {
	conf := config.MustInitConfig()

	logg := sl.NewLogger(conf.Env)

	logg.Info("start app",
		slog.String("env", conf.Env),
		slog.String("bond-report-service_app_host", conf.Clients.BondReportService.Host),
		slog.String("bond-report-service_app_port", conf.Clients.BondReportService.Port))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	repo := repository.MustInitNewStorage(ctx, conf, logg)

	logg.Info("initialize Tinkoff client", slog.String("addres", conf.Clients.TinkoffClient.GetTinkoffApiAddress()))
	tinkoffClient := tinkoffApi.NewClient(logg, conf.Clients.TinkoffClient.GetTinkoffApiAddress())

	logg.Info("initialize Moex client", slog.String("addres", conf.Clients.MoexClient.GetMoexAppAddress()))
	moexClient := moex.NewClient(logg, conf.Clients.MoexClient.GetMoexAppAddress())

	logg.Info("initialize CBR client", slog.String("addres", conf.Clients.CBRClient.GetCBRAppAddress()))
	cbrClient := cbr.New(logg, conf.Clients.CBRClient.GetCBRAppAddress())

	logg.Info("initialize Sber client", slog.String("addres", conf.SberConfigPath))
	sberClient, err := sber.NewClient(conf.RootPath, conf.SberConfigPath)
	if err != nil {
		logg.Error("could not create sber client", slog.String("error", err.Error()))
		return
	}

	logg.Info("initialize Service client")
	serviceClient := service.New(
		logg,
		tinkoffClient,
		moexClient,
		cbrClient,
		sberClient,
		repo)

	logg.Info("initialize Handlers")
	handl := handlers.NewHandlers(logg, serviceClient)

	logg.Info("initialize router gin")
	router := gin.New()
	router.Use(gin.Recovery())

	router.Use(handl.ContextHeaderTraceIdMiddleWare())
	router.Use(handl.LoggerMiddleware())
	router.Use(handl.AuthMiddleware())

	router.GET("/bondReportService/accounts", handl.GetAccountsList)
	router.GET("/bondReportService/getBondReportsByFifo", handl.GetBondReportsByFifo)
	router.GET("/bondReportService/getUSD", handl.GetUSD)
	router.GET("/bondReportService/getBondReports", handl.GetBondReports)
	router.GET("/bondReportService/getPortfolioStructure", handl.GetPortfolioStructure)
	router.GET("/bondReportService/getUnionPortfolioStructure", handl.GetUnionPortfolioStructure)
	router.GET("/bondReportService/getUnionPortfolioStructureWithSber", handl.GetUnionPortfolioStructureWithSber)

	address := conf.Clients.BondReportService.GetBondReportServiceAppAddress()
	logg.Info("run bond-report-service", slog.String("address", address))
	router.Run(address)
}
