package service

import (
	"errors"
	"math"
	"time"

	"main.go/lib/e"
	"main.go/service/service_models"
)

const (
	layout     = "2006-01-02"
	hoursInDay = 24
	daysInYear = 365
)

func (c *Client) CreateBondReport(reportPostions service_models.ReportPositions) (service_models.Report, error) {
	var resultReports service_models.Report
	for i := range reportPostions.CurrentPositions {
		position := reportPostions.CurrentPositions[i]
		switch position.Currency {
		case "rub":
			bondReport, err := c.createBondReportByCurrency(position)
			if err != nil {
				return resultReports, errors.New("service: GetBondReport" + err.Error())
			}
			resultReports.BondsInRUB = append(resultReports.BondsInRUB, bondReport)
		case "cny":
			bondReport, err := c.createBondReportByCurrency(position)
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

func (c *Client) CreateGeneralBondReport(resultBondPosition *service_models.ReportPositions, totalAmount float64) (_ service_models.GeneralBondReportPosition, err error) {
	defer func() { err = e.WrapIfErr("can't create general bond report", err) }()
	currentPostions := resultBondPosition.CurrentPositions
	var BondReporPosition service_models.GeneralBondReportPosition

	if len(currentPostions) == 0 {
		return BondReporPosition, errors.New("no input positions")
	}
	BondReporPosition.Ticker = currentPostions[0].Ticker
	BondReporPosition.Currencies = currentPostions[0].Currency
	BondReporPosition.Nominal = currentPostions[0].Nominal
	BondReporPosition.CurrentPrice = currentPostions[0].SellPrice
	BondReporPosition.Replaced = currentPostions[0].Replaced

	var sumOfPosition float64
	buyDate := currentPostions[0].BuyDate
	var sumOfQuantity float64
	var profit float64

	for _, position := range currentPostions {
		sumOfPosition += position.BuyPrice * position.Quantity
		compareDate := position.BuyDate
		if compareDate.Compare(buyDate) == -1 {
			buyDate = compareDate
		}
		sumOfQuantity += position.Quantity

		profitWithoutTax := getSecurityIncomeWithoutTax(position)
		totalTax := getTotalTaxFromPosition(position, profitWithoutTax)
		profit += getSecurityIncome(profitWithoutTax, totalTax)
	}
	BondReporPosition.PositionPrice = RoundFloat(sumOfPosition/sumOfQuantity, 2)
	BondReporPosition.BuyDate = buyDate
	BondReporPosition.Quantity = int64(sumOfQuantity)
	BondReporPosition.Profit = RoundFloat(profit, 2)
	BondReporPosition.ProfitInPercentage = RoundFloat((profit/sumOfPosition)*100, 2)
	BondReporPosition.PercentOfPortfolio = RoundFloat((sumOfPosition/totalAmount)*100, 2)

	moexBuyDateData, err := c.MoexApi.GetSpecifications(BondReporPosition.Ticker, buyDate)
	if err != nil {
		return BondReporPosition, err
	}
	date := time.Now()
	moexNowData, err := c.MoexApi.GetSpecifications(BondReporPosition.Ticker, date)
	if err != nil {
		return BondReporPosition, err
	}

	lastPriceDataPath := moexNowData.History.Data[0]

	if lastPriceDataPath.ShortName != nil {
		BondReporPosition.Name = *lastPriceDataPath.ShortName

	} else {
		BondReporPosition.Name = currentPostions[0].Name
	}

	duration := lastPriceDataPath.Duration
	if duration != nil {
		BondReporPosition.Duration = int64(*duration)
	}

	yieldToOffer := lastPriceDataPath.YieldToOffer
	yieldToMaturity := lastPriceDataPath.YieldToMaturity
	if yieldToOffer != nil {
		BondReporPosition.YieldToMaturity = *yieldToOffer
	} else {
		if yieldToMaturity != nil {
			BondReporPosition.YieldToMaturity = *yieldToMaturity
		}
	}

	maturityDate := lastPriceDataPath.MaturityDate
	offerDate := lastPriceDataPath.OfferDate
	buyBackDate := lastPriceDataPath.BuybackDate

	switch {
	case offerDate != nil && buyBackDate != nil:
		offerDateConv, err := time.Parse(layout, *offerDate)
		if err != nil {
			return BondReporPosition, err
		}
		buyBackDateConv, err := time.Parse(layout, *buyBackDate)
		if err != nil {
			return BondReporPosition, err
		}
		if offerDateConv.Compare(buyBackDateConv) == -1 {
			BondReporPosition.MaturityDate = offerDateConv
		} else {
			BondReporPosition.MaturityDate = buyBackDateConv
		}
	case offerDate != nil:
		offerDateConv, err := time.Parse(layout, *offerDate)
		if err != nil {
			return BondReporPosition, err
		}
		BondReporPosition.MaturityDate = offerDateConv
	case buyBackDate != nil:
		buyBackDateConv, err := time.Parse(layout, *buyBackDate)
		if err != nil {
			return BondReporPosition, err
		}
		BondReporPosition.MaturityDate = buyBackDateConv
	case maturityDate != nil:
		maturityDateConv, err := time.Parse(layout, *maturityDate)
		if err != nil {
			return BondReporPosition, err
		}
		BondReporPosition.MaturityDate = maturityDateConv
	}

	buyPriceDataPath := moexBuyDateData.History.Data[0]

	yieldToOfferOnPurchase := buyPriceDataPath.YieldToOffer
	yieldToMaturityOnPurchase := buyPriceDataPath.YieldToMaturity
	if yieldToOfferOnPurchase != nil {
		BondReporPosition.YieldToMaturityOnPurchase = *yieldToOfferOnPurchase
	} else {
		if yieldToMaturityOnPurchase != nil {
			BondReporPosition.YieldToMaturityOnPurchase = *yieldToMaturityOnPurchase
		}
	}

	return BondReporPosition, nil
}

func (c *Client) createBondReportByCurrency(position service_models.PositionByFIFO) (service_models.BondReport, error) {
	var bondReport service_models.BondReport
	moexBuyDateData, err := c.MoexApi.GetSpecifications(position.Ticker, position.BuyDate)
	if err != nil {
		return bondReport, errors.New("service: createBondReport" + err.Error())
	}
	date := time.Now()
	moexNowData, err := c.MoexApi.GetSpecifications(position.Ticker, date)
	if err != nil {
		return bondReport, errors.New("service: createBondReport" + err.Error())
	}

	bondReport = service_models.BondReport{
		Name:         position.Name,
		Ticker:       position.Ticker,
		BuyDate:      position.BuyDate.Format(layout),
		BuyPrice:     RoundFloat(position.BuyPrice, 2),
		CurrentPrice: RoundFloat(position.SellPrice, 2),
		Nominal:      position.Nominal,
	}
	lastPriceDataPath := moexNowData.History.Data[0]
	buyPriceDataPath := moexBuyDateData.History.Data[0]

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

func getProfit(sharePosition service_models.PositionByFIFO) (float64, error) {
	profitWithoutTax := getSecurityIncomeWithoutTax(sharePosition)
	totalTax := getTotalTaxFromPosition(sharePosition, profitWithoutTax)
	profit := getSecurityIncome(profitWithoutTax, totalTax)
	profitInPercentage, err := getProfitInPercentage(profit, sharePosition.BuyPrice, sharePosition.Quantity)
	if err != nil {
		return profitInPercentage, errors.New("service: GetProfit" + err.Error())
	}
	return profitInPercentage, nil
}

func getAnnualizedReturnInPercentage(p service_models.PositionByFIFO) (float64, error) {
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
