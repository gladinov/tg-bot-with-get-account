package tinkoffapi

import (
	"fmt"
	"time"

	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
	"main.go/lib/e"
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
	Operations  []*pb.OperationItem
}

type User struct {
	Token    string
	Accounts []*Account
}

func GetAcc(client *investgo.Client) (string, error) {
	usersService := client.NewUsersServiceClient()
	status := pb.AccountStatus_ACCOUNT_STATUS_ALL // ПОтом надо обработать закрытые счета(например ИИС)
	accsResp, err := usersService.GetAccounts(&status)
	var accStr string = "По данному аккаунту доступны следующие счета:"
	if err != nil {
		return "", e.Wrap("getAcc err", err)
	} else {
		accs := accsResp.GetAccounts()
		for _, acc := range accs {
			account := Account{
				Id:          acc.GetId(),
				Type:        acc.GetType(),
				Name:        acc.GetName(),
				Status:      int64(acc.GetStatus()),
				OpenedDate:  acc.GetOpenedDate().AsTime(),
				ClosedDate:  acc.GetClosedDate().AsTime(),
				AccessLevel: acc.GetAccessLevel(),
			}
			accStr += fmt.Sprintf("\n ID:%v, Type: %s, Name: %s, Status: %v \n", account.Id, account.Type, account.Name, account.Status)
		}
	}

	return accStr, nil
}
