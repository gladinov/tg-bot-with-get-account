package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"bonds-report-service/clients/tinkoffApi"

	"bonds-report-service/lib/e"
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

func (c *Client) TinkoffGetPortfolio(ctx context.Context, account tinkoffApi.Account) (_ tinkoffApi.Portfolio, err error) {
	const op = "service.TinkoffGetPortfolio"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)

	}()
	portfolioRequest := tinkoffApi.PortfolioRequest{
		AccountID:     account.Id,
		AccountStatus: account.Status,
	}

	switch account.Status {
	case 0:
		return tinkoffApi.Portfolio{}, ErrUnspecifiedAccount
	case 1:
		return tinkoffApi.Portfolio{}, ErrNewNotOpenYetAccount
	case 3:
		return tinkoffApi.Portfolio{}, ErrCloseAccount
	}

	if account.AccessLevel == 3 {
		return tinkoffApi.Portfolio{}, ErrNoAcces
	}
	if account.Id == "" {
		return tinkoffApi.Portfolio{}, ErrEmptyAccountIdInRequest
	}
	portfolio, err := c.Tinkoffapi.GetPortfolio(ctx, portfolioRequest)
	if err != nil {
		return tinkoffApi.Portfolio{}, fmt.Errorf("op:%s, %s", op, err)
	}
	return portfolio, nil
}

func (c *Client) TinkoffGetOperations(ctx context.Context, accountId string, fromDate time.Time) (_ []tinkoffApi.Operation, err error) {
	const op = "service.TinkoffGetPortfolio"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)

	}()
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
	tinkoffOperations, err := c.Tinkoffapi.GetOperations(ctx, operationRequest)
	if err != nil {
		return nil, e.WrapIfErr(fmt.Sprintf("op:%s,", op), err)
	}
	return tinkoffOperations, nil
}

func (c *Client) TinkoffGetBondActions(ctx context.Context, instrumentUid string) (_ tinkoffApi.BondIdentIdentifiers, err error) {
	const op = "service.TinkoffGetBondActions"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)

	}()

	if instrumentUid == "" {
		return tinkoffApi.BondIdentIdentifiers{}, ErrEmptyInstrumentUid
	}
	bondActions, err := c.Tinkoffapi.GetBondsActions(ctx, instrumentUid)
	if err != nil {
		return tinkoffApi.BondIdentIdentifiers{}, fmt.Errorf("op: %s, error: %s", op, err.Error())
	}
	return bondActions, nil
}

func (c *Client) TinkoffGetFutureBy(ctx context.Context, figi string) (_ tinkoffApi.Future, err error) {
	const op = "service.TinkoffGetFutureBy"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)

	}()
	if figi == "" {
		return tinkoffApi.Future{}, ErrEmptyFigi
	}
	future, err := c.Tinkoffapi.GetFutureBy(ctx, figi)
	if err != nil {
		return tinkoffApi.Future{}, fmt.Errorf("op: %s, error: %s", op, err.Error())
	}
	return future, nil
}

func (c *Client) TinkoffGetBondByUid(ctx context.Context, uid string) (_ tinkoffApi.Bond, err error) {
	const op = "service.TinkoffGetBondByUid"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)

	}()
	if uid == "" {
		return tinkoffApi.Bond{}, ErrEmptyUid
	}
	bond, err := c.Tinkoffapi.GetBondByUid(ctx, uid)
	if err != nil {
		return tinkoffApi.Bond{}, fmt.Errorf("op: %s, error: %s", op, err.Error())
	}
	return bond, nil
}

func (c *Client) TinkoffGetCurrencyBy(ctx context.Context, figi string) (_ tinkoffApi.Currency, err error) {
	const op = "service.TinkoffGetCurrencyBy"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)

	}()
	if figi == "" {
		return tinkoffApi.Currency{}, ErrEmptyFigi
	}
	currency, err := c.Tinkoffapi.GetCurrencyBy(ctx, figi)
	if err != nil {
		return tinkoffApi.Currency{}, fmt.Errorf("op: %s, error: %s", op, err.Error())
	}
	return currency, nil
}

func (c *Client) TinkoffGetBaseShareFutureValute(ctx context.Context, positionUid string) (_ tinkoffApi.BaseShareFutureValuteResponse, err error) {
	const op = "service.TinkoffGetBaseShareFutureValute"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)

	}()

	if positionUid == "" {
		return tinkoffApi.BaseShareFutureValuteResponse{}, ErrEmptyPositionUid
	}

	instrumentsShortResponce, err := c.Tinkoffapi.FindBy(ctx, positionUid)
	if err != nil {
		return tinkoffApi.BaseShareFutureValuteResponse{}, fmt.Errorf("op: %s, error: %s", op, err.Error())
	}
	if len(instrumentsShortResponce) == 0 {
		return tinkoffApi.BaseShareFutureValuteResponse{}, fmt.Errorf("op: %s, error:can't get base share future valute", op)
	}
	instrument := instrumentsShortResponce[0]
	if instrument.InstrumentType != "share" {
		return tinkoffApi.BaseShareFutureValuteResponse{}, fmt.Errorf("op: %s, instrument is not share", op)
	}
	if instrument.Figi == "" {
		return tinkoffApi.BaseShareFutureValuteResponse{}, ErrEmptyFigi
	}
	currency, err := c.Tinkoffapi.GetShareCurrencyBy(ctx, instrument.Figi)
	if err != nil {
		return tinkoffApi.BaseShareFutureValuteResponse{}, fmt.Errorf("op: %s, error: %s", op, err.Error())
	}
	var resp tinkoffApi.BaseShareFutureValuteResponse
	resp.Currency = currency.Currency
	return resp, nil
}

func (c *Client) TinkoffFindBy(ctx context.Context, query string) (_ []tinkoffApi.InstrumentShort, err error) {
	const op = "service.TinkoffFindBy"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)

	}()

	if query == "" {
		return nil, ErrEmptyQuery
	}
	resp, err := c.Tinkoffapi.FindBy(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("op: %s, error: %s", op, err.Error())
	}
	return resp, nil
}

func (c *Client) TinkoffGetLastPriceInPersentageToNominal(ctx context.Context, instrumentUid string) (_ tinkoffApi.LastPriceResponse, err error) {
	const op = "service.TinkoffGetLastPriceInPersentageToNominal"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)

	}()

	if instrumentUid == "" {
		return tinkoffApi.LastPriceResponse{}, ErrEmptyInstrumentUid
	}
	lastPrice, err := c.Tinkoffapi.GetLastPriceInPersentageToNominal(ctx, instrumentUid)
	if err != nil {
		return tinkoffApi.LastPriceResponse{}, fmt.Errorf("op: %s, error: %s", op, err.Error())
	}
	return lastPrice, nil
}
