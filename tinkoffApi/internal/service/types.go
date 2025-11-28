package service

import (
	"context"
	"errors"
)

var ErrCloseAccount = errors.New("close account haven't portffolio positions")
var ErrNoAcces = errors.New("this token no access to account")
var ErrEmptyAccountIdInRequest = errors.New("accountId could not be empty")
var ErrUnspecifiedAccount = errors.New("account is unspecified")
var ErrNewNotOpenYetAccount = errors.New("accountId is not opened yet")
var ErrEmptyQuery = errors.New("query could not be empty")
var ErrEmptyFigi = errors.New("figi could not be empty string")
var ErrEmptyUid = errors.New("uid could not be empty string")
var ErrEmptyPositionUid = errors.New("positionUid could not be empty string")
var ErrEmptyInstrumentUid = errors.New("instrumentUid could not be empty string")

type Service struct {
	InstrumentService InstrumentService
	PortfolioService  PortfolioService
	AnalyticsService  AnalyticsService
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=InstrumentService
type InstrumentService interface {
	GetClient(ctx context.Context) error
	FindBy(query string) ([]InstrumentShort, error)
	GetBondByUid(uid string) (Bond, error)
	GetCurrencyBy(figi string) (Currency, error)
	GetFutureBy(figi string) (Future, error)
	GetShareCurrencyBy(figi string) (ShareCurrencyByResponse, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=PortfolioService
type PortfolioService interface {
	GetClient(ctx context.Context) error
	GetAccounts() (map[string]Account, error)
	GetPortfolio(request PortfolioRequest) (Portfolio, error)
	GetOperations(request OperationsRequest) ([]Operation, error)
	MakeSafeGetOperationsRequest(request OperationsRequest) ([]Operation, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=AnalyticsService
type AnalyticsService interface {
	GetClient(ctx context.Context) error
	GetLastPriceInPersentageToNominal(instrumentUid string) (LastPriceResponse, error)
	GetAllAssetUids() (map[string]string, error)
	GetBondsActions(instrumentUid string) (BondIdentIdentifiers, error)
}
