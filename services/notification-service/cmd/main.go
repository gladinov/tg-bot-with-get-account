package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	sl "github.com/gladinov/mylogger"
	"github.com/gladinov/notification-service/internal/adapters/inbound/kafka"
	"github.com/gladinov/notification-service/internal/adapters/outbound/telegram"
	tgClient "github.com/gladinov/notification-service/internal/adapters/outbound/telegram"
	"github.com/gladinov/notification-service/internal/application/usecases"
	"github.com/gladinov/notification-service/internal/config"
	"github.com/twmb/franz-go/pkg/kgo"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	conf := config.MustInitConfig()

	logg := sl.NewLogger(conf.Env)

	logg.Info("start app",
		slog.String("env", conf.Env),
		slog.String("notification-service_app_host", conf.Clients.NotificationService.Host),
		slog.String("notification-service_app_port", conf.Clients.NotificationService.Port))

	// TODO: logg init kafkaClient
	logg.InfoContext(ctx, "initialize kafka client", slog.Any("host", conf.Kafka.Host), slog.Any("port", conf.Kafka.Port))
	kafkaClient, err := kgo.NewClient(
		kgo.SeedBrokers(conf.Kafka.GetKafkaAddress()),
	)
	if err != nil {
		logg.Error("haven't connect with kafka", slog.String("err", err.Error()))
		return
	}

	if err := kafkaClient.Ping(ctx); err != nil {
		logg.ErrorContext(ctx, "kafka not available", slog.Any("error", err))
		return
	}

	logg.Info("initialize Telegram client", slog.String("addres", conf.Clients.Telegram.Host))
	telegramClient := tgClient.New(logg, conf.Clients.Telegram.Host, conf.Clients.Telegram.Token)

	service := usecases.NewService(logg, telegramClient)

	handler := kafka.NewHandler(logg, service)

	consumer := kafka.NewConsumer(logg, kafkaClient, handler)

	errCh := make(chan error, 1)
	go func() {
		err := consumer.Run(ctx)
		errCh <- err
	}()

	select {
	case <-ctx.Done():
		logg.InfoContext(ctx, "Shutdown signal received")
	case er := <-errCh:
		logg.ErrorContext(ctx, "consumer error", slog.Any("error", er))
	}

	gracefulShutdown(ctx, logg, kafkaClient, telegramClient)
}

func gracefulShutdown(ctx context.Context, logg *slog.Logger, kafkaClient *kgo.Client, telegramClient *telegram.Client) {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logg.InfoContext(ctx, "close kafka")
	kafkaClient.LeaveGroupContext(shutdownCtx)
	kafkaClient.Close()
	logg.InfoContext(ctx, "close tg client")
	telegramClient.Close()

	<-shutdownCtx.Done()

	logg.InfoContext(ctx, "Server exited gracefully")
}
