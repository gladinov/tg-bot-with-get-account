package instrumentsclient

import (
	"bonds-report-service/internal/models/domain"
	dtoTinkoff "bonds-report-service/internal/models/dto/tinkoffApi"
)

func MapFutureToDomain(dto dtoTinkoff.Future) domain.Future {
	return domain.Future{
		Name:                    dto.Name,
		MinPriceIncrement:       MapQuotationToDomain(dto.MinPriceIncrement),
		MinPriceIncrementAmount: MapQuotationToDomain(dto.MinPriceIncrementAmount),
		AssetType:               dto.AssetType,
		BasicAssetPositionUid:   dto.BasicAssetPositionUid,
	}
}

func MapQuotationToDomain(dtoQuat dtoTinkoff.Quotation) domain.Quotation {
	return domain.Quotation{
		Units: dtoQuat.Units,
		Nano:  dtoQuat.Nano,
	}
}

func MapMoneyValueToDomain(dto dtoTinkoff.MoneyValue) domain.MoneyValue {
	return domain.MoneyValue{
		Currency: dto.Currency,
		Units:    dto.Units,
		Nano:     dto.Nano,
	}
}

func MapBondToDomain(dto dtoTinkoff.Bond) domain.Bond {
	return domain.Bond{
		AciValue: MapMoneyValueToDomain(dto.AciValue),
		Currency: dto.Currency,
		Nominal:  MapMoneyValueToDomain(dto.Nominal),
	}
}

func MapCurrencyToDomain(dto dtoTinkoff.Currency) domain.Currency {
	return domain.Currency{
		Isin: dto.Isin,
	}
}

func MapInstrumentShortToDomain(dto dtoTinkoff.InstrumentShort) domain.InstrumentShort {
	return domain.InstrumentShort{
		InstrumentType: dto.InstrumentType,
		Uid:            dto.Uid,
		Figi:           dto.Figi,
	}
}

func MapSliceInstrumentShortToDomain(dto []dtoTinkoff.InstrumentShort) []domain.InstrumentShort {
	out := make([]domain.InstrumentShort, 0, len(dto))
	for _, v := range dto {
		dom := MapInstrumentShortToDomain(v)
		out = append(out, dom)
	}
	return out
}

func MapShareCurrencyByResponseToDomain(dto dtoTinkoff.ShareCurrencyByResponse) domain.ShareCurrency {
	return domain.ShareCurrency{
		Currency: dto.Currency,
	}
}
