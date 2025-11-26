package repository

import (
	"context"
	"errors"
	"time"

	config "bonds-report-service/internal/configs"
	"bonds-report-service/internal/repository/postgreSQL"
	servicet_sqlite "bonds-report-service/internal/repository/sqlite"
	"bonds-report-service/internal/service/service_models"
	pathwd "bonds-report-service/lib/pathWD"
)

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
	SaveOperations(ctx context.Context, chatID int, accountId string, operations []service_models.Operation) error
	GetOperations(ctx context.Context, chatId int, assetUid string, accountId string) ([]service_models.Operation, error)
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
	SaveCurrency(ctx context.Context, currencies service_models.Currencies, date time.Time) error
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

func MustInitNewStorage(ctx context.Context, config config.Config, rootPath string) Storage {
	switch config.DbType {
	case postreSQL:
		serviceStorage, err := postgreSQL.NewStorage(config)
		if err != nil {
			panic(err)
		}
		err = serviceStorage.InitDB(ctx)
		if err != nil {
			panic(err)
		}

		return serviceStorage

	case SQLite:
		serviceStorageAbsolutPath, err := pathwd.PathFromWD(rootPath, config.ServiceStorageSQLLitePath)
		if err != nil {
			panic(err)
		}

		serviceStorage, err := servicet_sqlite.New(serviceStorageAbsolutPath)
		if err != nil {
			panic(err)
		}

		if err := serviceStorage.Init(ctx); err != nil {
			panic(err)
		}
		return serviceStorage
	default:
		panic(errors.New("possible init only SQLite or PostgreSQL databases"))
	}
}
