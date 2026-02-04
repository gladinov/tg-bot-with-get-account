package models

import "time"

type SpecificationsRequest struct {
	Ticker string    `json:"ticker"`
	Date   time.Time `json:"date"`
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

func (nf NullFloat64) IsHasValue() bool {
	if !nf.IsSet || nf.IsNull {
		return false
	}
	return true
}

func (ns NullString) IsHasValue() bool {
	if !ns.IsSet || ns.IsNull {
		return false
	}
	return true
}
