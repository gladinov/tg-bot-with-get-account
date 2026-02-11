package service_models

import (
	"time"
)

const (
	RubBonds      = "bondsInRub"
	ReplacedBonds = "replacedBonds"
	EuroBonds     = "euroBonds"
)

type ReportPositions struct {
	Quantity         float64
	CurrentPositions []PositionByFIFO
}

type PositionByFIFO struct {
	Name                  string
	Replaced              bool
	CurrencyIfReplaced    string
	BuyDate               time.Time
	SellDate              time.Time
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

type PortfolioPosition struct {
	InstrumentType string
	AssetUid       string
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

type MediaGroup struct {
	Reports []*ImageData `json:"reports"`
}

func NewMediaGroup() *MediaGroup {
	return &MediaGroup{
		Reports: make([]*ImageData, 0),
	}
}

type ImageData struct {
	Name    string `json:"name"`
	Data    []byte `json:"data"`
	Caption string `json:"caption"`
}

func NewImageData() *ImageData {
	return &ImageData{}
}

type AccountListResponce struct {
	Accounts string `json:"accounts,omitempty"`
}

type BondReportsRequest struct {
	ChatID int `json:"chatID,omitempty"`
}

type UsdResponce struct {
	Usd float64 `json:"usd,omitempty"`
}

type BondReportsResponce struct {
	Media [][]*MediaGroup `json:"media"`
}

type PortfolioStructureForEachAccountResponce struct {
	PortfolioStructures []string `json:"potftfolio"`
}

type UnionPortfolioStructureResponce struct {
	Report string `json:"report"`
}

type UnionPortfolioStructureWithSberResponce struct {
	Report string `json:"report"`
}
