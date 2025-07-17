package service

import (
	"errors"
	"math"
	"time"
)

const (
	layout     = "2006-01-02"
	hoursInDay = 24
	daysInYear = 365
)

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

func (c *Client) CreateBondReport(reportPostions ReportPositions) (Report, error) {
	var resultReports Report
	for i := range reportPostions.CurrentPositions {
		position := reportPostions.CurrentPositions[i]
		switch position.Currency {
		case "rub":
			bondReport, err := c.createBondReport(position)
			if err != nil {
				return resultReports, errors.New("service: GetBondReport" + err.Error())
			}
			resultReports.BondsInRUB = append(resultReports.BondsInRUB, bondReport)
		case "cny":
			bondReport, err := c.createBondReport(position)
			if err != nil {
				return resultReports, errors.New("service: GetBondReport" + err.Error())
			}
			resultReports.BondsInCNY = append(resultReports.BondsInCNY, bondReport)
		default:
			continue
		}
	}
	return resultReports, nil
}

func (c *Client) createBondReport(position SharePosition) (BondReport, error) {
	var bondReport BondReport
	moexBuyData, err := c.MoexApi.GetSpecifications(position.Ticker, position.BuyDate)
	if err != nil {
		return bondReport, errors.New("service: createBondReport" + err.Error())
	}
	date := time.Now()
	moexLastPriceData, err := c.MoexApi.GetSpecifications(position.Ticker, date)
	if err != nil {
		return bondReport, errors.New("service: createBondReport" + err.Error())
	}

	bondReport = BondReport{
		Name:         position.Name,
		Ticker:       position.Ticker,
		BuyDate:      position.BuyDate.Format(layout),
		BuyPrice:     RoundFloat(position.BuyPrice, 2),
		CurrentPrice: RoundFloat(position.SellPrice, 2),
		Nominal:      position.Nominal,
	}
	lastPriceDataPath := moexLastPriceData.History.Data[0]
	buyPriceDataPath := moexBuyData.History.Data[0]

	maturityDate := lastPriceDataPath.MaturityDate
	if maturityDate != nil {
		bondReport.MaturityDate = *maturityDate
	}

	offerDate := lastPriceDataPath.OfferDate
	if offerDate != nil {
		bondReport.OfferDate = *offerDate
	}

	duration := lastPriceDataPath.Duration
	if duration != nil {
		bondReport.Duration = int64(*duration)
	}

	yieldToMaturity := lastPriceDataPath.YieldToMaturity
	if yieldToMaturity != nil {
		bondReport.YieldToMaturity = RoundFloat(*yieldToMaturity, 2)
	}

	yieldToOffer := lastPriceDataPath.YieldToOffer
	if yieldToOffer != nil {
		bondReport.YieldToOffer = RoundFloat(*yieldToOffer, 2)
	}

	yieldToMaturityOnPurchase := buyPriceDataPath.YieldToMaturity
	if yieldToMaturityOnPurchase != nil {
		bondReport.YieldToMaturityOnPurchase = RoundFloat(*yieldToMaturityOnPurchase, 2)
	}

	yieldToOfferOnPurchase := buyPriceDataPath.YieldToOffer
	if yieldToOfferOnPurchase != nil {
		bondReport.YieldToOfferOnPurchase = RoundFloat(*yieldToOfferOnPurchase, 2)
	}

	profitInPercentage, err := getProfit(position)
	if err != nil {
		return bondReport, errors.New("service: createBondReport" + err.Error())
	}
	bondReport.Profit = profitInPercentage

	annualizedReturn, err := getAnnualizedReturnInPercentage(position)
	if err != nil {
		return bondReport, errors.New("service: createBondReport" + err.Error())
	}
	bondReport.AnnualizedReturn = annualizedReturn

	return bondReport, nil
}

func getProfit(sharePosition SharePosition) (float64, error) {
	profitWithoutTax := getSecurityIncomeWithoutTax(sharePosition)
	totalTax := getTotalTaxFromPosition(sharePosition, profitWithoutTax)
	profit := getSecurityIncome(profitWithoutTax, totalTax)
	profitInPercentage, err := getProfitInPercentage(profit, sharePosition.BuyPrice, sharePosition.Quantity)
	if err != nil {
		return profitInPercentage, errors.New("service: GetProfit" + err.Error())
	}
	return profitInPercentage, nil
}

func getAnnualizedReturnInPercentage(p SharePosition) (float64, error) {
	var totalReturn float64
	profitWithoutTax := getSecurityIncomeWithoutTax(p)
	totalTax := getTotalTaxFromPosition(p, profitWithoutTax)
	profit := getSecurityIncome(profitWithoutTax, totalTax)
	buyDate := p.BuyDate
	// Костыль. Надо переписать когда-нибудь, так как для закрытх позиций данная функция работать не будет
	sellDate := time.Now()
	timeDurationInDays := sellDate.Sub(buyDate).Hours() / float64(hoursInDay)
	// Если покупка и продажа были совершены в один день, то берем минимум один день
	timeDurationInYears := math.Max(1, timeDurationInDays) / float64(daysInYear)
	if p.BuyPrice != 0 || p.Quantity != 0 {
		totalReturn = profit / (p.BuyPrice * p.Quantity)
	} else {
		return 0, errors.New("service: getAnnualizedReturn : divide by zero")
	}
	annualizedReturn := math.Pow((1+totalReturn), (1/timeDurationInYears)) - 1
	annualizedReturnInPercentage := annualizedReturn * 100
	annualizedReturnInPercentageRound := RoundFloat(annualizedReturnInPercentage, 2)

	return annualizedReturnInPercentageRound, nil

}
