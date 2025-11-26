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
	"bonds-report-service/pkg/app"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	app.MustInitialize()
	rootPath := app.MustGetRoot()
	conf := config.MustInitConfig(rootPath)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	repo := repository.MustInitNewStorage(ctx, conf, rootPath)

	tinkoffClient := tinkoffApi.NewClient(conf.TinkoffApiHost)

	moexClient := moex.NewClient(conf.MoexHost)

	cbrClient := cbr.New(conf.CBRHost)

	sberClient, err := sber.NewClient(rootPath, conf.SberConfigPath)
	if err != nil {
		log.Fatalf("could not create sber client: %s", err.Error())
	}

	serviceClient := service.New(tinkoffClient, moexClient, cbrClient, sberClient, repo)

	handl := handlers.NewHandlers(serviceClient)

	router := gin.Default()

	router.Use(gin.Logger())

	router.GET("/bondReportService/accounts", handlers.AuthMiddleware(), handl.GetAccountsList)
	router.POST("/bondReportService/getBondReportsByFifo", handlers.AuthMiddleware(), handl.GetBondReportsByFifo)
	router.GET("/bondReportService/getUSD", handlers.AuthMiddleware(), handl.GetUSD)
	router.POST("/bondReportService/getBondReports", handlers.AuthMiddleware(), handl.GetBondReports)
	router.GET("/bondReportService/getPortfolioStructure", handlers.AuthMiddleware(), handl.GetPortfolioStructure)
	router.GET("/bondReportService/getUnionPortfolioStructure", handlers.AuthMiddleware(), handl.GetUnionPortfolioStructure)
	router.GET("/bondReportService/getUnionPortfolioStructureWithSber", handlers.AuthMiddleware(), handl.GetUnionPortfolioStructureWithSber)

	router.Run(conf.BondReportServiceHost)

}
