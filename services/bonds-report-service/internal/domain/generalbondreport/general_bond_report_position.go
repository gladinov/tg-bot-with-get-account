package generalbondreport

import (
	"bonds-report-service/internal/domain"
	report "bonds-report-service/internal/domain/report_position"
	"bonds-report-service/internal/utils"
	"time"

	"github.com/gladinov/e"
)

const (
	layout     = "2006-01-02"
	hoursInDay = 24
	daysInYear = 365
)

const (
	threeYearInHours = 26304 // Три года в часах
	baseTaxRate      = 0.13  // Налог с продажи ЦБ
)

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

func (p *GeneralBondReportPosition) CreateGeneralBondReportPosition(
	currentPositions []report.PositionByFIFO,
	totalAmount float64,
	moexBuyDateData domain.ValuesMoex,
	moexNowData domain.ValuesMoex,
	firstBondsBuyDate time.Time,
) (err error) {
	const op = "service.CreateGeneralBondReport"
	// Переменная ткущих позиций формата []report.PositionByFIFO
	// Создаем переменные

	var name string
	var duration int64
	var yieldToMaturity float64

	if len(currentPositions) == 0 {
		return ErrEmptyPositons
	}

	ticker := currentPositions[0].Ticker
	currency := currentPositions[0].Currency
	nominal := currentPositions[0].Nominal
	sellPrice := currentPositions[0].SellPrice
	replaced := currentPositions[0].Replaced

	sumOfPositionsRes := getSumOfPositions(currentPositions)

	positionsPrice := utils.RoundFloat(sumOfPositionsRes.SumOfPositions/sumOfPositionsRes.SumOfQuantity, 2)
	quantity := int64(sumOfPositionsRes.SumOfQuantity)
	profitOfAllPositions := utils.RoundFloat(sumOfPositionsRes.ProfitOfAllPositions, 2)
	profitInPercentage := utils.RoundFloat((profitOfAllPositions/sumOfPositionsRes.SumOfPositions)*100, 2)
	percentOfPortfolio := utils.RoundFloat((sumOfPositionsRes.SumOfPositions/totalAmount)*100, 2)

	if moexNowData.ShortName.IsHasValue() {
		name = moexNowData.ShortName.Value
	} else {
		name = currentPositions[0].Name
	}

	if moexNowData.Duration.IsHasValue() {
		duration = int64(moexNowData.Duration.Value)
	}

	if moexNowData.YieldToOffer.IsHasValue() {
		yieldToMaturity = moexNowData.YieldToOffer.Value
	} else {
		if moexNowData.YieldToMaturity.IsHasValue() {
			yieldToMaturity = moexNowData.YieldToMaturity.Value
		}
	}

	maturityDate, err := getFirstMuturityDate(moexNowData.BuybackDate, moexNowData.OfferDate, moexNowData.MaturityDate)
	if err != nil {
		return e.WrapIfErr("failed to get first maturity date", err)
	}

	p.Ticker = ticker
	p.Currencies = currency
	p.Nominal = nominal
	p.CurrentPrice = sellPrice
	p.Replaced = replaced
	p.PositionPrice = positionsPrice
	p.BuyDate = firstBondsBuyDate
	p.Quantity = quantity
	p.Profit = profitOfAllPositions
	p.ProfitInPercentage = profitInPercentage
	p.PercentOfPortfolio = percentOfPortfolio
	p.Name = name
	p.Duration = duration
	p.YieldToMaturity = yieldToMaturity
	p.MaturityDate = maturityDate

	if moexBuyDateData.YieldToOffer.IsHasValue() {
		yieldToOfferOnPurchase := moexBuyDateData.YieldToOffer.Value
		p.YieldToMaturityOnPurchase = yieldToOfferOnPurchase
	} else {
		if moexBuyDateData.YieldToMaturity.IsHasValue() {
			yieldToMaturityOnPurchase := moexBuyDateData.YieldToMaturity.Value
			p.YieldToMaturityOnPurchase = yieldToMaturityOnPurchase
		}
	}

	return nil
}
