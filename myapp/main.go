package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"main.go/clients/cbr"
	"main.go/clients/moex"
	"main.go/clients/sber"
	tgClient "main.go/clients/telegram"
	tinkoffapi "main.go/clients/tinkoffApi"
	event_consumer "main.go/consumer/event-consumer"
	"main.go/events/telegram"
	"main.go/internal/config"
	pathwd "main.go/lib/pathWD"
	loggAdapter "main.go/logger"
	"main.go/pkg/app"
	"main.go/service"
	storageInterface "main.go/service/storage"
	"main.go/service/storage/postgreSQL"
	servicet_sqlite "main.go/service/storage/sqlite"
	"main.go/storage/sqlite"
)

const (
	moexHost       = "localhost:8081"
	tinkoffApiHost = "localhost:8082"
	cbrHost        = "localhost:8083"
	tgBotHost      = "api.telegram.org"
	storageSqlPath = "/data/sqlite/storage.db"
	sberConfigPath = "/configs/sber.yaml"
	batchSize      = 100
)

const (
	postreSQL = "postgreSQL"
	SQLite    = "SQLite"
)

func main() {
	app.MustInitialize()
	rootPath := app.MustGetRoot()

	cnfg := config.MustInitConfig(rootPath)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

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

	sber, err := sber.NewClient(rootPath, sberConfigPath)
	if err != nil {
		logger.Fatalf("create sber client failed: %s", err.Error())
	}
	tinkoffApiClient := tinkoffapi.NewClient(tinkoffApiHost)

	serviceStorage, err := NewServiceStorage(ctx, cnfg, rootPath)
	if err != nil {
		logger.Fatalf("can't create service_storage: %s", err.Error())
	}
	defer func() {
		if serviceStorage != nil {
			serviceStorage.CloseDB()
		}
	}()
	// TODO: Change to REdis
	storageAbsolutPath, err := pathwd.PathFromWD(rootPath, storageSqlPath)
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
	// TODO: Change to REdis

	serviceClient := service.New(
		tinkoffApiClient,
		moexApi,
		cbrApi,
		sber,
		serviceStorage)

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

func NewServiceStorage(ctx context.Context, config config.Config, rootPath string) (storageInterface.Storage, error) {
	switch config.DbType {
	case postreSQL:
		serviceStorage, err := postgreSQL.NewStorage()
		if err != nil {
			return nil, err
		}
		err = serviceStorage.InitDB(ctx)
		if err != nil {
			return nil, err
		}
		return serviceStorage, nil

	case SQLite:
		serviceStorageAbsolutPath, err := pathwd.PathFromWD(rootPath, config.ServiceStorageSQLLitePath)
		if err != nil {
			return nil, err
		}

		serviceStorage, err := servicet_sqlite.New(serviceStorageAbsolutPath)
		if err != nil {
			return nil, err
		}

		if err := serviceStorage.Init(ctx); err != nil {
			return nil, err
		}
		return serviceStorage, nil
	default:
		return nil, errors.New("Possible init only SQLite or PostgreSQL databases")
	}
}
