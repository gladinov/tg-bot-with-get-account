package service

import (
	"errors"
	"fmt"
	"time"

	"main.go/clients/tinkoffApi"
	"main.go/lib/e"
)

var ErrCloseAccount = errors.New("close account haven't portffolio positions")
var ErrNoAcces = errors.New("this token no access to account")

func (c *Client) TinkoffGetPortfolio(account tinkoffApi.Account) (tinkoffApi.Portfolio, error) {
	const op = "service.TinkoffGetPortfolio"
	portfolioRequest := tinkoffApi.PortfolioRequest{
		AccountID:     account.Id,
		AccountStatus: account.Status,
	}
	if account.Status == 3 {
		return tinkoffApi.Portfolio{}, ErrCloseAccount
	}
	if account.AccessLevel == 3 {
		return tinkoffApi.Portfolio{}, ErrNoAcces
	}
	portfolio, err := c.Tinkoffapi.GetPortfolio(portfolioRequest)
	if err != nil {
		return tinkoffApi.Portfolio{}, fmt.Errorf("op:%s, %s", op, err)
	}
	return portfolio, nil
}

func (c *Client) TinkoffGetOperations(accountId string, fromDate time.Time) ([]tinkoffApi.Operation, error) {
	const op = "service.TinkoffGetPortfolio"
	if fromDate.Compare(time.Now()) == 1 {
		return nil, fmt.Errorf("op:%s, from can't be more than the current date", op)
	}
	operationRequest := tinkoffApi.OperationsRequest{
		AccountID: accountId,
		Date:      fromDate,
	}
	tinkoffOperations, err := c.Tinkoffapi.GetOperations(operationRequest)
	if err != nil {
		return nil, e.WrapIfErr(fmt.Sprintf("op:%s,", op), err)
	}
	return tinkoffOperations, nil
}
