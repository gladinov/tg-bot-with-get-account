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

	tinkoffApiClient := service.New(ctx, logger, cnfgs.TinkoffApiConfig)
	// TODO: Подключить Redis и забирать токен оттуда

	handlrs := hanlders.NewHandlers(tinkoffApiClient)

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
	e.POST("/tinkoff/basesharecurrency", handlrs.GetBaseShareFutureValute)
	e.POST("/tinkoff/findby", handlrs.FindBy)
	e.POST("/tinkoff/bondactions", handlrs.GetBondsActions)
	e.POST("/tinkoff/lastprice", handlrs.GetLastPriceInPersentageToNominal)

	e.Start(cnfgs.Config.HttpServer)

}
