package moex

import (
	"encoding/json"
	"errors"
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

type Values struct {
	ShortName       *string  `json:"SHORTNAME"`
	TradeDate       *string  `json:"TRADEDATE"`    // Торговая дата(на момент которой рассчитаны остальные данные)
	MaturityDate    *string  `json:"MATDATE"`      // Дата погашения
	OfferDate       *string  `json:"OFFERDATE"`    // Дата Оферты
	BuybackDate     *string  `json:"BUYBACKDATE"`  // дата обратного выкупа
	YieldToMaturity *float64 `json:"YIELDCLOSE"`   // Доходность к погашению при покупке
	YieldToOffer    *float64 `json:"YIELDTOOFFER"` // Доходность к оферте при покупке
	FaceValue       *float64 `json:"FACEVALUE"`
	FaceUnit        *float64 `json:"FACEUNIT"` // номинальная стоимость облигации
	Duration        *float64 `json:"DURATION"` // дюрация (средневзвешенный срок платежей)

}

func (d *Values) UnmarshalJSON(data []byte) error {
	dataSlice := make([]any, 10)
	err := json.Unmarshal(data, &dataSlice)
	if err != nil {
		return errors.New("CustomFloat64: UnmarshalJSON: " + err.Error())
	}
	d.TradeDate = checkStringNull(dataSlice[0])
	d.MaturityDate = checkStringNull(dataSlice[1])
	d.OfferDate = checkStringNull(dataSlice[2])
	d.BuybackDate = checkStringNull(dataSlice[3])
	d.YieldToMaturity = checkFloat64Null(dataSlice[4])
	d.YieldToOffer = checkFloat64Null(dataSlice[5])
	d.FaceValue = checkFloat64Null(dataSlice[6])
	d.FaceUnit = checkFloat64Null(dataSlice[7])
	d.Duration = checkFloat64Null(dataSlice[8])
	d.ShortName = checkStringNull(dataSlice[9])

	return nil
}

func checkFloat64Null(a any) *float64 {
	if FloatVal, ok := a.(float64); ok {
		return &FloatVal
	} else {
		return nil
	}
}

func checkStringNull(a any) *string {
	if StringVal, ok := a.(string); ok {
		return &StringVal
	} else {
		return nil
	}
}
