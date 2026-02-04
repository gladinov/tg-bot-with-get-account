package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gladinov/cryptotoken"
	"github.com/gladinov/mylogger"
	"github.com/gladinov/traceidgenerator"
	bondreportservice "main.go/clients/bondReportService"
	tgClient "main.go/clients/telegram"
	tinkoffapi "main.go/clients/tinkoffApi"
	event_consumer "main.go/internal/app/consumer/event-consumer"
	"main.go/internal/app/events/telegram"
	"main.go/internal/config"
	storage "main.go/internal/repository"
	"main.go/internal/repository/redis"
	tokenauth "main.go/internal/tokenAuth"
)

const (
	batchSize = 100
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	conf := config.MustInitConfig()

	_ = traceidgenerator.Must()

	logg := mylogger.NewLogger(conf.Env)

	logg.Info("start app",
		slog.String("env", conf.Env),
	)

	logg.Info("initialize Redis client", slog.String("addres", conf.RedisHTTPServer.GetAddress()))
	redis, err := redis.NewClient(ctx, conf)
	if err != nil {
		logg.Error("haven't connect with redis", slog.String("err", err.Error()))
		return
	}

	logg.Info("initialize TokenCrypter client")
	tokenCrypter := cryptotoken.NewTokenCrypter(conf.Key)

	logg.Info("initialize Telegram client", slog.String("addres", conf.ClientsHosts.TelegramHost))
	telegrammClient := tgClient.New(logg, conf.ClientsHosts.TelegramHost, conf.Token)

	logg.Info("initialize Tinkoff client", slog.String("addres", conf.ClientsHosts.GetTinkoffApiAddress()))
	tinkoffApiClient := tinkoffapi.NewClient(logg, conf.ClientsHosts.GetTinkoffApiAddress())

	logg.Info("initialize User storage",
		slog.String("dbType", conf.DbType),
		slog.String("address", conf.PostgresHost.GetAdress()),
	)
	userStorage, err := storage.NewStorage(ctx, conf)
	if err != nil {
		logg.Error("can't create user_storage", slog.String("err", err.Error()))
		return
	}
	defer func() {
		if userStorage != nil {
			userStorage.CloseDB()
		}
	}()

	logg.Info("initialize bondReportService client", slog.String("addres", conf.ClientsHosts.GetBondReportAddress()))
	bondReportServiceClient := bondreportservice.New(logg, conf.ClientsHosts.GetBondReportAddress())

	logg.Info("initialize TokenAuthService")
	tokenAuthService := tokenauth.NewTokenAuthService(
		logg,
		redis, // TODO: Переместить redis cashe из слоя service в слой repo
		userStorage,
		tinkoffApiClient,
		tokenCrypter)

	logg.Info("initialize Processor")
	processor := telegram.NewProccesor(
		logg,
		telegrammClient,
		tinkoffApiClient,
		bondReportServiceClient,
		tokenAuthService,
	)

	logg.Info("initialize Fetcher")
	fetcher := telegram.NewFetcher(logg, telegrammClient)

	logg.Info("service started")
	consumer := event_consumer.New(logg, fetcher, processor, batchSize)

	if err := consumer.Start(); err != nil {
		logg.Error("service is stopped")
		return
	}
}
