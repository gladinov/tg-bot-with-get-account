package analyticsclient

import (
	"bonds-report-service/internal/models/domain"
	dtoTinkoff "bonds-report-service/internal/models/dto/tinkoffApi"
)

func MapBondIdentIdentifiers(dtoBondId dtoTinkoff.BondIdentIdentifiers) domain.BondIdentIdentifiers {
	return domain.BondIdentIdentifiers{
		Ticker:          dtoBondId.Ticker,
		ClassCode:       dtoBondId.ClassCode,
		Name:            dtoBondId.Name,
		Nominal:         MapMoneyValue(dtoBondId.Nominal),
		NominalCurrency: dtoBondId.NominalCurrency,
		Replaced:        dtoBondId.Replaced,
	}
}

func MapMoneyValue(dtoMoneyValue dtoTinkoff.MoneyValue) domain.MoneyValue {
	return domain.MoneyValue{
		Currency: dtoMoneyValue.Currency,
		Units:    dtoMoneyValue.Units,
		Nano:     dtoMoneyValue.Nano,
	}
}

func MapLastPriceResponseToDomain(dtoResp dtoTinkoff.LastPriceResponse) domain.LastPrice {
	return domain.LastPrice{
		LastPrice: MapQuotationToDomain(dtoResp.LastPrice),
	}
}

func MapQuotationToDomain(dtoQuat dtoTinkoff.Quotation) domain.Quotation {
	return domain.Quotation{
		Units: dtoQuat.Units,
		Nano:  dtoQuat.Nano,
	}
}
