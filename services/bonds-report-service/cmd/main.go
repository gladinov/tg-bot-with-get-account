package main

import (
	"bonds-report-service/internal/app"
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
	"github.com/gladinov/traceidgenerator"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	conf := config.MustInitConfig()

	logg := sl.NewLogger(conf.Env)

	logg.Info("start app",
		slog.String("env", conf.Env),
		slog.String("bond-report-service_app_host", conf.Clients.BondReportService.Host),
		slog.String("bond-report-service_app_port", conf.Clients.BondReportService.Port))

	_ = traceidgenerator.Must()

	repo := repository.MustInitNewStorage(ctx, conf, logg)
	// TODO: close db

	tinkoffClient := app.InitTinkoffApiClient(logg, conf.Clients.TinkoffClient.GetTinkoffApiAddress())

	moexClient := app.InitTiMoexClient(logg, conf.Clients.MoexClient.GetMoexAppAddress())

	cbrClient := app.InitCBRClient(logg, conf.Clients.CBRClient.GetCBRAppAddress())

	sberClient, err := app.InitSberClient(logg, &conf)
	if err != nil {
		logg.Error("could not create sber client", slog.String("error", err.Error()))
		return
	}

	externalApis := service.NewExternalApis(moexClient, cbrClient, sberClient)

	uidProvider := app.InitUidProvider(logg, repo, tinkoffClient.Analytics)

	logg.Info("initialize Service client")
	serviceClient := service.NewClient(
		logg,
		tinkoffClient,
		externalApis,
		repo,
		uidProvider)

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
