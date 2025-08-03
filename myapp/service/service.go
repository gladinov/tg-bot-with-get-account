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
	reportPath = "service/vizualize/tables/report.png"
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

func (c *Client) GetBondReportsByFifo(chatID int, token string) (err error) {
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
		err = c.updateOperations(chatID, account.Id, account.OpenedDate)
		if err != nil {
			return err
		}

		portfolio, err := c.Tinkoffapi.GetPortf(account.Id, account.Status)
		if err != nil && !errors.Is(err, tinkoffApi.ErrCloseAccount) {
			return err
		}

		portfolioPositions, err := c.TransformPositions(account.Id, portfolio.Positions)
		if err != nil {
			return err
		}
		err = c.Storage.DeleteBondReport(context.Background(), chatID, account.Id)
		if err != nil {
			return err
		}
		bondsInRub := make([]service_models.BondReport, 0)
		for _, v := range portfolioPositions {
			if v.InstrumentType == "bond" {
				operationsDb, err := c.Storage.GetOperations(context.Background(), chatID, v.AssetUid, account.Id)
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
				bondsInRub = append(bondsInRub, bondReport.BondsInRUB...)
			}
		}
		err = c.Storage.SaveBondReport(context.Background(), chatID, account.Id, bondsInRub)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) GetBondReportsWithEachGeneralPosition(chatID int, token string) (err error) {
	defer func() { err = e.WrapIfErr("can't get general bond report", err) }()
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
		err = c.updateOperations(chatID, account.Id, account.OpenedDate)
		if err != nil {
			return err
		}

		portfolio, err := c.Tinkoffapi.GetPortf(account.Id, account.Status)
		if err != nil && !errors.Is(err, tinkoffApi.ErrCloseAccount) {
			return err
		}

		portfolioPositions, err := c.TransformPositions(account.Id, portfolio.Positions)
		if err != nil {
			return err
		}
		err = c.Storage.DeleteGeneralBondReport(context.Background(), chatID, account.Id)
		if err != nil {
			return err
		}
		bondsInRub := make([]service_models.GeneralBondReporPosition, 0)
		for _, v := range portfolioPositions {
			if v.InstrumentType == "bond" {
				operationsDb, err := c.Storage.GetOperations(context.Background(), chatID, v.AssetUid, account.Id)
				if err != nil {
					return err
				}
				resultBondPosition, err := c.ProcessOperations(operationsDb)
				if err != nil {
					return err
				}
				bondReport, err := c.CreateGeneralBondReport(resultBondPosition, portfolio.TotalAmount)
				if err != nil {
					return err
				}

				bondsInRub = append(bondsInRub, bondReport)
			}
		}
		err = c.Storage.SaveGeneralBondReport(context.Background(), chatID, account.Id, bondsInRub)
		if err != nil {
			return err
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

func (c *Client) updateOperations(chatID int, accountId string, openDate time.Time) (err error) {
	defer func() { err = e.WrapIfErr("can't updateOperations", err) }()
	fromDate, err := c.Storage.LastOperationTime(context.Background(), chatID, accountId)
	fromDate = fromDate.Add(time.Microsecond * 1)

	if err != nil {
		if errors.Is(err, service_models.ErrNoOpperations) {
			fromDate = openDate
		} else {
			return err
		}
	}

	tinkoffOperations, err := c.Tinkoffapi.GetOperations(accountId, fromDate)
	if err != nil {
		return err
	}
	operations := c.TransOperations(tinkoffOperations)

	err = c.Storage.SaveOperations(context.Background(), chatID, accountId, operations)
	if err != nil {
		return err
	}
	return nil
}
