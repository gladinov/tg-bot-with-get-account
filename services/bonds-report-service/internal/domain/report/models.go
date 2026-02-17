package report

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
