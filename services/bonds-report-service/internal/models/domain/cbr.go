package domain

import "time"

type CurrencyCBR struct {
	Date      time.Time
	NumCode   string
	CharCode  string
	Nominal   int
	Name      string
	Value     float64
	VunitRate float64
}

type CurrenciesCBR struct {
	CurrenciesMap map[string]CurrencyCBR
}

func NewCurrencies(mapCurr map[string]CurrencyCBR) *CurrenciesCBR {
	return &CurrenciesCBR{
		CurrenciesMap: mapCurr,
	}
}
