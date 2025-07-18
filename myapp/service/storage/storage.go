package service_storage

import (
	"context"

	"main.go/service/service_models"
)

type Storage interface {
	OperationStorage
	BondReportStorage
	CurrencyStorage
}

type OperationStorage interface {
	SaveOperations(ctx context.Context, chatID int, accountId string, operations []service_models.Operation) error
	GetOperations(ctx context.Context, chatId int, assetUid string, accountId string) ([]service_models.Operation, error)
}

type BondReportStorage interface {
	SaveBondReport(ctx context.Context, chatID int, accountId string, bondReport []service_models.BondReport) error
}

type CurrencyStorage interface {
}
