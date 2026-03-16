package app

import (
	"bonds-report-service/internal/application/ports"
	config "bonds-report-service/internal/configs"
	"bonds-report-service/internal/infrastructure/repository/postgreSQL"
	"context"
	"fmt"
	"log/slog"
)

func MustInitNewStorage(ctx context.Context, config config.Config, logg *slog.Logger) ports.Storage {
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
