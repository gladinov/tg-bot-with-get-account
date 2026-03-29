package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"tinkoffApi/internal/configs"
	"tinkoffApi/internal/handlers"
	redisClient "tinkoffApi/internal/repository/redis"
	"tinkoffApi/internal/service"
	loggeradapter "tinkoffApi/lib/logger/loggerAdapter"

	"github.com/gladinov/cryptotoken"

	sl "github.com/gladinov/mylogger"
	"github.com/gladinov/traceidgenerator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	_ = traceidgenerator.Must()

	confs := configs.MustInitConfigs()

	logg := sl.NewLogger(confs.Config.Env)

	logg.Info("start app",
		slog.String("env", confs.Config.Env),
		slog.String("cbr_app_host", confs.Config.TinkoffApiAppHost),
		slog.String("cbr_app_port", confs.Config.TinkoffApiAppPort))

	logg.Info("initialize logger adapter")
	loggAdapter := loggeradapter.NewLoggerAdapter(logg)

	logg.Info("initialize analyticsService")
	analyticsService := service.NewAnalyticsServiceClient(confs.TinkoffApiConfig, loggAdapter)
	logg.Info("initialize portfolioService")
	portfolioService := service.NewPortfolioServiceClient(confs.TinkoffApiConfig, loggAdapter)
	logg.Info("initialize instrumentService")
	instrumentService := service.NewInstrumentServiceClient(confs.TinkoffApiConfig, loggAdapter)

	logg.Info("initialize serviceClient")
	serviceClient := service.NewService(
		analyticsService,
		portfolioService,
		instrumentService)

	logg.Info("initialize tokenCrypter")
	tokenCrypter := cryptotoken.NewTokenCrypter(confs.Config.Key)

	logg.Info("initialize redis", slog.String("adress", confs.Config.RedisHTTPServer.GetAddress()))
	redis, err := redisClient.NewClient(ctx, confs.Config)
	if err != nil {
		logg.Error("haven't connect with redis", slog.String("error", err.Error()))
	}

	logg.Info("initialize handlers")
	handlrs := handlers.NewHandlers(logg, serviceClient, tokenCrypter, redis)

	logg.Info("initialize router echo")
	router := echo.New()

	router.Use(middleware.CORS())
	router.Use(handlrs.ContextHeaderTraceIdMiddleWare)
	router.Use(handlrs.LoggerMiddleWare)
	router.Use(handlrs.CheckTokenFromRedisByChatIDMiddleWare)

	router.GET("/tinkoff/checktoken", handlrs.CheckToken, handlrs.CheckTokenFromHeadersMiddleWare)
	router.GET("/tinkoff/accounts", handlrs.GetAccounts)
	router.POST("/tinkoff/portfolio", handlrs.GetPortfolio)
	router.POST("/tinkoff/operations", handlrs.GetOperations)
	router.GET("/tinkoff/allassetsuid", handlrs.GetAllAssetUids)
	router.POST("/tinkoff/future", handlrs.GetFutureBy)
	router.POST("/tinkoff/bond", handlrs.GetBondBy)
	router.POST("/tinkoff/currency", handlrs.GetCurrencyBy)
	router.POST("/tinkoff/share/currency", handlrs.GetShareCurrencyBy)
	router.POST("/tinkoff/findby", handlrs.FindBy)
	router.POST("/tinkoff/bondactions", handlrs.GetBondsActions)
	router.POST("/tinkoff/lastprice", handlrs.GetLastPriceInPersentageToNominal)

	address := confs.Config.GetTinkoffAppAddress()

	httpSrv := &http.Server{
		Addr:         address,
		Handler:      router,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  60 * time.Second,
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
	gracefulShutdown(ctx, logg, httpSrv)
}

func gracefulShutdown(ctx context.Context, logg *slog.Logger, httpSrv *http.Server) {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		logg.ErrorContext(ctx, "Forced shutdown", slog.Any("err", err))
	}
	logg.InfoContext(ctx, "Server exited gracefully")
}
