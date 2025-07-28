package service_models

import (
	"errors"
	"time"
)

var ErrEmptyUids = errors.New("no uids")
var ErrNoCurrency = errors.New("no currency")
var ErrNoOpperations = errors.New("no operations")

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

type ReportPositions struct {
	Quantity         float64
	CurrentPositions []SharePosition
	ClosePositions   []SharePosition
}

type SharePosition struct {
	Name                  string
	BuyDate               time.Time
	SellDate              time.Time
	BuyOperationID        string
	SellOperationID       string
	Quantity              float64
	Figi                  string
	InstrumentType        string
	InstrumentUid         string
	Ticker                string
	ClassCode             string
	Nominal               float64
	BuyPrice              float64
	SellPrice             float64 // Для открытых позиций.Текущая цена с биржи
	BuyPayment            float64
	SellPayment           float64
	Currency              string
	BuyAccruedInt         float64 // НКД при покупке
	SellAccruedInt        float64
	PartialEarlyRepayment float64 // Частичное досрочное гашение
	TotalCoupon           float64
	TotalDividend         float64
	TotalComission        float64
	PaidTax               float64 // Фактически уплаченный налог(Часть налога будет уплачена в конце года, либо при выводе средств)
	TotalTax              float64 // Налог рассчитанный
	PositionProfit        float64 // С учетом рассчитанных налогов(TotalTax)
	ProfitInPercentage    float64 // В процентах строковая переменная
}

type Report struct {
	BondsInRUB []BondReport
	BondsInCNY []BondReport
}

type BondReport struct {
	Name                      string
	Ticker                    string
	MaturityDate              string // Дата погашения
	OfferDate                 string
	Duration                  int64
	BuyDate                   string
	BuyPrice                  float64
	YieldToMaturityOnPurchase float64 // Доходность к погашению при покупке
	YieldToOfferOnPurchase    float64 // Доходность к оферте при покупке
	YieldToMaturity           float64 // Текущая доходность к погашению
	YieldToOffer              float64 // Текущая доходность к оферте
	CurrentPrice              float64
	Nominal                   float64
	Profit                    float64 // Результат инвестиции
	AnnualizedReturn          float64 // Годовая доходность
}

type PortfolioPosition struct {
	InstrumentType string
	AssetUid       string
}

type Currency struct {
	Date      time.Time
	NumCode   string
	CharCode  string
	Nominal   int
	Name      string
	Value     float64
	VunitRate float64
}

type Currencies struct {
	CurrenciesMap map[string]Currency
}
