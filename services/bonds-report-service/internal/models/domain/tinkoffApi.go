package domain

import (
	"math"
	"time"
)

const (
	RubBonds      = "bondsInRub"
	ReplacedBonds = "replacedBonds"
	EuroBonds     = "euroBonds"
)

type BondIdentIdentifiers struct {
	Ticker          string
	ClassCode       string
	Name            string
	Nominal         MoneyValue
	NominalCurrency string
	Replaced        bool
}

type LastPrice struct {
	LastPrice Quotation
}

type MoneyValue struct {
	Currency string
	Units    int64
	Nano     int32
}

func (x *MoneyValue) GetCurrency() string {
	if x != nil {
		return x.Currency
	}
	return ""
}

func (x *MoneyValue) GetUnits() int64 {
	if x != nil {
		return x.Units
	}
	return 0
}

func (x *MoneyValue) GetNano() int32 {
	if x != nil {
		return x.Nano
	}
	return 0
}

func (mv *MoneyValue) ToFloat() float64 {
	num := float64(mv.Units) + float64(mv.Nano)*math.Pow10(-9)
	num = num * math.Pow10(9)
	num = math.Round(num)
	num = num / math.Pow10(9)
	return num
}

type Quotation struct {
	Units int64
	Nano  int32
}

func (x *Quotation) GetUnits() int64 {
	if x != nil {
		return x.Units
	}
	return 0
}

func (x *Quotation) GetNano() int32 {
	if x != nil {
		return x.Nano
	}
	return 0
}

func (q *Quotation) ToFloat() float64 {
	num := float64(q.Units) + float64(q.Nano)*math.Pow10(-9)
	num = num * math.Pow10(9)
	num = math.Round(num)
	num = num / math.Pow10(9)
	return num
}

type Future struct {
	Name                    string
	MinPriceIncrement       Quotation
	MinPriceIncrementAmount Quotation
	AssetType               string
	BasicAssetPositionUid   string
}

type Portfolio struct {
	Positions   []PortfolioPosition
	TotalAmount MoneyValue
}

type PortfolioPosition struct {
	Figi                     string
	InstrumentType           string
	Quantity                 Quotation
	AveragePositionPrice     MoneyValue
	ExpectedYield            Quotation
	CurrentNkd               MoneyValue
	CurrentPrice             MoneyValue
	AveragePositionPriceFifo MoneyValue
	Blocked                  bool
	BlockedLots              Quotation
	PositionUid              string
	InstrumentUid            string
	VarMargin                MoneyValue
	ExpectedYieldFifo        Quotation
	DailyYield               MoneyValue
	Ticker                   string
}

type PortfolioPositionsWithAssetUid struct {
	InstrumentType string
	AssetUid       string
}

type Bond struct {
	AciValue MoneyValue
	Currency string
	Nominal  MoneyValue
}

type Currency struct {
	Isin string
}

type InstrumentShortList []InstrumentShort

func (list InstrumentShortList) ValidateAndGetFirstShare() (InstrumentShort, error) {
	if len(list) == 0 {
		return InstrumentShort{}, ErrEmptyInstrumentShortResponce
	}

	inst := list[0]

	if inst.InstrumentType != share {
		return InstrumentShort{}, ErrInstrumentNotShare
	}

	if inst.Figi == "" {
		return InstrumentShort{}, ErrEmptyFigi
	}

	return inst, nil
}

type InstrumentShort struct {
	InstrumentType string
	Uid            string
	Figi           string
}

type Operation struct {
	BrokerAccountID   string
	Currency          string
	OperationID       string
	ParentOperationID string
	Name              string
	Date              time.Time
	Type              int64
	Description       string
	InstrumentUid     string
	Figi              string
	InstrumentType    string
	InstrumentKind    string
	PositionUid       string
	Payment           MoneyValue
	Price             MoneyValue
	Commission        MoneyValue
	Yield             MoneyValue
	YieldRelative     Quotation
	AccruedInt        MoneyValue
	QuantityDone      int64
	AssetUid          string
}

type OperationWithoutCustomTypes struct {
	BrokerAccountID   string
	Currency          string
	OperationID       string
	ParentOperationID string
	Name              string
	Date              time.Time
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

type Account struct {
	ID          string
	Type        string
	Name        string
	Status      int64
	OpenedDate  time.Time
	ClosedDate  time.Time
	AccessLevel int64
}

func (a Account) ValidateForPortfolio() error {
	switch a.Status {
	case 0:
		return ErrUnspecifiedAccount
	case 1:
		return ErrNewNotOpenYetAccount
	case 3:
		return ErrCloseAccount
	}

	if a.AccessLevel == 3 {
		return ErrNoAcces
	}

	if a.ID == "" {
		return ErrEmptyAccountIdInRequest
	}

	return nil
}

type ShareCurrency struct {
	Currency string
}
