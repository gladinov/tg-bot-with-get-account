package service

import (
	"context"
	"errors"

	"github.com/russianinvestments/invest-api-go-sdk/investgo"
)

var (
	ErrCloseAccount            = errors.New("close account haven't portffolio positions")
	ErrNoAcces                 = errors.New("this token no access to account")
	ErrEmptyAccountIdInRequest = errors.New("accountId could not be empty")
	ErrUnspecifiedAccount      = errors.New("account is unspecified")
	ErrNewNotOpenYetAccount    = errors.New("accountId is not opened yet")
	ErrEmptyQuery              = errors.New("query could not be empty")
	ErrEmptyFigi               = errors.New("figi could not be empty string")
	ErrEmptyUid                = errors.New("uid could not be empty string")
	ErrEmptyPositionUid        = errors.New("positionUid could not be empty string")
	ErrEmptyInstrumentUid      = errors.New("instrumentUid could not be empty string")
)

type Service struct {
	InstrumentService InstrumentService
	PortfolioService  PortfolioService
	AnalyticsService  AnalyticsService
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=InstrumentService
type InstrumentService interface {
	GetClient(ctx context.Context) (*investgo.Client, error)
	FindBy(client *investgo.Client, query string) ([]InstrumentShort, error)
	GetBondByUid(client *investgo.Client, uid string) (Bond, error)
	GetCurrencyBy(client *investgo.Client, figi string) (Currency, error)
	GetFutureBy(client *investgo.Client, figi string) (Future, error)
	GetShareCurrencyBy(client *investgo.Client, figi string) (ShareCurrencyByResponse, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=PortfolioService
type PortfolioService interface {
	GetClient(ctx context.Context) (*investgo.Client, error)
	GetAccounts(client *investgo.Client) (map[string]Account, error)
	GetPortfolio(client *investgo.Client, request PortfolioRequest) (Portfolio, error)
	GetOperations(client *investgo.Client, request OperationsRequest) ([]Operation, error)
	MakeSafeGetOperationsRequest(client *investgo.Client, request OperationsRequest) ([]Operation, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=AnalyticsService
type AnalyticsService interface {
	GetClient(ctx context.Context) (*investgo.Client, error)
	GetLastPriceInPersentageToNominal(client *investgo.Client, instrumentUid string) (LastPriceResponse, error)
	GetAllAssetUids(client *investgo.Client) (map[string]string, error)
	GetBondsActions(client *investgo.Client, instrumentUid string) (BondIdentIdentifiers, error)
}
