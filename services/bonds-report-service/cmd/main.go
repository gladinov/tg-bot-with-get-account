package main

import (
	"bonds-report-service/internal/app"
	"bonds-report-service/internal/application/usecases"
	config "bonds-report-service/internal/configs"
	"bonds-report-service/internal/handlers"
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

	repo := app.MustInitNewStorage(ctx, conf, logg)
	defer repo.CloseDB()

	tinkoffApiHelper := app.InitTinkoffApiHelper(logg, conf.Clients.TinkoffClient.GetTinkoffApiAddress())

	moexClient := app.InitTiMoexClient(logg, conf.Clients.MoexClient.GetMoexAppAddress())

	cbrClient := app.InitCBRClient(logg, conf.Clients.CBRClient.GetCBRAppAddress())

	sberClient, err := app.InitSberClient(logg, &conf)
	if err != nil {
		logg.Error("could not create sber client", slog.String("error", err.Error()))
		return
	}

	bondReporter := app.InitBondReportProcessor(logg)

	cbrCurrencyGetter := app.InitCBRCurrencyGetter(logg, cbrClient, repo)

	generalBondReporter := app.InitGeneralReportProcessor(logg)

	moexSpecificationGetter := app.InitMoexSpecificationGetter(logg, moexClient)

	reportProcessor := app.InitReportProcessor(logg)

	uidProvider := app.InitUidProvider(logg, repo, tinkoffApiHelper.Analytics)

	operationsUpdater := app.InitOperationsUpdater(logg, tinkoffApiHelper, repo)

	positionProcessor := app.InitPositionProcessor(logg, uidProvider)

	reportLineBuilder := app.InitReportLineBuilder(logg, tinkoffApiHelper, cbrCurrencyGetter)

	dividerbyassettype := app.InitDividerByAssetType(logg, tinkoffApiHelper, cbrCurrencyGetter)

	externalApis := usecases.NewExternalApis(moexClient, cbrClient, sberClient)

	helpers := usecases.NewHelpers(bondReporter,
		cbrCurrencyGetter,
		generalBondReporter,
		moexSpecificationGetter,
		reportProcessor,
		tinkoffApiHelper,
		operationsUpdater,
		positionProcessor,
		reportLineBuilder,
		dividerbyassettype,
	)

	logg.Info("initialize Service client")
	serviceClient := usecases.NewService(
		logg,
		conf.WorkersNubmer,
		externalApis,
		helpers,
		repo,
	)

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
