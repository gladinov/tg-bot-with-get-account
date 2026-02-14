package domain

import "time"

type MediaGroup struct {
	Reports []*ImageData
}

func NewMediaGroup() *MediaGroup {
	return &MediaGroup{
		Reports: make([]*ImageData, 0),
	}
}

type ImageData struct {
	Name    string
	Data    []byte
	Caption string
}

func NewImageData() *ImageData {
	return &ImageData{}
}

type AccountListResponce struct {
	Accounts string
}

type BondReportsResponce struct {
	Media [][]*MediaGroup
}

type PortfolioStructureForEachAccountResponce struct {
	PortfolioStructures []string
}

type UnionPortfolioStructureResponce struct {
	Report string
}

type UnionPortfolioStructureWithSberResponce struct {
	Report string
}

type UsdResponce struct {
	Usd float64
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

type GeneralBondReports struct {
	RubBondsReport      map[TickerTimeKey]GeneralBondReportPosition
	EuroBondsReport     map[TickerTimeKey]GeneralBondReportPosition
	ReplacedBondsReport map[TickerTimeKey]GeneralBondReportPosition
}

type GeneralBondReportPosition struct {
	Name                      string
	Ticker                    string
	Replaced                  bool
	Currencies                string
	Quantity                  int64
	PercentOfPortfolio        float64
	MaturityDate              time.Time // дата погашения или выкупа или опциона
	Duration                  int64
	BuyDate                   time.Time
	PositionPrice             float64 // Средняя цена позиции
	YieldToMaturityOnPurchase float64 // Доходность при покупке до даты погашения или выкупа или опциона
	YieldToMaturity           float64 // Текущая доходность к погашению или выкупу или опциону
	CurrentPrice              float64
	Nominal                   float64
	Profit                    float64 // Результат инвестиции
	ProfitInPercentage        float64
}

type TickerTimeKey struct {
	Ticker string
	Time   time.Time
}

type PortfolioByTypeAndCurrency struct {
	AllAssets        float64
	BondsAssets      BondsAssets
	SharesAssets     SharesAssets
	EtfsAssets       EtfsAssets
	FuturesAssets    FuturesAssets
	CurrenciesAssets CurrenciesAssets
}

type BondsAssets struct {
	SumOfAssets      float64
	AssetsByCurrency map[string]*AssetByParam
}

type AssetByParam struct {
	SumOfAssets float64
}

func NewAssetsByParam() *AssetByParam {
	return &AssetByParam{}
}

type SharesAssets struct {
	SumOfAssets      float64
	AssetsByCurrency map[string]*AssetByParam
}

type EtfsAssets struct {
	SumOfAssets      float64
	AssetsByCurrency map[string]*AssetByParam
}

type FuturesAssets struct {
	SumOfAssets  float64
	AssetsByType AssetsByType
}

type AssetsByType struct {
	Commodity FuturesType
	Currency  FuturesType
	Security  FuturesType
	Index     FuturesType
}

type FuturesType struct {
	SumOfAssets      float64
	AssetsByCurrency map[string]*AssetByParam
}

type CurrenciesAssets struct {
	SumOfAssets      float64
	AssetsByCurrency map[string]*AssetByParam
}

func NewPortfolioByTypeAndCurrency() *PortfolioByTypeAndCurrency {
	return &PortfolioByTypeAndCurrency{
		AllAssets: 0,
		BondsAssets: BondsAssets{
			AssetsByCurrency: make(map[string]*AssetByParam),
		},
		SharesAssets: SharesAssets{
			AssetsByCurrency: make(map[string]*AssetByParam),
		},
		EtfsAssets: EtfsAssets{
			AssetsByCurrency: make(map[string]*AssetByParam),
		},
		FuturesAssets: FuturesAssets{
			AssetsByType: AssetsByType{
				Commodity: FuturesType{
					AssetsByCurrency: make(map[string]*AssetByParam),
				},
				Currency: FuturesType{
					AssetsByCurrency: make(map[string]*AssetByParam),
				},
				Security: FuturesType{
					AssetsByCurrency: make(map[string]*AssetByParam),
				},
				Index: FuturesType{
					AssetsByCurrency: make(map[string]*AssetByParam),
				},
			},
		},
		CurrenciesAssets: CurrenciesAssets{
			AssetsByCurrency: make(map[string]*AssetByParam),
		},
	}
}

type OperationsRequest struct {
	AccountID string
	FromDate  time.Time
}

func NewOperationsRequest(accountID string, fromDate time.Time) OperationsRequest {
	return OperationsRequest{
		AccountID: accountID,
		FromDate:  fromDate,
	}
}

func (r OperationsRequest) Validate(now time.Time) error {
	if r.AccountID == "" {
		return ErrEmptyAccountID
	}
	if r.FromDate.After(now) {
		return ErrInvalidFromDate
	}
	return nil
}

type ReportLine struct {
	Operation  []OperationWithoutCustomTypes
	Bond       BondIdentIdentifiers
	LastPrice  LastPrice
	Vunit_rate Rate
}

func NewReportLine(op []OperationWithoutCustomTypes, bond BondIdentIdentifiers, price LastPrice, vunit_rate Rate) ReportLine {
	return ReportLine{
		Operation:  op,
		Bond:       bond,
		LastPrice:  price,
		Vunit_rate: vunit_rate,
	}
}

type Rate struct {
	IsoCurrencyName string
	Vunit_Rate      NullFloat64
}
