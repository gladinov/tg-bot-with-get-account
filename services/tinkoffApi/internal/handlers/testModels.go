package handlers

import (
	"time"
	"tinkoffApi/internal/service"
)

var accountHappyPath = `{
    "2007907898": {
        "id": "2007907898",
        "type": "ACCOUNT_TYPE_TINKOFF",
        "name": "Брокерский счёт",
        "status": 2,
        "openedDate": "2019-05-21T00:00:00Z",
        "closedDate": "1970-01-01T00:00:00Z",
        "accessLevel": 2
    },
    "2012259491": {
        "id": "2012259491",
        "type": "ACCOUNT_TYPE_TINKOFF_IIS",
        "name": "ИИС",
        "status": 3,
        "openedDate": "2019-11-12T00:00:00Z",
        "closedDate": "2023-03-31T00:00:00Z",
        "accessLevel": 2
    },
    "2016119489": {
        "id": "2016119489",
        "type": "ACCOUNT_TYPE_TINKOFF",
        "name": "Брокерский счет 1",
        "status": 2,
        "openedDate": "2023-02-06T00:00:00Z",
        "closedDate": "1970-01-01T00:00:00Z",
        "accessLevel": 2
    }
}`

var portffolioHappyPath = `{
    "positions": [
        {
            "figi": "TCS00A107UB5",
            "instrumentType": "bond",
            "quantity": {
                "units": 69
            },
            "averagePositionPrice": {
                "currency": "rub",
                "units": 1005,
                "nano": 800000000
            },
            "expectedYield": {
                "units": -469,
                "nano": -200000000
            },
            "currentNkd": {
                "currency": "rub",
                "units": 37,
                "nano": 810000000
            },
            "currentPrice": {
                "currency": "rub",
                "units": 999
            },
            "averagePositionPriceFifo": {
                "currency": "rub",
                "units": 1005,
                "nano": 800000000
            },
            "blockedLots": {},
            "positionUid": "98ce5918-9479-4af8-acb1-62250ee65744",
            "instrumentUid": "01e0c046-dfbe-4385-8a9f-3ccbfdeb3c4d",
            "varMargin": {},
            "expectedYieldFifo": {
                "units": -469,
                "nano": -200000000
            },
            "dailyYield": {
                "currency": "rub",
                "units": -13,
                "nano": -800000000
            },
            "ticker": "RU000A107UB5"
        }
    ],
    "totalAmount": {
        "currency": "rub",
        "units": 1627727,
        "nano": 490000000
    }
}`

var TestPortfolio = service.Portfolio{
	Positions: PortfolioPositions,
	TotalAmount: service.MoneyValue{
		Currency: "rub",
		Units:    1627727,
		Nano:     490000000,
	},
}

var PortfolioPositions = []service.PortfolioPositions{
	{
		Figi:           "TCS00A107UB5",
		InstrumentType: "bond",
		Quantity: service.Quotation{
			Units: 69,
		},
		AveragePositionPrice: service.MoneyValue{
			Currency: "rub",
			Units:    1005,
			Nano:     800000000,
		},
		ExpectedYield: service.Quotation{
			Units: -469,
			Nano:  -200000000,
		},
		CurrentNkd: service.MoneyValue{
			Currency: "rub",
			Units:    37,
			Nano:     810000000,
		},
		CurrentPrice: service.MoneyValue{
			Currency: "rub",
			Units:    999,
		},
		AveragePositionPriceFifo: service.MoneyValue{
			Currency: "rub",
			Units:    1005,
			Nano:     800000000,
		},
		BlockedLots:   service.Quotation{},
		PositionUid:   "98ce5918-9479-4af8-acb1-62250ee65744",
		InstrumentUid: "01e0c046-dfbe-4385-8a9f-3ccbfdeb3c4d",
		VarMargin:     service.MoneyValue{},
		ExpectedYieldFifo: service.Quotation{
			Units: -469,
			Nano:  -200000000,
		},
		DailyYield: service.MoneyValue{
			Currency: "rub",
			Units:    -13,
			Nano:     -800000000,
		},
		Ticker: "RU000A107UB5",
	},
}

var portfolioRequest = service.PortfolioRequest{
	AccountID:     "2007907898",
	AccountStatus: 2,
}

var emptyAccounts = make(map[string]service.Account)

var accountsServiceHappyPath = map[string]service.Account{
	"2007907898": {
		Id:          "2007907898",
		Type:        "ACCOUNT_TYPE_TINKOFF",
		Name:        "Брокерский счёт",
		Status:      2,
		OpenedDate:  time.Date(2019, time.May, 21, 0, 0, 0, 0, time.UTC),
		ClosedDate:  time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
		AccessLevel: 2,
	},
	"2012259491": {
		Id:          "2012259491",
		Type:        "ACCOUNT_TYPE_TINKOFF_IIS",
		Name:        "ИИС",
		Status:      3,
		OpenedDate:  time.Date(2019, time.November, 12, 0, 0, 0, 0, time.UTC),
		ClosedDate:  time.Date(2023, time.March, 31, 0, 0, 0, 0, time.UTC),
		AccessLevel: 2,
	},
	"2016119489": {
		Id:          "2016119489",
		Type:        "ACCOUNT_TYPE_TINKOFF",
		Name:        "Брокерский счет 1",
		Status:      2,
		OpenedDate:  time.Date(2023, time.February, 6, 0, 0, 0, 0, time.UTC),
		ClosedDate:  time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
		AccessLevel: 2,
	},
}

