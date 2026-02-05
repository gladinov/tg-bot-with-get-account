package moex

import (
	"bonds-report-service/internal/models/domain"
	"bonds-report-service/internal/models/dto/moex"
)

func MapValueFromDTOToDomain(dtoValues moex.Values) domain.ValuesMoex {
	return domain.ValuesMoex{
		ShortName:       MapNullStringFromDTOToDomain(dtoValues.ShortName),
		TradeDate:       MapNullStringFromDTOToDomain(dtoValues.TradeDate),
		MaturityDate:    MapNullStringFromDTOToDomain(dtoValues.MaturityDate),
		OfferDate:       MapNullStringFromDTOToDomain(dtoValues.OfferDate),
		BuybackDate:     MapNullStringFromDTOToDomain(dtoValues.BuybackDate),
		YieldToMaturity: MapNullFloat64FromDTOToDomain(dtoValues.YieldToMaturity),
		YieldToOffer:    MapNullFloat64FromDTOToDomain(dtoValues.YieldToOffer),
		FaceValue:       MapNullFloat64FromDTOToDomain(dtoValues.FaceValue),
		FaceUnit:        MapNullStringFromDTOToDomain(dtoValues.FaceUnit),
		Duration:        MapNullFloat64FromDTOToDomain(dtoValues.Duration),
	}
}

func MapNullStringFromDTOToDomain(dtoNullString moex.NullString) domain.NullString {
	return domain.NewNullString(dtoNullString.Value, dtoNullString.IsSet, dtoNullString.IsNull)
}

func MapNullFloat64FromDTOToDomain(dtoNullFLoat64 moex.NullFloat64) domain.NullFloat64 {
	return domain.NewNullFloat64(dtoNullFLoat64.Value, dtoNullFLoat64.IsSet, dtoNullFLoat64.IsNull)
}
