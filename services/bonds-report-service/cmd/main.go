package main

import (
	"bonds-report-service/internal/app"
	"bonds-report-service/internal/application/ports"
	"bonds-report-service/internal/application/usecases"
	config "bonds-report-service/internal/configs"
	"bonds-report-service/internal/handlers"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	sl "github.com/gladinov/mylogger"
	"github.com/gladinov/traceidgenerator"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	conf := config.MustInitConfig()

	logg := sl.NewLogger(conf.Env)

	logg.Info("start app",
		slog.String("env", conf.Env),
		slog.String("bond-report-service_app_host", conf.Clients.BondReportService.Host),
		slog.String("bond-report-service_app_port", conf.Clients.BondReportService.Port))

	_ = traceidgenerator.Must()

	repo := app.MustInitNewStorage(ctx, conf, logg)

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

	dividerByAssetType := app.InitDividerByAssetType(logg, tinkoffApiHelper, cbrCurrencyGetter, conf.WorkersNubmer)

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
		dividerByAssetType,
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

	httpSrv := &http.Server{
		Addr:    address,
		Handler: router,
	}

	errCh := make(chan error, 1)

	go func() {
		logg.Info("run bond-report-service", slog.String("address", address))
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()
	select {
	case <-ctx.Done():
		logg.InfoContext(ctx, "Shutdown signal received")
	case err = <-errCh:
		logg.ErrorContext(ctx, "server stopped with error", slog.Any("error", err))
	}
	gracefulShutdown(ctx, logg, httpSrv, repo)
}

func gracefulShutdown(ctx context.Context, logg *slog.Logger, httpSrv *http.Server, repo ports.Storage) {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		logg.ErrorContext(ctx, "Forced shutdown", slog.Any("err", err))
	}
	logg.InfoContext(ctx, "close DB")
	repo.CloseDB()
	logg.InfoContext(ctx, "Server exited gracefully")
}
