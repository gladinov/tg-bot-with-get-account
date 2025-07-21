package service_storage

import (
	"context"
	"time"

	"main.go/service/service_models"
)

type Storage interface {
	OperationStorage
	BondReportStorage
	CurrencyStorage
	UidsStorage
}

type OperationStorage interface {
	SaveOperations(ctx context.Context, chatID int, accountId string, operations []service_models.Operation) error
	GetOperations(ctx context.Context, chatId int, assetUid string, accountId string) ([]service_models.Operation, error)
}

type BondReportStorage interface {
	SaveBondReport(ctx context.Context, chatID int, accountId string, bondReport []service_models.BondReport) error
}

type CurrencyStorage interface {
	SaveCurrency(ctx context.Context, currencies service_models.Currencies) error
	GetCurrency(ctx context.Context, currency string, date time.Time) (float64, error)
}

type UidsStorage interface {
	SaveUids(ctx context.Context, uids map[string]string) error
	IsUpdatedUids(ctx context.Context) (bool, error)
	GetUid(ctx context.Context, instrumentUid string) (string, error)
}
