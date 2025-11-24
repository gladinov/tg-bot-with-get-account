package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	postgresMigrator "main.go/migrators/postgres"
	sqliteMigrator "main.go/migrators/sqlite"

	"main.go/clients/cbr"
	"main.go/clients/moex"
	"main.go/clients/sber"
	tgClient "main.go/clients/telegram"
	tinkoffapi "main.go/clients/tinkoffApi"
	event_consumer "main.go/consumer/event-consumer"
	"main.go/events/telegram"
	"main.go/internal/config"
	"main.go/lib/cryptoToken"
	pathwd "main.go/lib/pathWD"
	loggAdapter "main.go/logger"
	"main.go/pkg/app"
	"main.go/service"
	storageInterface "main.go/service/storage"
	"main.go/service/storage/postgreSQL"
	servicet_sqlite "main.go/service/storage/sqlite"
	"main.go/storage"
	"main.go/storage/postgres"
	"main.go/storage/sqlite"
)

const (
	moexHost        = "localhost:8081"
	tinkoffApiHost  = "localhost:8082"
	cbrHost         = "localhost:8083"
	tgBotHost       = "api.telegram.org"
	storageSqlPath  = "/data/sqlite/storage.db"
	sberConfigPath  = "/configs/sber.yaml"
	redisConfigPath = "/configs/redisConfig.yaml"
	batchSize       = 100
)

const (
	postreSQL = "postgreSQL"
	SQLite    = "SQLite"
)

func main() {
	app.MustInitialize()
	rootPath := app.MustGetRoot()

	cnfg := config.MustInitConfig(rootPath)
	userStorageConfig := config.MustInitStorageConfig(rootPath)
	
	//redisConfig := config.MustInitRedisConfig(rootPath, redisConfigPath)

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

	key := os.Getenv("KEY")
	if key == "" {
		log.Fatal("KEY environment variable is required")
	}

	tokenCrypter := cryptoToken.NewTokenCrypter(key)

	telegrammClient := tgClient.New(tgBotHost, token)

	logger := loggAdapter.SetupLogger()

	moexApi := moex.NewClient(moexHost)

	cbrApi := cbr.New(cbrHost)

	sberClient, err := sber.NewClient(rootPath, sberConfigPath)
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

	userStorage, err := NewStorage(ctx, userStorageConfig, rootPath)
	if err != nil {
		logger.Fatalf("can't create user_storage: %s", err.Error())
	}
	defer func() {
		if userStorage != nil {
			userStorage.CloseDB()
		}
	}()

	serviceClient := service.New(
		tinkoffApiClient,
		moexApi,
		cbrApi,
		sberClient,
		serviceStorage)

	processor := telegram.NewProccesor(
		tokenCrypter,
		telegrammClient,
		userStorage,
		serviceClient,
	)

	fetcher := telegram.NewFetcher(telegrammClient)

	logger.Printf("service started")

	consumer := event_consumer.New(fetcher, processor, batchSize)

	if err := consumer.Start(); err != nil {
		logger.Fatalf("service is stopped")
	}
}

func NewServiceStorage(ctx context.Context, config config.Config, rootPath string) (storageInterface.Storage, error) {
	switch config.DbType {
	case postreSQL:
		serviceStorage, err := postgreSQL.NewStorage(config)
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

func NewStorage(ctx context.Context, config config.UserStorageConfig, rootPath string) (storage.Storage, error) {
	switch config.DbType {
	case postreSQL:
		storage, err := postgres.NewStorage(config)
		if err != nil {
			return nil, err
		}
		postgresMigrator.MustMigratePostgres(rootPath, config)

		err = storage.Init(ctx)
		if err != nil {
			return nil, err
		}

		return storage, nil
	case SQLite:
		storageAbsolutPath, err := pathwd.PathFromWD(rootPath, storageSqlPath)
		if err != nil {
			return nil, err
		}
		storage, err := sqlite.New(storageAbsolutPath)
		if err != nil {
			return nil, err
		}
		sqliteMigrator.MustMigrateSqllite(rootPath, storageAbsolutPath, config.MigrationsSqllitePath)

		err = storage.Init(ctx)
		if err != nil {
			return nil, err
		}

		return storage, nil
	default:
		return nil, errors.New("Possible init only SQLite or PostgreSQL databases")
	}
}
