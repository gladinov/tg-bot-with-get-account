package mapper

import "bonds-report-service/internal/domain"

func MapOperationToOperationWithoutCustomTypes(operations []domain.Operation) []domain.OperationWithoutCustomTypes {
	const op = "service.TransformOperations"

	transformOperations := make([]domain.OperationWithoutCustomTypes, 0, len(operations))
	for _, v := range operations {
		transformOperation := domain.OperationWithoutCustomTypes{
			Currency:          v.Currency,
			BrokerAccountID:   v.BrokerAccountID,
			OperationID:       v.OperationID,
			ParentOperationID: v.ParentOperationID,
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
