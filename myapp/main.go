package main

import (
	"context"
	"flag"
	"log"
	"os"

	"main.go/clients/cbr"
	"main.go/clients/moex"
	tgClient "main.go/clients/telegram"
	tinkoffapi "main.go/clients/tinkoffApi"
	event_consumer "main.go/consumer/event-consumer"
	"main.go/events/telegram"
	pathwd "main.go/lib/pathWD"
	loggAdapter "main.go/logger"
	"main.go/service"
	servicet_sqlite "main.go/service/storage/sqlite"
	"main.go/storage/sqlite"
)

const (
	moexHost  = "iss.moex.com"
	cbrHost   = "www.cbr.ru"
	tgBotHost = "api.telegram.org"
	// storagePath            = "storage"
	storageSqlPath         = "/data/sqlite/storage.db"
	service_storageSqlPath = "/data/sqlite/service_storage.db"
	batchSize              = 100
	token                  = "7758843053:AAGSURIkq8xJYio8-m9WCHP9eIDWEqPMu9c"
)

func main() {
	telegrammClient := tgClient.New(tgBotHost, token)

	logger := loggAdapter.SetupLogger()

	// TODO: delete block. begin
	cwd, _ := os.Getwd()

	logger.Infof("work dir path is :%s", cwd)
	// end

	moexApi := moex.New(moexHost)

	cbrApi := cbr.New(cbrHost)

	tinkoffApiClient := tinkoffapi.New(context.TODO(), logger)

	storageAbsolutPath, err := pathwd.PathFromWD(storageSqlPath)
	if err != nil {
		logger.Fatalf("can't create absolute storare path by: %s", storageSqlPath)
	}

	storage, err := sqlite.New(storageAbsolutPath)
	if err != nil {
		logger.Fatalf("can't connect to storage err:%s:", err.Error())
	}

	if err := storage.Init(context.TODO()); err != nil {
		logger.Fatalf("can't init storage ")
	}

	logger.Infof("storage sucsess init in path: %s", storageAbsolutPath)

	service_storageAbsolutPath, err := pathwd.PathFromWD(service_storageSqlPath)
	if err != nil {
		logger.Fatalf("can't create absolute service_storare path by: %s", service_storageSqlPath)
	}

	service_storage, err := servicet_sqlite.New(service_storageAbsolutPath)
	if err != nil {
		logger.Fatalf("can't connect to storage err:%s:", err.Error())
	}

	if err := service_storage.Init(context.TODO()); err != nil {
		logger.Fatalf("can't init storage ")
	}

	logger.Infof("storage sucsess init in path: %s", service_storageAbsolutPath)

	serviceClient := service.New(
		tinkoffApiClient,
		moexApi,
		cbrApi,
		service_storage)

	processor := telegram.NewProccesor(
		telegrammClient,
		storage,
		serviceClient,
	)

	fetcher := telegram.NewFetcher(telegrammClient)

	logger.Printf("service started")

	consumer := event_consumer.New(fetcher, processor, batchSize)

	if err := consumer.Start(); err != nil {
		logger.Fatalf("service is stopped")
	}
}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
