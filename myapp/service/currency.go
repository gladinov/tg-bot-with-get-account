package service

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"main.go/clients/cbr"
	"main.go/lib/e"
	"main.go/service/service_models"
)

func (c *Client) GetCurrencyFromCB(charCode string, date time.Time) (vunit_rate float64, err error) {
	defer func() { err = e.WrapIfErr("can't get currencies from CB", err) }()
	vunit_rate, err = c.Storage.GetCurrency(context.Background(), charCode, date)
	if err == nil {
		return vunit_rate, nil
	}
	if err != nil && !errors.Is(err, service_models.ErrNoCurrency) {
		return vunit_rate, err
	}
	currenciesFromCB, err := c.CbrApi.GetAllCurrencies(date)
	if err != nil {
		return vunit_rate, err
	}
	currencies, err := transformCurrenciesFromCB(currenciesFromCB)
	if err != nil {
		return vunit_rate, err
	}
	charCode = strings.ToLower(charCode)
	vunit_rate = currencies.CurrenciesMap[charCode].VunitRate

	err = c.Storage.SaveCurrency(context.Background(), *currencies, date)
	if err != nil {
		return vunit_rate, err
	}
	return vunit_rate, nil

}

func transformCurrenciesFromCB(inCurrencies cbr.CurrenciesResponce) (_ *service_models.Currencies, err error) {
	defer func() { err = e.WrapIfErr("can't transform currencies from CB", err) }()
	outCurrencies := &service_models.Currencies{
		CurrenciesMap: make(map[string]service_models.Currency),
	}

	date, err := time.Parse(layoutCurr, inCurrencies.Date)
	if err != nil {
		return outCurrencies, err
	}
	for _, inCurrency := range inCurrencies.Currencies {
		var outCurrency service_models.Currency
		outCurrency.Date = date
		outCurrency.NumCode = inCurrency.NumCode
		outCurrency.CharCode = strings.ToLower(inCurrency.CharCode)
		outCurrency.Nominal, err = strconv.Atoi(inCurrency.Nominal)
		if err != nil {
			return outCurrencies, err
		}
		outCurrency.Name = inCurrency.Name
		value := strings.ReplaceAll(inCurrency.Value, ",", ".")
		outCurrency.Value, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return outCurrencies, err
		}
		vunitRate := strings.ReplaceAll(inCurrency.VunitRate, ",", ".")
		outCurrency.VunitRate, err = strconv.ParseFloat(vunitRate, 64)
		if err != nil {
			return outCurrencies, err
		}
		charCode := strings.ToLower(inCurrency.CharCode)
		outCurrencies.CurrenciesMap[charCode] = outCurrency

	}
	return outCurrencies, nil
}
