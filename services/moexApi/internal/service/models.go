package service

import (
	"encoding/json"
	"fmt"
	"time"
)

type SpecificationsRequest struct {
	Ticker string    `json:"ticker"`
	Date   time.Time `json:"date"`
}

type SpecificationsResponce struct {
	History *History `json:"history"`
}

type History struct {
	Data []Values `json:"data"`
}
type Nullable interface {
	IsSet() bool
	IsNull() bool
}

type NullString struct {
	Value  string `json:"value"`
	IsSet  bool   `json:"isSet"`
	IsNull bool   `json:"isNull"`
}

func (ns NullString) GetValue() string {
	return ns.Value
}

func (ns NullString) GetIsSet() bool {
	return ns.IsSet
}

func (ns NullString) GetIsNull() bool {
	return ns.IsNull
}

type NullFloat64 struct {
	Value  float64 `json:"value"`
	IsSet  bool    `json:"isSet"`
	IsNull bool    `json:"isNull"`
}

func (nf NullFloat64) GetValue() float64 {
	return nf.Value
}

func (nf NullFloat64) GetIsSet() bool {
	return nf.IsSet
}

func (nf NullFloat64) GetIsNull() bool {
	return nf.IsNull
}

type Values struct {
	ShortName       NullString  `json:"SHORTNAME"`
	TradeDate       NullString  `json:"TRADEDATE"`    // Торговая дата(на момент которой рассчитаны остальные данные)
	MaturityDate    NullString  `json:"MATDATE"`      // Дата погашения
	OfferDate       NullString  `json:"OFFERDATE"`    // Дата Оферты
	BuybackDate     NullString  `json:"BUYBACKDATE"`  // дата обратного выкупа
	YieldToMaturity NullFloat64 `json:"YIELDCLOSE"`   // Доходность к погашению при покупке
	YieldToOffer    NullFloat64 `json:"YIELDTOOFFER"` // Доходность к оферте при покупке
	FaceValue       NullFloat64 `json:"FACEVALUE"`
	FaceUnit        NullString  `json:"FACEUNIT"` // номинальная стоимость облигации
	Duration        NullFloat64 `json:"DURATION"` // дюрация (средневзвешенный срок платежей)

}

func (v *Values) UnmarshalJSON(data []byte) error {
	// const op = "service.Values.UnmarshalJSON"

	// Парсим как массив
	var dataSlice []any
	if err := json.Unmarshal(data, &dataSlice); err != nil {
		return fmt.Errorf("cannot unmarshal array: %w", err)
	}

	if len(dataSlice) < 10 {
		return fmt.Errorf("expected at least 10 elements in array, got %d", len(dataSlice))
	}

	tradeDate, err := parseNullString(dataSlice[0])
	if err != nil {
		return fmt.Errorf("element 0 (TRADEDATE): %w", err)
	}
	v.TradeDate = tradeDate

	maturityDate, err := parseNullString(dataSlice[1])
	if err != nil {
		return fmt.Errorf("element 1 (MATDATE): %w", err)
	}
	v.MaturityDate = maturityDate

	offerDate, err := parseNullString(dataSlice[2])
	if err != nil {
		return fmt.Errorf("element 2 (OFFERDATE): %w", err)
	}
	v.OfferDate = offerDate

	buybackDate, err := parseNullString(dataSlice[3])
	if err != nil {
		return fmt.Errorf("element 3 (BUYBACKDATE): %w", err)
	}
	v.BuybackDate = buybackDate

	yieldToMaturity, err := parseNullFloat64(dataSlice[4])
	if err != nil {
		return fmt.Errorf("element 4 (YIELDCLOSE): %w", err)
	}
	v.YieldToMaturity = yieldToMaturity

	yieldToOffer, err := parseNullFloat64(dataSlice[5])
	if err != nil {
		return fmt.Errorf("element 5 (YIELDTOOFFER): %w", err)
	}
	v.YieldToOffer = yieldToOffer

	faceValue, err := parseNullFloat64(dataSlice[6])
	if err != nil {
		return fmt.Errorf("element 6 (FACEVALUE): %w", err)
	}
	v.FaceValue = faceValue

	faceUnit, err := parseNullString(dataSlice[7])
	if err != nil {
		return fmt.Errorf("element 7 (FACEUNIT): %w", err)
	}
	v.FaceUnit = faceUnit

	duration, err := parseNullFloat64(dataSlice[8])
	if err != nil {
		return fmt.Errorf("element 8 (DURATION): %w", err)
	}
	v.Duration = duration

	shortName, err := parseNullString(dataSlice[9])
	if err != nil {
		return fmt.Errorf("element 9 (SHORTNAME): %w", err)
	}
	v.ShortName = shortName

	return nil
}

func parseNullString(input any) (NullString, error) {
	const op = "service.parseNullString"
	var ns NullString
	ns.IsSet = true
	if input == nil {
		ns.IsNull = true
		return ns, nil
	}
	res, ok := input.(string)
	if !ok {
		return NullString{}, fmt.Errorf("op:%s, err: could not convert any to string", op)
	}
	ns.Value = res
	return ns, nil
}

func parseNullFloat64(input any) (NullFloat64, error) {
	const op = "service.parseNullFloat64"
	var nf NullFloat64
	nf.IsSet = true

	if input == nil {
		nf.IsNull = true
		return nf, nil
	}
	res, ok := input.(float64)
	if !ok {
		return NullFloat64{}, fmt.Errorf("op:%s, err: could not convert any to float64", op)
	}
	nf.Value = res
	return nf, nil
}
