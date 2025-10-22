package service

import (
	"encoding/json"
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
	value  string
	isSet  bool
	isNull bool
}

func (ns NullString) Value() string {
	return ns.value
}

func (ns NullString) IsSet() bool {
	return ns.isSet
}

func (ns NullString) IsNull() bool {
	return ns.isNull
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
	ns.isSet = true

	if string(data) == "null" {
		ns.isNull = true
		ns.value = ""
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	ns.value = s
	ns.isNull = false
	return nil
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.isSet {
		return []byte("null"), nil // или возвращать ошибку
	}
	if ns.isNull {
		return []byte("null"), nil
	}
	return json.Marshal(ns.value)
}

type NullFloat64 struct {
	value  float64
	isSet  bool
	isNull bool
}

func (nf NullFloat64) Value() float64 {
	return nf.value
}

func (nf NullFloat64) IsSet() bool {
	return nf.isSet
}

func (nf NullFloat64) IsNull() bool {
	return nf.isNull
}

func (nf *NullFloat64) UnmarshalJSON(data []byte) error {
	nf.isSet = true

	if string(data) == "null" {
		nf.isNull = true
		nf.value = 0
		return nil
	}

	var f float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}

	nf.value = f
	nf.isNull = false
	return nil
}

func (nf NullFloat64) MarshalJSON() ([]byte, error) {
	if !nf.isSet {
		return []byte("null"), nil
	}
	if nf.isNull {
		return []byte("null"), nil
	}
	return json.Marshal(nf.value)
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
	FaceUnit        NullFloat64 `json:"FACEUNIT"` // номинальная стоимость облигации
	Duration        NullFloat64 `json:"DURATION"` // дюрация (средневзвешенный срок платежей)

}
