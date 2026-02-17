package app

import (
	service "bonds-report-service/internal/application"
	config "bonds-report-service/internal/configs"
	"bonds-report-service/internal/infrastructure/repository/postgreSQL"
	"context"
	"fmt"
	"log/slog"
)

func MustInitNewStorage(ctx context.Context, config config.Config, logg *slog.Logger) service.Storage {
	const op = "repository.MustInitNewStorage"
	logg.Debug(fmt.Sprintf("start %s", op))

	serviceStorage, err := postgreSQL.NewStorage(logg, config)
	if err != nil {
		logg.Debug("failed to create PostgreSQL storage", "err", err)
		panic(err)
	}
	err = serviceStorage.InitDB(ctx)
	if err != nil {
		logg.Debug("failed to init PostgreSQL database", "err", err)
		panic(err)
	}
	logg.Info("PostgreSQL storage initialized successfully")
	return serviceStorage
}
