package repository

import (
	"bonds-report-service/internal/models/domain"
	"bonds-report-service/internal/repository/postgreSQL"
	"bonds-report-service/internal/service/service_models"
	"context"
	"fmt"
	"log/slog"
	"time"

	config "bonds-report-service/internal/configs"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=Storage
type Storage interface {
	OperationStorage
	BondReportStorage
	GeneralBondReportStorage
	CurrencyStorage
	UidsStorage
	CloseStorage
}

type OperationStorage interface {
	LastOperationTime(ctx context.Context, chatID int, accountId string) (time.Time, error)
	SaveOperations(ctx context.Context, chatID int, accountId string, operations []domain.OperationWithoutCustomTypes) error
	GetOperations(ctx context.Context, chatId int, assetUid string, accountId string) ([]domain.OperationWithoutCustomTypes, error)
}

type BondReportStorage interface {
	DeleteBondReport(ctx context.Context, chatID int, accountId string) (err error)
	SaveBondReport(ctx context.Context, chatID int, accountId string, bondReport []service_models.BondReport) error
}

type GeneralBondReportStorage interface {
	DeleteGeneralBondReport(ctx context.Context, chatID int, accountId string) (err error)
	SaveGeneralBondReport(ctx context.Context, chatID int, accountId string, bondReport []service_models.GeneralBondReportPosition) error
}

type CurrencyStorage interface {
	SaveCurrency(ctx context.Context, currencies domain.CurrenciesCBR, date time.Time) error
	GetCurrency(ctx context.Context, currency string, date time.Time) (float64, error)
}

type UidsStorage interface {
	SaveUids(ctx context.Context, uids map[string]string) error
	IsUpdatedUids(ctx context.Context) (time.Time, error)
	GetUid(ctx context.Context, instrumentUid string) (string, error)
}

type CloseStorage interface {
	CloseDB()
}

const (
	postreSQL = "postgreSQL"
	SQLite    = "SQLite"
)

func MustInitNewStorage(ctx context.Context, config config.Config, logg *slog.Logger) Storage {
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
