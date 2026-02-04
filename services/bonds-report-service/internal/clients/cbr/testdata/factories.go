package factories

import (
	"bonds-report-service/internal/models"
	"encoding/json"
	"time"
)

func NewCurrencyRequest(date ...time.Time) *models.CurrencyRequest {
	d := time.Now()
	if len(date) > 0 {
		d = date[0]
	}
	return models.NewCurrencyRequest(d)
}

func NewCurrency() models.Currency {
	c := models.Currency{
		NumCode:   "840",
		CharCode:  "USD",
		Nominal:   "1",
		Name:      "US Dollar",
		Value:     "75.12",
		VunitRate: "1.0",
	}
	return c
}

func NewCurrenciesResponce(date string, currencies ...models.Currency) *models.CurrenciesResponce {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	return models.NewCurrenciesResponce(date, currencies)
}

func NewHTTPResponse(statusCode int, body interface{}) *models.HTTPResponse {
	var bodyBytes []byte
	switch v := body.(type) {
	case string:
		bodyBytes = []byte(v)
	default:
		b, _ := json.Marshal(v)
		bodyBytes = b
	}
	return models.NewHTTPResponse(statusCode, bodyBytes)
}
