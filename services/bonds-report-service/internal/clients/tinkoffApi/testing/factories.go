package factories

import (
	"bonds-report-service/internal/models/domain"
	"math"
	"time"

	tinkoffDto "bonds-report-service/internal/models/dto/tinkoffApi"
)

func NewTestLastPrice(numb float64) domain.LastPrice {
	return domain.LastPrice{
		LastPrice: QuotationFromFloat(numb),
	}
}

func QuotationFromFloat(v float64) domain.Quotation {
	const nanoFactor = 1_000_000_000

	v = math.Round(v*nanoFactor) / nanoFactor

	units := int64(v)
	nano := int32(math.Round((v - float64(units)) * nanoFactor))

	if nano == nanoFactor {
		units++
		nano = 0
	}

	return domain.Quotation{
		Units: units,
		Nano:  nano,
	}
}

func NewTestLastPriceDto(numb float64) tinkoffDto.LastPriceResponse {
	return tinkoffDto.LastPriceResponse{
		LastPrice: QuotationDtoFromFloat(numb),
	}
}

func QuotationDtoFromFloat(v float64) tinkoffDto.Quotation {
	const nanoFactor = 1_000_000_000

	v = math.Round(v*nanoFactor) / nanoFactor

	units := int64(v)
	nano := int32(math.Round((v - float64(units)) * nanoFactor))

	if nano == nanoFactor {
		units++
		nano = 0
	}

	return tinkoffDto.Quotation{
		Units: units,
		Nano:  nano,
	}
}

func NewFutureDTO() tinkoffDto.Future {
	return tinkoffDto.Future{
		Name: "Test Future",
		MinPriceIncrement: tinkoffDto.Quotation{
			Units: 1,
			Nano:  0,
		},
		MinPriceIncrementAmount: tinkoffDto.Quotation{
			Units: 10,
			Nano:  0,
		},
		AssetType:             "commodity",
		BasicAssetPositionUid: "basic-uid-123",
	}
}

func NewBond() tinkoffDto.Bond {
	return tinkoffDto.Bond{
		AciValue: tinkoffDto.MoneyValue{
			Currency: "RUB",
			Units:    10,
			Nano:     500_000_000,
		},
		Currency: "RUB",
		Nominal: tinkoffDto.MoneyValue{
			Currency: "RUB",
			Units:    1000,
			Nano:     0,
		},
	}
}

func NewInstrumentShortTest() tinkoffDto.InstrumentShort {
	return tinkoffDto.InstrumentShort{
		InstrumentType: "bond",
		Uid:            "UID123",
		Figi:           "FIGI123",
	}
}

func NewPortfolioPosition() tinkoffDto.PortfolioPositions {
	return tinkoffDto.PortfolioPositions{
		Figi:           "TEST_FIGI",
		InstrumentType: "share",
		Quantity: tinkoffDto.Quotation{
			Units: 10,
			Nano:  0,
		},
		AveragePositionPrice: tinkoffDto.MoneyValue{
			Currency: "RUB",
			Units:    100,
			Nano:     0,
		},
		ExpectedYield: tinkoffDto.Quotation{
			Units: 5,
			Nano:  0,
		},
		CurrentNkd: tinkoffDto.MoneyValue{
			Currency: "RUB",
			Units:    0,
			Nano:     0,
		},
		CurrentPrice: tinkoffDto.MoneyValue{
			Currency: "RUB",
			Units:    110,
			Nano:     0,
		},
		AveragePositionPriceFifo: tinkoffDto.MoneyValue{
			Currency: "RUB",
			Units:    100,
			Nano:     0,
		},
		Blocked: false,
		BlockedLots: tinkoffDto.Quotation{
			Units: 0,
			Nano:  0,
		},
		PositionUid:   "POSITION_UID",
		InstrumentUid: "INSTRUMENT_UID",
		VarMargin: tinkoffDto.MoneyValue{
			Currency: "RUB",
			Units:    0,
			Nano:     0,
		},
		ExpectedYieldFifo: tinkoffDto.Quotation{
			Units: 5,
			Nano:  0,
		},
		DailyYield: tinkoffDto.MoneyValue{
			Currency: "RUB",
			Units:    1,
			Nano:     0,
		},
		Ticker: "TEST",
	}
}

func NewPortfolio() tinkoffDto.Portfolio {
	position := NewPortfolioPosition()

	return tinkoffDto.Portfolio{
		Positions: []tinkoffDto.PortfolioPositions{position},
		TotalAmount: tinkoffDto.MoneyValue{
			Currency: "RUB",
			Units:    1000,
			Nano:     0,
		},
	}
}

func NewOperationDTO() tinkoffDto.Operation {
	return tinkoffDto.Operation{
		BrokerAccountID:   "ACC123",
		Currency:          "USD",
		OperationID:       "OP123",
		ParentOperationID: "POP123",
		Name:              "Buy Stock",
		Date:              time.Now(),
		Type:              1,
		Description:       "Test operation",
		InstrumentUid:     "INST123",
		Figi:              "FIGI123",
		InstrumentType:    "Stock",
		InstrumentKind:    "Common",
		PositionUid:       "POS123",
		Payment:           NewMoneyValueDTO(),
		Price:             NewMoneyValueDTO(),
		Commission:        NewMoneyValueDTO(),
		Yield:             NewMoneyValueDTO(),
		YieldRelative:     NewQuotationDTO(),
		AccruedInt:        NewMoneyValueDTO(),
		QuantityDone:      100,
		AssetUid:          "ASSET123",
	}
}

func NewMoneyValueDTO() tinkoffDto.MoneyValue {
	return tinkoffDto.MoneyValue{
		Currency: "USD",
		Units:    100,
		Nano:     500_000_000,
	}
}

func NewQuotationDTO() tinkoffDto.Quotation {
	return tinkoffDto.Quotation{
		Units: 1,
		Nano:  250_000_000,
	}
}
