package main

import (
	"context"
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

	cnfg := config.MustInitConfig()

	//// //  for local
	//envPath := filepath.Join(rootPath, ".env")
	//err := godotenv.Load(envPath)
	//if err != nil {
	//	logger.Printf("Error loading .env file. Erorr: %v", err.Error())
	//}
	//

	redis, err := redis.NewClient(ctx, cnfg)
	if err != nil {
		logger.Fatalf("haven't connect with redis")
	}

	tokenCrypter := cryptoToken.NewTokenCrypter(cnfg.Key)

	telegrammClient := tgClient.New(cnfg.ClientsHosts.TelegramHost, cnfg.Token)

	tinkoffApiClient := tinkoffapi.NewClient(cnfg.ClientsHosts.GetTinkoffApiAddress())

	userStorage, err := storage.NewStorage(ctx, cnfg)
	if err != nil {
		logger.Fatalf("can't create user_storage: %s", err.Error())
	}
	defer func() {
		if userStorage != nil {
			userStorage.CloseDB()
		}
	}()

	bondReportServiceClient := bondreportservice.New(cnfg.ClientsHosts.GetBondReportAddress())

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
