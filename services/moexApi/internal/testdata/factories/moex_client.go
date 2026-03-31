package factories

import (
	"encoding/json"
	"moex/internal/models"
)

func NewValuesArray() []any {
	return []any{
		"2025-01-15", // TradeDate
		"2035-01-15", // MaturityDate
		"",           // OfferDate
		"",           // BuybackDate
		8.42,         // YieldToMaturity
		8.30,         // YieldToOffer
		1000.0,       // FaceValue
		"RUB",        // FaceUnit
		365.5,        // Duration
		"OFZ 26238",  // ShortName
	}
}

func NewSpecificationsResponseJSON() []byte {
	resp := map[string]any{
		"history": map[string]any{
			"data": []any{
				NewValuesArray(),
			},
		},
	}

	b, err := json.Marshal(resp)
	if err != nil {
		panic(err) // в тестах удобно паниковать, чтобы сразу заметить ошибку
	}
	return b
}

func NS(val string) models.NullString {
	return models.NullString{
		Value:  val,
		IsSet:  true,
		IsNull: false,
	}
}

func NF(val float64) models.NullFloat64 {
	return models.NullFloat64{
		Value:  val,
		IsSet:  true,
		IsNull: false,
	}
}

func NewValues() models.Values {
	return models.Values{
		ShortName:       NS("OFZ 26238"),
		TradeDate:       NS("2025-01-15"),
		MaturityDate:    NS("2035-01-15"),
		OfferDate:       NS(""),
		BuybackDate:     NS(""),
		YieldToMaturity: NF(8.42),
		YieldToOffer:    NF(8.30),
		FaceValue:       NF(1000),
		FaceUnit:        NS("RUB"),
		Duration:        NF(365.5),
	}
}

func NewSpecificationsResponse() models.SpecificationsResponce {
	return models.SpecificationsResponce{
		History: &models.History{
			Data: []models.Values{
				NewValues(),
			},
		},
	}
}
