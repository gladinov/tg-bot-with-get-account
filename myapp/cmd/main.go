package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	bondreportservice "main.go/clients/bondReportService"
	tgClient "main.go/clients/telegram"
	tinkoffapi "main.go/clients/tinkoffApi"
	event_consumer "main.go/internal/app/consumer/event-consumer"
	"main.go/internal/app/events/telegram"
	"main.go/internal/config"
	storage "main.go/internal/repository"
	"main.go/internal/repository/redis"
	"main.go/lib/cryptoToken"
	loggAdapter "main.go/logger"
)

const (
	batchSize = 100
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	logger := loggAdapter.SetupLogger()
	// for local run in terminal
	//logger := loggAdapter.SetupLogger()
	//app.MustInitialize()
	//rootPath := app.MustGetRoot()

	// docker run
	rootPath := os.Getenv("ROOT_PATH")
	if rootPath == "" {
		logger.Fatalf("ROOT_PATH environment variable is required")
	}

	cnfg := config.MustInitConfig(rootPath)

	//// //  for local
	//envPath := filepath.Join(rootPath, ".env")
	//err := godotenv.Load(envPath)
	//if err != nil {
	//	logger.Printf("Error loading .env file. Erorr: %v", err.Error())
	//}
	//
	token := os.Getenv("LOCAL_BOT_TOKEN")
	if token == "" {
		logger.Fatalf("BOT_TOKEN environment variable is required")
	}

	key := os.Getenv("KEY")
	if key == "" {
		logger.Fatalf("KEY environment variable is required")
	}

	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	if postgresPassword == "" {
		logger.Fatalf("POSTGRES_PASSWORD environment variable is required")
	}

	postgresDB := os.Getenv("POSTGRES_DB")
	if postgresDB == "" {
		logger.Fatalf("POSTGRES_DB environment variable is required")
	}

	postgresUser := os.Getenv("POSTGRES_USER")
	if postgresUser == "" {
		logger.Fatalf("POSTGRES_USER environment variable is required")
	}

	postgresHost := os.Getenv("POSTGRES_HOST")
	if postgresHost == "" {
		logger.Fatalf("POSTGRES_HOST environment variable is required")
	}

	cnfg.PostgresHost.Host = postgresHost
	cnfg.PostgresHost.Password = postgresPassword
	cnfg.PostgresHost.Dbname = postgresDB
	cnfg.PostgresHost.User = postgresUser

	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		logger.Fatalf("REDIS_PASSWORD environment variable is required")
	}

	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		logger.Fatalf("REDIS_HOST environment variable is required")
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		logger.Fatalf("REDIS_PORT environment variable is required")
	}

	cnfg.RedisHTTPServer.Password = redisPassword
	cnfg.RedisHTTPServer.Address = redisHost + ":" + redisPort

	fmt.Println(cnfg)

	redis, err := redis.NewClient(ctx, cnfg)
	if err != nil {
		logger.Fatalf("haven't connect with redis")
	}

	tokenCrypter := cryptoToken.NewTokenCrypter(key)

	telegrammClient := tgClient.New(cnfg.ClientsHosts.TelegramHost, token)

	tinkoffApiClient := tinkoffapi.NewClient(cnfg.ClientsHosts.TinkoffApiHost)

	userStorage, err := storage.NewStorage(ctx, cnfg, rootPath)
	if err != nil {
		logger.Fatalf("can't create user_storage: %s", err.Error())
	}
	defer func() {
		if userStorage != nil {
			userStorage.CloseDB()
		}
	}()

	bondReportServiceClient := bondreportservice.New(cnfg.ClientsHosts.BondReportServiceHost)

	processor := telegram.NewProccesor(
		tokenCrypter,
		telegrammClient,
		tinkoffApiClient,
		bondReportServiceClient,
		redis,
		userStorage,
	)

	fetcher := telegram.NewFetcher(telegrammClient)

	logger.Printf("service started")

	consumer := event_consumer.New(fetcher, processor, batchSize)

	if err := consumer.Start(); err != nil {
		logger.Fatalf("service is stopped")
	}
}
