package models

import (
	"time"
)

type CurrencyRequest struct {
	Date time.Time `json:"date,omitempty"`
}

func NewCurrencyRequest(date time.Time) *CurrencyRequest {
	return &CurrencyRequest{Date: date}
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

func NewCurrenciesResponce(date string, currs []Currency) *CurrenciesResponce {
	return &CurrenciesResponce{
		Date:       date,
		Currencies: currs,
	}
}

type HTTPResponse struct {
	StatusCode int
	Body       []byte
}

func NewHTTPResponse(statusCode int, body []byte) *HTTPResponse {
	return &HTTPResponse{
		StatusCode: statusCode,
		Body:       body,
	}
}
