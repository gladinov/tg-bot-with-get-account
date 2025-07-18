package main

import (
	"context"
	"flag"
	"log"

	"main.go/clients/cbr"
	"main.go/clients/moex"
	tgClient "main.go/clients/telegram"
	tinkoffapi "main.go/clients/tinkoffApi"
	event_consumer "main.go/consumer/event-consumer"
	"main.go/events/telegram"
	loggAdapter "main.go/logger"
	"main.go/service"
	servicet_sqlite "main.go/service/storage/sqlite"
	"main.go/storage/sqlite"
)

const (
	moexHost               = "iss.moex.com"
	cbrHost                = "www.cbr.ru"
	tgBotHost              = "api.telegram.org"
	storagePath            = "storage"
	storageSqlPath         = "data/sqlite/storage.db"
	service_storageSqlPath = "data/sqlite/service_storage.db"
	batchSize              = 100
)

func main() {
	telegrammClient := tgClient.New(tgBotHost, mustToken())

	logger := loggAdapter.SetupLogger()

	moexApi := moex.New(moexHost)

	cbrApi := cbr.New(cbrHost)

	tinkoffApiClient := tinkoffapi.New(context.TODO(), logger)

	storage, err := sqlite.New(storageSqlPath)
	if err != nil {
		logger.Fatalf("can't connect to storage:")
	}

	if err := storage.Init(context.TODO()); err != nil {
		logger.Fatalf("can't init storage ")
	}

	service_storage, err := servicet_sqlite.New(service_storageSqlPath)
	if err != nil {
		logger.Fatalf("can't connect to storage:")
	}

	if err := service_storage.Init(context.TODO()); err != nil {
		logger.Fatalf("can't init storage ")
	}

	serviceClient := service.New(tinkoffApiClient, moexApi, cbrApi, service_storage)

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
