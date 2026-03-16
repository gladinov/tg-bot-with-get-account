package updateoperations

import (
	tinkoffHelper "bonds-report-service/internal/application/helpers/tinkoff"
	"bonds-report-service/internal/application/ports"
	"bonds-report-service/internal/domain"
	"bonds-report-service/internal/domain/mapper"
	"bonds-report-service/internal/utils/logging"
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/gladinov/e"
)

type Updater struct {
	logger        *slog.Logger
	Storage       ports.Storage
	TinkoffHelper *tinkoffHelper.TinkoffHelper
}

func NewUpdater(logger *slog.Logger, storage ports.Storage, tinkoffHelper *tinkoffHelper.TinkoffHelper) *Updater {
	return &Updater{
		logger:        logger,
		Storage:       storage,
		TinkoffHelper: tinkoffHelper,
	}
}

func (h *Updater) UpdateOperations(ctx context.Context,
	chatID int,
	accountID string,
	openDate time.Time,
) (err error) {
	const op = "updateoperations.updateOperations"

	defer logging.LogOperation_Debug(ctx, h.logger, op, &err)()

	fromDate, err := h.Storage.LastOperationTime(ctx, chatID, accountID)
	// TODO: Если fromDate будет больше time.Now, то будет ошибка.
	// Есть вероятность такой ошибки, поэтому при тестировании функции нужно придумать другой способ вызова функции по последней операции
	fromDate = fromDate.Add(time.Microsecond * 1)

	if err != nil {
		if errors.Is(err, domain.ErrNoOpperations) {
			fromDate = openDate
		} else {
			return e.WrapIfErr("can't get last op from storage", err)
		}
	}

	opRequest := domain.NewOperationsRequest(accountID, fromDate)

	tinkoffOperations, err := h.TinkoffHelper.TinkoffGetOperations(ctx, opRequest)
	if err != nil {
		return e.WrapIfErr("can't get operations from tinkoff", err)
	}

	// НЕ domain mapper
	operations := mapper.MapOperationToOperationWithoutCustomTypes(tinkoffOperations)

	err = h.Storage.SaveOperations(ctx, chatID, accountID, operations)
	if err != nil {
		return e.WrapIfErr("can't save ops to Storage", err)
	}
	return nil
}
