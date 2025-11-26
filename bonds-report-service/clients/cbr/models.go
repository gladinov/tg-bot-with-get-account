package cbr

import (
	"time"
)

type CurrencyRequest struct {
	Date time.Time `json:"date,omitempty"`
}

type Currency struct {
	NumCode   string `json:"numCode,omitempty"`
	CharCode  string `json:"charCode,omitempty"`
	Nominal   string `json:"nominal,omitempty"`
	Name      string `json:"name,omitempty"`
	Value     string `json:"value,omitempty"`
	VunitRate string `json:"vunitRate,omitempty"`
}

type CurrenciesResponce struct {
	Date       string     `json:"date,omitempty"`
	Currencies []Currency `json:"valute,omitempty"`
}

type HTTPResponse struct {
	StatusCode int
	Body       []byte
}
