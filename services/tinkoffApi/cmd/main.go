package main

import (
	"context"
	"log/slog"
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
	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(handlrs.ContextHeaderTraceIdMiddleWare)
	e.Use(handlrs.LoggerMiddleWare)
	e.Use(handlrs.CheckTokenFromRedisByChatIDMiddleWare)

	e.GET("/tinkoff/checktoken", handlrs.CheckToken, handlrs.CheckTokenFromHeadersMiddleWare)
	e.GET("/tinkoff/accounts", handlrs.GetAccounts)
	e.POST("/tinkoff/portfolio", handlrs.GetPortfolio)
	e.POST("/tinkoff/operations", handlrs.GetOperations)
	e.GET("/tinkoff/allassetsuid", handlrs.GetAllAssetUids)
	e.POST("/tinkoff/future", handlrs.GetFutureBy)
	e.POST("/tinkoff/bond", handlrs.GetBondBy)
	e.POST("/tinkoff/currency", handlrs.GetCurrencyBy)
	e.POST("/tinkoff/share/currency", handlrs.GetShareCurrencyBy)
	e.POST("/tinkoff/findby", handlrs.FindBy)
	e.POST("/tinkoff/bondactions", handlrs.GetBondsActions)
	e.POST("/tinkoff/lastprice", handlrs.GetLastPriceInPersentageToNominal)

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := e.Shutdown(shutdownCtx); err != nil {
			logg.Error("Failed to shutdown server:", slog.String("error", err.Error()))
		}
	}()
	address := confs.Config.GetTinkoffAppAddress()
	logg.Info("run tinkoffApiApp", slog.String("address", address))
	e.Start(address)
}