var operationRequest = service.OperationsRequest{
	AccountID: "2007907898",
	Date:      time.Date(2025, time.October, 21, 0, 0, 0, 0, time.UTC),
}

var operationBadRequest = service.OperationsRequest{
	AccountID: "",
	Date:      time.Date(2025, time.October, 21, 0, 0, 0, 0, time.UTC),
}

var happyPathOperationResponseInBytes = `[
    {
        "brokerAccountId": "2007907898",
        "currency": "rub",
        "operationId": "28596",
        "date": "2025-11-01T11:02:55.028Z",
        "type": 27,
        "description": "Списание вариационной маржи",
        "instrumentKind": "INSTRUMENT_TYPE_UNSPECIFIED",
        "payment": {
            "currency": "rub",
            "units": -523
        },
        "price": {
            "currency": "rub"
        },
        "commission": {},
        "yield": {},
        "yieldRelative": {},
        "accruedInt": {}
    }]`

var happyPathOperationResponse = []service.Operation{
	{
		BrokerAccountId: "2007907898",
		Currency:        "rub",
		Operation_Id:    "28596",
		Date:            time.Date(2025, 11, 1, 11, 2, 55, 28000000, time.UTC),
		Type:            27,
		Description:     "Списание вариационной маржи",
		InstrumentKind:  "INSTRUMENT_TYPE_UNSPECIFIED",
		Payment: service.MoneyValue{
			Currency: "rub",
			Units:    -523,
		},
		Price: service.MoneyValue{
			Currency: "rub",
		},
		Commission:    service.MoneyValue{},
		Yield:         service.MoneyValue{},
		YieldRelative: service.Quotation{},
		AccruedInt:    service.MoneyValue{},
	},
}

var happyPathAllAssetsUid = map[string]string{"000029ae-00a2-441c-a9bf-9bfd2988e706": "f1c6effc-4595-48eb-8523-d4a9a1d1dd4f"}

var happyPathAllAssetsUidInBytes = `{"000029ae-00a2-441c-a9bf-9bfd2988e706": "f1c6effc-4595-48eb-8523-d4a9a1d1dd4f"}`

var happyPathBondsActionsInBytes = `{
    "ticker": "RU000A101XD8",
    "classCode": "TQCB",
    "name": "МаксимаТелеком выпуск 1",
    "nominal": {
        "currency": "rub",
        "units": 1000
    },
    "nominalCurrency": "rub"
}`

var happyPathBondsActions = service.BondIdentIdentifiers{
	Ticker:          "RU000A101XD8",
	ClassCode:       "TQCB",
	Name:            "МаксимаТелеком выпуск 1",
	Nominal:         service.MoneyValue{Currency: "rub", Units: 1000},
	NominalCurrency: "rub",
}

var happyPathLastPrice = service.LastPriceResponse{
	LastPrice: service.Quotation{
		Units: 101, Nano: 290000000,
	},
}

var happyPathLastPriceInBytes = `{
    "lastPrice": {
        "units": 101,
        "nano": 290000000
    }
}`

var happpyPathFindByInBytes = `[
    {
        "instrumentType": "share",
        "uid": "76721c1c-52a9-4b45-987e-d075f651f1b1",
        "figi": "BBG000RJL816"
    }
]`

var happpyPathFindBy = []service.InstrumentShort{{
	InstrumentType: "share",
	Uid:            "76721c1c-52a9-4b45-987e-d075f651f1b1",
	Figi:           "BBG000RJL816",
}}

var happpyPathBondByInBytes = `{
    "aciValue": {
        "currency": "rub",
        "units": 66,
        "nano": 900000000
    },
    "currency": "rub",
    "nominal": {
        "currency": "rub",
        "units": 1000
    }
}`

var happpyPathBondBy = service.Bond{
	AciValue: service.MoneyValue{
		Currency: "rub",
		Units:    66,
		Nano:     900000000,
	},
	Currency: "rub",
	Nominal: service.MoneyValue{
		Currency: "rub",
		Units:    1000,
	},
}

var happpyPathCurrencyBy = service.Currency{
	Isin: "cny",
}

var happpyPathCurrencyByInBytes = `{
    "isin": "cny"
}`

var happyPathShareCurrency = service.ShareCurrencyByResponse{
	Currency: "rub",
}

var happyPathShareCurrencyInBytes = `{
    "currency": "rub"
}`

var happyPathFutureBy = service.Future{
	Name:                    "CNY-3.23 Курс Юань - Рубль",
	MinPriceIncrement:       service.Quotation{},
	MinPriceIncrementAmount: service.Quotation{},
	AssetType:               "TYPE_CURRENCY",
	BasicAssetPositionUid:   "176c3dbf-b346-48a6-b20c-daa9d028f031",
}

var happyPathFutureByInBytes = `{
    "name": "CNY-3.23 Курс Юань - Рубль",
    "minPriceIncrement": {},
    "minPriceIncrementAmount": {},
    "assetType": "TYPE_CURRENCY",
    "basicAssetPositionUid": "176c3dbf-b346-48a6-b20c-daa9d028f031"
}`
