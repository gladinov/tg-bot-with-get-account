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
var ErrEmptyAccountIdInRequest = errors.New("accountId could not be empty")
var ErrUnspecifiedAccount = errors.New("account is unspecified")
var ErrNewNotOpenYetAccount = errors.New("accountId is not opened yet")
var ErrEmptyInstrumentUid = errors.New("instrumentUid could not be empty string")
var ErrEmptyFigi = errors.New("figi could not be empty string")
var ErrEmptyQuery = errors.New("query could not be empty")
var ErrEmptyUid = errors.New("uid could not be empty string")
var ErrEmptyPositionUid = errors.New("positionUid could not be empty string")

func (c *Client) TinkoffGetPortfolio(account tinkoffApi.Account) (tinkoffApi.Portfolio, error) {
	const op = "service.TinkoffGetPortfolio"
	portfolioRequest := tinkoffApi.PortfolioRequest{
		AccountID:     account.Id,
		AccountStatus: account.Status,
	}
	switch account.Status {
	case 0:
		return tinkoffApi.Portfolio{}, ErrCloseAccount
	case 1:
		return tinkoffApi.Portfolio{}, ErrUnspecifiedAccount
	case 3:
		return tinkoffApi.Portfolio{}, ErrNewNotOpenYetAccount
	}

	if account.AccessLevel == 3 {
		return tinkoffApi.Portfolio{}, ErrNoAcces
	}
	if account.Id == "" {
		return tinkoffApi.Portfolio{}, ErrEmptyAccountIdInRequest
	}
	portfolio, err := c.Tinkoffapi.GetPortfolio(portfolioRequest)
	if err != nil {
		return tinkoffApi.Portfolio{}, fmt.Errorf("op:%s, %s", op, err)
	}
	return portfolio, nil
}

func (c *Client) TinkoffGetOperations(accountId string, fromDate time.Time) ([]tinkoffApi.Operation, error) {
	const op = "service.TinkoffGetPortfolio"
	now := time.Now().UTC()
	if accountId == "" {
		return nil, fmt.Errorf("%s: empty account ID", op)
	}
	if fromDate.After(now) {
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

func (c *Client) TinkoffGetBondActions(instrumentUid string) (tinkoffApi.BondIdentIdentifiers, error) {
	const op = "service.TinkoffGetBondActions"
	if instrumentUid == "" {
		return tinkoffApi.BondIdentIdentifiers{}, ErrEmptyInstrumentUid
	}
	bondActions, err := c.Tinkoffapi.GetBondsActions(instrumentUid)
	if err != nil {
		return tinkoffApi.BondIdentIdentifiers{}, fmt.Errorf("op: %s, error: %s", op, err.Error())
	}
	return bondActions, nil
}

func (c *Client) TinkoffGetFutureBy(figi string) (tinkoffApi.Future, error) {
	const op = "service.TinkoffGetFutureBy"
	if figi == "" {
		return tinkoffApi.Future{}, ErrEmptyFigi
	}
	future, err := c.Tinkoffapi.GetFutureBy(figi)
	if err != nil {
		return tinkoffApi.Future{}, fmt.Errorf("op: %s, error: %s", op, err.Error())
	}
	return future, nil
}

func (c *Client) TinkoffGetBondByUid(uid string) (tinkoffApi.Bond, error) {
	const op = "service.TinkoffGetBondByUid"
	if uid == "" {
		return tinkoffApi.Bond{}, ErrEmptyUid
	}
	bond, err := c.Tinkoffapi.GetBondByUid(uid)
	if err != nil {
		return tinkoffApi.Bond{}, fmt.Errorf("op: %s, error: %s", op, err.Error())
	}
	return bond, nil
}

func (c *Client) TinkoffGetCurrencyBy(figi string) (tinkoffApi.Currency, error) {
	const op = "service.TinkoffGetCurrencyBy"
	if figi == "" {
		return tinkoffApi.Currency{}, ErrEmptyFigi
	}
	currency, err := c.Tinkoffapi.GetCurrencyBy(figi)
	if err != nil {
		return tinkoffApi.Currency{}, fmt.Errorf("op: %s, error: %s", op, err.Error())
	}
	return currency, nil
}

func (c *Client) TinkoffGetBaseShareFutureValute(positionUid string) (tinkoffApi.BaseShareFutureValuteResponse, error) {
	const op = "service.TinkoffGetBaseShareFutureValute"
	if positionUid == "" {
		return tinkoffApi.BaseShareFutureValuteResponse{}, ErrEmptyPositionUid
	}
	currency, err := c.Tinkoffapi.GetBaseShareFutureValute(positionUid)
	if err != nil {
		return tinkoffApi.BaseShareFutureValuteResponse{}, fmt.Errorf("op: %s, error: %s", op, err.Error())
	}
	return currency, nil
}

func (c *Client) TinkoffFindBy(query string) ([]tinkoffApi.InstrumentShort, error) {
	const op = "service.TinkoffFindBy"
	if query == "" {
		return nil, ErrEmptyQuery
	}
	resp, err := c.Tinkoffapi.FindBy(query)
	if err != nil {
		return nil, fmt.Errorf("op: %s, error: %s", op, err.Error())
	}
	return resp, nil
}

func (c *Client) TinkoffGetLastPriceInPersentageToNominal(instrumentUid string) (tinkoffApi.LastPriceResponse, error) {
	const op = "service.TinkoffGetLastPriceInPersentageToNominal"
	if instrumentUid == "" {
		return tinkoffApi.LastPriceResponse{}, ErrEmptyInstrumentUid
	}
	lastPrice, err := c.Tinkoffapi.GetLastPriceInPersentageToNominal(instrumentUid)
	if err != nil {
		return tinkoffApi.LastPriceResponse{}, fmt.Errorf("op: %s, error: %s", op, err.Error())
	}
	return lastPrice, nil
}
