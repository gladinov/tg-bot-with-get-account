package service

import (
	"context"
	"errors"
	"fmt"
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
		fromDate, err := c.Storage.LastOperationTime(context.Background(), chatID, account.Id)
		fromDate = fromDate.Add(time.Microsecond * 1)

		if err != nil {
			if errors.Is(err, service_models.ErrNoOpperations) {
				fromDate = account.OpenedDate
			} else {
				return err
			}
		}

		err = c.Tinkoffapi.GetOpp(&account, fromDate)
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
		err = c.Storage.DeleteBondReport(context.Background(), chatID, account.Id)
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
