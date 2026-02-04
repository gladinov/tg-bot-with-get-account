package service

import (
	"bonds-report-service/internal/clients/tinkoffApi"
	"bonds-report-service/internal/service/service_models"
	"log/slog"
	"time"
)

func (c *Client) TransOperations(operations []tinkoffApi.Operation) []service_models.Operation {
	const op = "service.TransformOperations"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
		)
	}()

	transformOperations := make([]service_models.Operation, 0)
	for _, v := range operations {
		transformOperation := service_models.Operation{
			Currency:          v.Currency,
			BrokerAccountId:   v.BrokerAccountId,
			Operation_Id:      v.Operation_Id,
			ParentOperationId: v.ParentOperationId,
			Name:              v.Name,
			Date:              v.Date,
			Type:              v.Type,
			Description:       v.Description,
			InstrumentUid:     v.InstrumentUid,
			Figi:              v.Figi,
			InstrumentType:    v.InstrumentType,
			InstrumentKind:    v.InstrumentKind,
			PositionUid:       v.PositionUid,
			Payment:           v.Payment.ToFloat(),
			Price:             v.Price.ToFloat(),
			Commission:        v.Commission.ToFloat(),
			Yield:             v.Yield.ToFloat(),
			YieldRelative:     v.YieldRelative.ToFloat(),
			AccruedInt:        v.AccruedInt.ToFloat(),
			QuantityDone:      float64(v.QuantityDone),
			AssetUid:          v.AssetUid,
		}

		transformOperations = append(transformOperations, transformOperation)
	}
	return transformOperations
}
