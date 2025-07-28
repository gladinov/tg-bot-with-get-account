package tinkoffApi

import (
	"errors"
	"time"

	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
)

type Account struct {
	Id          string
	Type        pb.AccountType
	Name        string
	Status      int64
	OpenedDate  time.Time
	ClosedDate  time.Time
	AccessLevel pb.AccessLevel
	Portfolio   []*pb.PortfolioPosition
}

func (c *Client) GetAcc() (map[string]Account, error) {
	usersService := c.Client.NewUsersServiceClient()
	accounts := make(map[string]Account)
	status := pb.AccountStatus_ACCOUNT_STATUS_ALL // ПОтом надо обработать закрытые счета(например ИИС)
	accsResp, err := usersService.GetAccounts(&status)
	if err != nil {
		return nil, errors.New("GetAcc: operationsService.GetOperationsByCursor" + err.Error())
	} else {
		accs := accsResp.GetAccounts()
		for _, acc := range accs {
			account := Account{Id: acc.GetId(),
				Name:       acc.GetName(),
				OpenedDate: acc.GetOpenedDate().AsTime(),
				ClosedDate: acc.GetClosedDate().AsTime(),
				Status:     int64(acc.GetStatus()),
			}
			accounts[acc.GetId()] = account
		}
	}

	return accounts, nil
}
