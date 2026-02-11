package service

import (
	"bonds-report-service/internal/models/domain"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/gladinov/e"
)

var (
	ErrCloseAccount            = errors.New("close account haven't portffolio positions")
	ErrNoAcces                 = errors.New("this token no access to account")
	ErrEmptyAccountIdInRequest = errors.New("accountId could not be empty")
	ErrUnspecifiedAccount      = errors.New("account is unspecified")
	ErrNewNotOpenYetAccount    = errors.New("accountId is not opened yet")
	ErrEmptyInstrumentUid      = errors.New("instrumentUid could not be empty string")
	ErrEmptyFigi               = errors.New("figi could not be empty string")
	ErrEmptyQuery              = errors.New("query could not be empty")
	ErrEmptyUid                = errors.New("uid could not be empty string")
	ErrEmptyPositionUid        = errors.New("positionUid could not be empty string")
)

func (c *Client) TinkoffGetPortfolio(ctx context.Context, account domain.Account) (_ domain.Portfolio, err error) {
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

	switch account.Status {
	case 0:
		return domain.Portfolio{}, ErrUnspecifiedAccount
	case 1:
		return domain.Portfolio{}, ErrNewNotOpenYetAccount
	case 3:
		return domain.Portfolio{}, ErrCloseAccount
	}

	if account.AccessLevel == 3 {
		return domain.Portfolio{}, ErrNoAcces
	}
	if account.ID == "" {
		return domain.Portfolio{}, ErrEmptyAccountIdInRequest
	}
	portfolio, err := c.Tinkoffapi.PortfolioTinkoffClient.GetPortfolio(ctx, account.ID, account.Status)
	if err != nil {
		return domain.Portfolio{}, fmt.Errorf("op:%s, %s", op, err)
	}
	return portfolio, nil
}

func (c *Client) TinkoffGetOperations(ctx context.Context, accountId string, fromDate time.Time) (_ []domain.Operation, err error) {
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
	tinkoffOperations, err := c.Tinkoffapi.PortfolioTinkoffClient.GetOperations(ctx, accountId, fromDate)
	if err != nil {
		return nil, e.WrapIfErr(fmt.Sprintf("op:%s,", op), err)
	}
	return tinkoffOperations, nil
}

func (c *Client) TinkoffGetBondActions(ctx context.Context, instrumentUid string) (_ domain.BondIdentIdentifiers, err error) {
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
		return domain.BondIdentIdentifiers{}, ErrEmptyInstrumentUid
	}
	bondActions, err := c.Tinkoffapi.AnalyticsTinkoffClient.GetBondsActions(ctx, instrumentUid)
	if err != nil {
		return domain.BondIdentIdentifiers{}, fmt.Errorf("op: %s, error: %s", op, err.Error())
	}
	return bondActions, nil
}

func (c *Client) TinkoffGetFutureBy(ctx context.Context, figi string) (_ domain.Future, err error) {
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
		return domain.Future{}, ErrEmptyFigi
	}
	future, err := c.Tinkoffapi.InstrumentsTinkoffClient.GetFutureBy(ctx, figi)
	if err != nil {
		return domain.Future{}, fmt.Errorf("op: %s, error: %s", op, err.Error())
	}
	return future, nil
}

func (c *Client) TinkoffGetBondByUid(ctx context.Context, uid string) (_ domain.Bond, err error) {
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
		return domain.Bond{}, ErrEmptyUid
	}
	bond, err := c.Tinkoffapi.InstrumentsTinkoffClient.GetBondByUid(ctx, uid)
	if err != nil {
		return domain.Bond{}, fmt.Errorf("op: %s, error: %s", op, err.Error())
	}
	return bond, nil
}

func (c *Client) TinkoffGetCurrencyBy(ctx context.Context, figi string) (_ domain.Currency, err error) {
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
		return domain.Currency{}, ErrEmptyFigi
	}
	currency, err := c.Tinkoffapi.InstrumentsTinkoffClient.GetCurrencyBy(ctx, figi)
	if err != nil {
		return domain.Currency{}, fmt.Errorf("op: %s, error: %s", op, err.Error())
	}
	return currency, nil
}

func (c *Client) TinkoffGetBaseShareFutureValute(ctx context.Context, positionUid string) (_ domain.ShareCurrency, err error) {
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
		return domain.ShareCurrency{}, ErrEmptyPositionUid
	}

	instrumentsShortResponce, err := c.Tinkoffapi.InstrumentsTinkoffClient.FindBy(ctx, positionUid)
	if err != nil {
		return domain.ShareCurrency{}, fmt.Errorf("op: %s, error: %s", op, err.Error())
	}
	if len(instrumentsShortResponce) == 0 {
		return domain.ShareCurrency{}, fmt.Errorf("op: %s, error:can't get base share future valute", op)
	}
	instrument := instrumentsShortResponce[0]
	if instrument.InstrumentType != "share" {
		return domain.ShareCurrency{}, fmt.Errorf("op: %s, instrument is not share", op)
	}
	if instrument.Figi == "" {
		return domain.ShareCurrency{}, ErrEmptyFigi
	}
	currency, err := c.Tinkoffapi.InstrumentsTinkoffClient.GetShareCurrencyBy(ctx, instrument.Figi)
	if err != nil {
		return domain.ShareCurrency{}, fmt.Errorf("op: %s, error: %s", op, err.Error())
	}

	return currency, nil
}

func (c *Client) TinkoffFindBy(ctx context.Context, query string) (_ []domain.InstrumentShort, err error) {
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
	resp, err := c.Tinkoffapi.InstrumentsTinkoffClient.FindBy(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("op: %s, error: %s", op, err.Error())
	}
	return resp, nil
}

func (c *Client) TinkoffGetLastPriceInPersentageToNominal(ctx context.Context, instrumentUid string) (_ domain.LastPrice, err error) {
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
		return domain.LastPrice{}, ErrEmptyInstrumentUid
	}
	lastPrice, err := c.Tinkoffapi.AnalyticsTinkoffClient.GetLastPriceInPersentageToNominal(ctx, instrumentUid)
	if err != nil {
		return domain.LastPrice{}, fmt.Errorf("op: %s, error: %s", op, err.Error())
	}
	return lastPrice, nil
}
