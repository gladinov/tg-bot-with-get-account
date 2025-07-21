package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"main.go/clients/cbr"
	"main.go/clients/moex"
	"main.go/clients/tinkoffApi"
	"main.go/lib/e"
	"main.go/service/service_models"
	service_storage "main.go/service/storage"
)

const (
	layoutCurr = "02.01.2006"
)

type Client struct {
	Tinkoffapi *tinkoffApi.Client
	MoexApi    *moex.Client
	CbrApi     *cbr.Client
	Storage    service_storage.Storage
}

func New(tinkoffApiClient *tinkoffApi.Client, moexClient *moex.Client, CbrClient *cbr.Client, storage service_storage.Storage) *Client {
	return &Client{
		Tinkoffapi: tinkoffApiClient,
		MoexApi:    moexClient,
		CbrApi:     CbrClient,
		Storage:    storage,
	}
}

func (c *Client) GetBondReports(chatID int, token string) (err error) {
	defer func() { err = e.WrapIfErr("can't get bond reports", err) }()
	client := c.Tinkoffapi

	err = client.FillClient(token)
	if err != nil {
		return err
	}

	accounts, err := c.Tinkoffapi.GetAcc()
	if err != nil {
		return err
	}

	for _, account := range accounts {
		err := c.Tinkoffapi.GetOpp(&account)
		if err != nil {
			return err
		}
		operations := c.TransOperations(account.Operations)

		err = c.Storage.SaveOperations(context.Background(), chatID, account.Id, operations)
		if err != nil {
			return err
		}

		err = c.Tinkoffapi.GetPortf(&account)
		if err != nil {
			return err
		}

		portfolio, err := c.TransPositions(&account)
		if err != nil {
			return err
		}

		for _, v := range portfolio.BondPositions {
			operationsDb, err := c.Storage.GetOperations(context.Background(), chatID, v.Identifiers.AssetUid, account.Id)
			if err != nil {
				return err
			}
			resultBondPosition, err := c.ProcessOperations(operationsDb)
			if err != nil {
				return err
			}
			bondReport, err := c.CreateBondReport(*resultBondPosition)
			if err != nil {
				return err
			}
			err = c.Storage.SaveBondReport(context.Background(), chatID, account.Id, bondReport.BondsInRUB)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Client) GetAccounts(token string) (answ string, err error) {
	defer func() { err = e.WrapIfErr("can't get accounts", err) }()
	client := c.Tinkoffapi
	var accStr string = "По данному аккаунту доступны следующие счета:"
	err = client.FillClient(token)
	if err != nil {
		return "", err
	}

	accs, err := c.Tinkoffapi.GetAcc()
	if err != nil {
		return "", err
	}
	for _, account := range accs {
		accStr += fmt.Sprintf("\n ID:%v, Type: %s, Name: %s, Status: %v \n", account.Id, account.Type, account.Name, account.Status)
	}

	return accStr, nil
}

func (c *Client) GetUsd() (float64, error) {
	usd, err := c.GetCurrencyFromCB("usd", time.Now())
	if err != nil {
		return 0, e.WrapIfErr("usd get error", err)
	}
	return usd, nil
}

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

	err = c.Storage.SaveCurrency(context.Background(), *currencies)
	if err != nil {
		return vunit_rate, err
	}
	return vunit_rate, nil

}

func transformCurrenciesFromCB(inCurrencies cbr.ValCurs) (_ *service_models.Currencies, err error) {
	defer func() { err = e.WrapIfErr("can't transform currencies from CB", err) }()
	outCurrencies := &service_models.Currencies{
		CurrenciesMap: make(map[string]service_models.Currency),
	}

	date, err := time.Parse(layoutCurr, inCurrencies.Date)
	if err != nil {
		return outCurrencies, err
	}
	for _, inCurrency := range inCurrencies.Valute {
		var outCurrency service_models.Currency
		outCurrency.Date = date
		outCurrency.NumCode = inCurrency.NumCode
		outCurrency.CharCode = strings.ToLower(inCurrency.CharCode)
		outCurrency.Nominal, err = strconv.Atoi(inCurrency.Nominal)
		if err != nil {
			return outCurrencies, err
		}
		outCurrency.Name = inCurrency.Name
		value := strings.Replace(inCurrency.Value, ",", ".", -1)
		outCurrency.Value, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return outCurrencies, err
		}
		vunitRate := strings.Replace(inCurrency.VunitRate, ",", ".", -1)
		outCurrency.VunitRate, err = strconv.ParseFloat(vunitRate, 64)
		if err != nil {
			return outCurrencies, err
		}
		charCode := strings.ToLower(inCurrency.CharCode)
		outCurrencies.CurrenciesMap[charCode] = outCurrency

	}
	return outCurrencies, nil
}
