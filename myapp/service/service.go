package service

import (
	"time"

	"main.go/clients/cbr"
	"main.go/clients/moex"
	"main.go/clients/tinkoffApi"
)

type Client struct {
	Tinkoffapi *tinkoffApi.Client
	MoexApi    *moex.Client
	CbrApi     *cbr.Client
}

func New(tinkoffApiClient *tinkoffApi.Client, moexClient *moex.Client, CbrClient *cbr.Client) *Client {
	return &Client{
		Tinkoffapi: tinkoffApiClient,
		MoexApi:    moexClient,
		CbrApi:     CbrClient,
	}
}

type Operation struct {
	BrokerAccountId   string
	Currency          string
	Operation_Id      string
	ParentOperationId string
	Name              string
	Date              time.Time // Время в UTC
	Type              int64
	Description       string
	InstrumentUid     string
	Figi              string
	InstrumentType    string
	InstrumentKind    string
	PositionUid       string
	Payment           float64
	Price             float64
	Commission        float64
	Yield             float64
	YieldRelative     float64
	AccruedInt        float64
	QuantityDone      float64
	AssetUid          string
}

type Portfolio struct {
	PortfolioPositions []PortfolioPosition
	BondPositions      []Bond
}

type PortfolioPosition struct {
	AccountId                string
	Figi                     string
	InstrumentType           string
	Currency                 string
	Quantity                 float64
	AveragePositionPrice     float64
	ExpectedYield            float64
	CurrentNkd               float64
	CurrentPrice             float64
	AveragePositionPriceFifo float64
	Blocked                  bool
	BlockedLots              float64
	PositionUid              string
	InstrumentUid            string
	AssetUid                 string
	VarMargin                float64
	ExpectedYieldFifo        float64
	DailyYield               float64
}

type Bond struct {
	Identifiers              Identifiers
	Name                     string  // GetBondsActionsFromPortfolio
	InstrumentType           string  // T_Api_Getportfolio
	Currency                 string  // T_Api_Getportfolio
	Quantity                 float64 // T_Api_Getportfolio
	AveragePositionPrice     float64 // T_Api_Getportfolio
	ExpectedYield            float64 // T_Api_Getportfolio
	CurrentNkd               float64 // T_Api_Getportfolio
	CurrentPrice             float64 // T_Api_Getportfolio
	AveragePositionPriceFifo float64 // T_Api_Getportfolio
	Blocked                  bool    // T_Api_Getportfolio
	ExpectedYieldFifo        float64 // T_Api_Getportfolio
	DailyYield               float64 // T_Api_Getportfolio
}

type Identifiers struct {
	Ticker        string // GetBondsActionsFromPortfolio
	ClassCode     string // GetBondsActionsFromPortfolio
	Figi          string // T_Api_Getportfolio
	InstrumentUid string // T_Api_Getportfolio
	PositionUid   string // T_Api_Getportfolio
	AssetUid      string // GetBondsActionsFromPortfolio
}
