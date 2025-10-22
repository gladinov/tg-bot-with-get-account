package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
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
	// moexHost  = "iss.moex.com"
	moexHost  = "localhost:8081"
	cbrHost   = "www.cbr.ru"
	tgBotHost = "api.telegram.org"
	// storagePath            = "storage"
	storageSqlPath         = "/data/sqlite/storage.db"
	service_storageSqlPath = "/data/sqlite/service_storage.db"
	batchSize              = 100
)

func main() {
	// //  for local
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading .env file. Erorr: %v", err.Error())
	}

	token := os.Getenv("LOCAL_BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN environment variable is required")
	}

	telegrammClient := tgClient.New(tgBotHost, token)

	logger := loggAdapter.SetupLogger()

	moexApi := moex.NewClient(moexHost)

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
