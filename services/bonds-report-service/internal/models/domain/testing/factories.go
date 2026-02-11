package factories

import (
	"bonds-report-service/internal/models/domain"
	"time"
)

var times time.Time = time.Time{}

func NewOperation() domain.Operation {
	return domain.Operation{
		Currency:          "USD",
		BrokerAccountID:   "acc1",
		OperationID:       "op1",
		ParentOperationID: "parent1",
		Name:              "Buy",
		Date:              times,
		Type:              1,
		Description:       "Test operation",
		InstrumentUid:     "instr1",
		Figi:              "figi1",
		InstrumentType:    "stock",
		InstrumentKind:    "kind1",
		PositionUid:       "pos1",
		Payment:           domain.MoneyValue{Units: 100, Nano: 500_000_000},
		Price:             domain.MoneyValue{Units: 101, Nano: 0},
		Commission:        domain.MoneyValue{Units: 1, Nano: 0},
		Yield:             domain.MoneyValue{Units: 2, Nano: 250_000_000},
		YieldRelative:     domain.Quotation{Units: 1, Nano: 500_000_000},
		AccruedInt:        domain.MoneyValue{Units: 3, Nano: 0},
		QuantityDone:      10,
		AssetUid:          "asset1",
	}
}

// OperationWithoutCustomTypes
func NewOperationWithoutCustomTypes() domain.OperationWithoutCustomTypes {
	return domain.OperationWithoutCustomTypes{
		Currency:          "USD",
		BrokerAccountID:   "acc1",
		OperationID:       "op1",
		ParentOperationID: "parent1",
		Name:              "Buy",
		Date:              times,
		Type:              1,
		Description:       "Test operation",
		InstrumentUid:     "instr1",
		Figi:              "figi1",
		InstrumentType:    "stock",
		InstrumentKind:    "kind1",
		PositionUid:       "pos1",
		Payment:           100.5,
		Price:             101,
		Commission:        1,
		Yield:             2.25,
		YieldRelative:     1.5,
		AccruedInt:        3,
		QuantityDone:      10,
		AssetUid:          "asset1",
	}
}
