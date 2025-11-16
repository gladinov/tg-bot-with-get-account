package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"
	"tinkoffApi/internal/configs"
	"tinkoffApi/internal/hanlders"
	"tinkoffApi/internal/service"

	"tinkoffApi/pkg/app"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	app.MustInitialize()
	rootPath := app.MustGetRoot()

	// TODO: Подключить Redis и забирать токен оттуда
	cnfgs := configs.MustInitConfigs(rootPath)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	zapConfig.EncoderConfig.TimeKey = "time"
	l, err := zapConfig.Build()
	logger := l.Sugar()
	defer func() {
		err := logger.Sync()
		if err != nil {
			log.Print(err.Error())
		}
	}()
	if err != nil {
		log.Fatalf("logger creating error %v", err)
	}

	analyticsService := service.NewAnalyticsServiceClient(cnfgs.TinkoffApiConfig, logger)
	portfolioService := service.NewPortfolioServiceClient(cnfgs.TinkoffApiConfig, logger)
	instrumentService := service.NewInstrumentServiceClient(cnfgs.TinkoffApiConfig, logger)

	serviceClient := service.NewService(
		analyticsService,
		portfolioService,
		instrumentService)

	handlrs := hanlders.NewHandlers(serviceClient)

	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

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
			logger.Error("Failed to shutdown server:", err)
		}
	}()
	e.Start(cnfgs.Config.HttpServer)

}
