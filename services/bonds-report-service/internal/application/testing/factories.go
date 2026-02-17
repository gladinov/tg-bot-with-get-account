package factories

import (
	"bonds-report-service/internal/domain"
	"time"
)

func NewValuesMoex() domain.ValuesMoex {
	return domain.ValuesMoex{
		ShortName:       domain.NewNullString("TestShortName", true, false),
		TradeDate:       domain.NewNullString("2026-02-12", true, false),
		MaturityDate:    domain.NewNullString("2030-02-12", true, false),
		OfferDate:       domain.NewNullString("2026-02-01", true, false),
		BuybackDate:     domain.NewNullString("2029-12-31", true, false),
		YieldToMaturity: domain.NewNullFloat64(7.5, true, false),
		YieldToOffer:    domain.NewNullFloat64(6.8, true, false),
		FaceValue:       domain.NewNullFloat64(1000, true, false),
		FaceUnit:        domain.NewNullString("RUB", true, false),
		Duration:        domain.NewNullFloat64(4.0, true, false),
	}
}

func NewCurrencyCBR(charCode string, value float64) domain.CurrencyCBR {
	return domain.CurrencyCBR{
		Date:      time.Date(2026, 2, 12, 0, 0, 0, 0, time.UTC),
		NumCode:   "840",
		CharCode:  charCode,
		Nominal:   1,
		Name:      "US Dollar",
		Value:     value,
		VunitRate: value,
	}
}

func NewCurrenciesCBR() domain.CurrenciesCBR {
	currencies := make(map[string]domain.CurrencyCBR)
	currencies["usd"] = NewCurrencyCBR("USD", 74.5)
	currencies["eur"] = NewCurrencyCBR("EUR", 80.0)
	return domain.CurrenciesCBR{
		CurrenciesMap: currencies,
	}
}

func NewOpenAccount() domain.Account {
	return domain.Account{
		ID:          "test-account-id",
		Type:        "broker",
		Name:        "Test Account",
		Status:      2,
		OpenedDate:  time.Date(2024, 2, 12, 0, 0, 0, 0, time.UTC),
		ClosedDate:  time.Date(2026, 2, 12, 0, 0, 0, 0, time.UTC),
		AccessLevel: 1,
	}
}

// NewPortfolioPosition создаёт фиктивную позицию для портфеля
func NewPortfolioPosition() domain.PortfolioPosition {
	return domain.PortfolioPosition{
		Figi:                     "figi_test_1",
		InstrumentType:           "share",
		Quantity:                 domain.NewQuotation(10, 0), // 10 шт.
		AveragePositionPrice:     domain.NewMoneyValue("RUB", 1000, 0),
		ExpectedYield:            domain.NewQuotation(50, 0),
		CurrentNkd:               domain.NewMoneyValue("RUB", 5, 0),
		CurrentPrice:             domain.NewMoneyValue("RUB", 1050, 0),
		AveragePositionPriceFifo: domain.NewMoneyValue("RUB", 1000, 0),
		Blocked:                  false,
		BlockedLots:              domain.NewQuotation(0, 0),
		PositionUid:              "pos_uid_1",
		InstrumentUid:            "instr_uid_1",
		VarMargin:                domain.NewMoneyValue("RUB", 0, 0),
		ExpectedYieldFifo:        domain.NewQuotation(50, 0),
		DailyYield:               domain.NewMoneyValue("RUB", 10, 0),
		Ticker:                   "TICKER1",
	}
}

// NewPortfolio создаёт фиктивный портфель с одной или несколькими позициями
func NewPortfolio(positions ...domain.PortfolioPosition) domain.Portfolio {
	if len(positions) == 0 {
		positions = []domain.PortfolioPosition{NewPortfolioPosition()}
	}

	total := domain.MoneyValue{
		Units:    0,
		Nano:     0,
		Currency: "RUB",
	}

	// Суммируем текущие цены позиций для total
	for _, p := range positions {
		total.Units += p.CurrentPrice.Units
		total.Nano += p.CurrentPrice.Nano
	}

	return domain.Portfolio{
		Positions:   positions,
		TotalAmount: total,
	}
}

func NewOperationsRequest() domain.OperationsRequest {
	return domain.OperationsRequest{
		AccountID: "account-123",
		FromDate:  time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
	}
}

func NewOperation() domain.Operation {
	return domain.Operation{
		BrokerAccountID:   "account-123",
		Currency:          "USD",
		OperationID:       "op-1",
		ParentOperationID: "parent-1",
		Name:              "Buy",
		Date:              time.Now(),
		Type:              1,
		Description:       "Test operation",
		InstrumentUid:     "instr-1",
		Figi:              "figi-1",
		InstrumentType:    "Stock",
		InstrumentKind:    "Equity",
		PositionUid:       "pos-1",
		Payment:           domain.NewMoneyValue("USD", 100, 0),
		Price:             domain.NewMoneyValue("USD", 50, 0),
		Commission:        domain.NewMoneyValue("USD", 1, 0),
		Yield:             domain.NewMoneyValue("USD", 2, 0),
		YieldRelative:     domain.NewQuotation(0, 0),
		AccruedInt:        domain.NewMoneyValue("USD", 0, 0),
		QuantityDone:      2,
		AssetUid:          "asset-1",
	}
}

func NewBondIdentIdentifiers() domain.BondIdentIdentifiers {
	return domain.BondIdentIdentifiers{
		Ticker:          "TSTBOND",
		ClassCode:       "TST",
		Name:            "Test Bond",
		Nominal:         domain.NewMoneyValue("USD", 1000, 0),
		NominalCurrency: "USD",
		Replaced:        false,
	}
}

func NewFuture() domain.Future {
	return domain.Future{
		Name:                    "Test Future",
		MinPriceIncrement:       domain.NewQuotation(1, 0),
		MinPriceIncrementAmount: domain.NewQuotation(100, 0),
		AssetType:               "Future",
		BasicAssetPositionUid:   "pos-123",
	}
}

func NewBond() domain.Bond {
	return domain.Bond{
		AciValue: domain.NewMoneyValue("RUB", 15, 50_000_000), // 15.05
		Currency: "RUB",
		Nominal:  domain.NewMoneyValue("RUB", 1000, 0),
	}
}

func NewCurrency() domain.Currency {
	return domain.Currency{
		Isin: "US1234567890",
	}
}

func NewShareCurrency() domain.ShareCurrency {
	return domain.ShareCurrency{
		Currency: "USD",
	}
}

func NewInstrumentShort() domain.InstrumentShort {
	return domain.InstrumentShort{
		InstrumentType: "share",
		Uid:            "uid-123",
		Figi:           "figi-123",
	}
}

func NewInstrumentShortList(instruments ...domain.InstrumentShort) domain.InstrumentShortList {
	return instruments
}

func NewLastPrice() domain.LastPrice {
	return domain.LastPrice{
		LastPrice: domain.Quotation{
			Units: 100,
			Nano:  0,
		},
	}
}
