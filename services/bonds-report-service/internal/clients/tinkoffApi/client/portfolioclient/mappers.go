package portfolioclient

import (
	"bonds-report-service/internal/models/domain"
	dtoTinkoff "bonds-report-service/internal/models/dto/tinkoffApi"
)

func MapAccountToDomain(a dtoTinkoff.Account) domain.Account {
	return domain.Account{
		ID:          a.ID,
		Type:        a.Type,
		Name:        a.Name,
		Status:      a.Status,
		OpenedDate:  a.OpenedDate,
		ClosedDate:  a.ClosedDate,
		AccessLevel: a.AccessLevel,
	}
}

func MapAccountsToDomain(accounts map[string]dtoTinkoff.Account) map[string]domain.Account {
	out := make(map[string]domain.Account, len(accounts))

	for key, value := range accounts {
		out[key] = MapAccountToDomain(value)
	}

	return out
}

func MapPortfolioPositionToDomain(dto dtoTinkoff.PortfolioPositions) domain.PortfolioPosition {
	return domain.PortfolioPosition{
		Figi:                     dto.Figi,
		InstrumentType:           dto.InstrumentType,
		Quantity:                 MapQuotationToDomain(dto.Quantity),
		AveragePositionPrice:     MapMoneyValueToDomain(dto.AveragePositionPrice),
		ExpectedYield:            MapQuotationToDomain(dto.ExpectedYield),
		CurrentNkd:               MapMoneyValueToDomain(dto.CurrentNkd),
		CurrentPrice:             MapMoneyValueToDomain(dto.CurrentPrice),
		AveragePositionPriceFifo: MapMoneyValueToDomain(dto.AveragePositionPriceFifo),
		Blocked:                  dto.Blocked,
		BlockedLots:              MapQuotationToDomain(dto.BlockedLots),
		PositionUid:              dto.PositionUid,
		InstrumentUid:            dto.InstrumentUid,
		VarMargin:                MapMoneyValueToDomain(dto.VarMargin),
		ExpectedYieldFifo:        MapQuotationToDomain(dto.ExpectedYieldFifo),
		DailyYield:               MapMoneyValueToDomain(dto.DailyYield),
		Ticker:                   dto.Ticker,
	}
}

func MapPortfolioToDomain(dto dtoTinkoff.Portfolio) domain.Portfolio {
	positions := make([]domain.PortfolioPosition, len(dto.Positions))
	for i, p := range dto.Positions {
		positions[i] = MapPortfolioPositionToDomain(p)
	}
	return domain.Portfolio{
		Positions:   positions,
		TotalAmount: MapMoneyValueToDomain(dto.TotalAmount),
	}
}

func MapMoneyValueToDomain(dto dtoTinkoff.MoneyValue) domain.MoneyValue {
	return domain.MoneyValue{
		Currency: dto.Currency,
		Units:    dto.Units,
		Nano:     dto.Nano,
	}
}

func MapQuotationToDomain(dtoQuat dtoTinkoff.Quotation) domain.Quotation {
	return domain.Quotation{
		Units: dtoQuat.Units,
		Nano:  dtoQuat.Nano,
	}
}

func MapOperationToDomain(op dtoTinkoff.Operation) domain.Operation {
	return domain.Operation{
		BrokerAccountID:   op.BrokerAccountID,
		Currency:          op.Currency,
		OperationID:       op.OperationID,
		ParentOperationID: op.ParentOperationID,
		Name:              op.Name,
		Date:              op.Date,
		Type:              op.Type,
		Description:       op.Description,
		InstrumentUid:     op.InstrumentUid,
		Figi:              op.Figi,
		InstrumentType:    op.InstrumentType,
		InstrumentKind:    op.InstrumentKind,
		PositionUid:       op.PositionUid,
		Payment:           MapMoneyValueToDomain(op.Payment),
		Price:             MapMoneyValueToDomain(op.Price),
		Commission:        MapMoneyValueToDomain(op.Commission),
		Yield:             MapMoneyValueToDomain(op.Yield),
		YieldRelative:     MapQuotationToDomain(op.YieldRelative),
		AccruedInt:        MapMoneyValueToDomain(op.AccruedInt),
		QuantityDone:      op.QuantityDone,
		AssetUid:          op.AssetUid,
	}
}

func MapOperationsToDomain(ops []dtoTinkoff.Operation) []domain.Operation {
	out := make([]domain.Operation, 0, len(ops))

	for _, op := range ops {
		out = append(out, MapOperationToDomain(op))
	}

	return out
}
