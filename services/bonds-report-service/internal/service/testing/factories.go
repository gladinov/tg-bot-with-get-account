package factories

import (
	"bonds-report-service/internal/models/domain"
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
