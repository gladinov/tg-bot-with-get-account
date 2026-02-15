package bondreport

import (
	"bonds-report-service/internal/models/domain"
	bondreport "bonds-report-service/internal/models/domain/bondReport"
	"bonds-report-service/internal/models/domain/report"
	"bonds-report-service/internal/utils"
	"bonds-report-service/internal/utils/logging"
	"context"
	"errors"
	"log/slog"
	"math"
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

func (s *BondReporter) CreateBondReport(ctx context.Context, reportPostions report.ReportPositions) (_ domain.Report, err error) {
	const op = "service.CreateBondReport"

	var resultReports domain.Report
	for i := range reportPostions.CurrentPositions {
		position := reportPostions.CurrentPositions[i]
		switch position.Currency {
		case "rub":
			bondReport, err := s.createBondReportByCurrency(ctx, position)
			if err != nil {
				return resultReports, errors.New("service: GetBondReport" + err.Error())
			}
			resultReports.BondsInRUB = append(resultReports.BondsInRUB, bondReport)
		case "cny":
			bondReport, err := s.createBondReportByCurrency(ctx, position)
			if err != nil {
				return resultReports, errors.New("service: GetBondReport" + err.Error())
			}
			resultReports.BondsInCNY = append(resultReports.BondsInCNY, bondReport)
		default:
			s.logger.WarnContext(ctx, "unkown currency", slog.String("op", op), slog.String("currency", position.Currency))
			continue
		}
	}
	return resultReports, nil
}

func (s *BondReporter) CreateGeneralBondReport(ctx context.Context, resultBondPosition *report.ReportPositions, totalAmount float64) (_ domain.GeneralBondReportPosition, err error) {
	const op = "service.CreateGeneralBondReport"
	defer logging.LogOperation_Debug(ctx, s.logger, op, &err)()
	// Переменная ткущих позиций формата []report.PositionByFIFO
	currentPostions := resultBondPosition.CurrentPositions
	// Создаем переменные
	var BondReporPosition domain.GeneralBondReportPosition
	var name string
	var duration int64
	var yieldToMaturity float64

	firstBondsBuyDate := currentPostions[0].BuyDate

	if len(currentPostions) == 0 {
		return domain.GeneralBondReportPosition{}, bondreport.ErrEmptyPositons
	}

	ticker := currentPostions[0].Ticker
	currency := currentPostions[0].Currency
	nominal := currentPostions[0].Nominal
	sellPrice := currentPostions[0].SellPrice
	replaced := currentPostions[0].Replaced

	sumOfPositionsRes := GetSumOfPositions(currentPostions)

	positionsPrice := utils.RoundFloat(sumOfPositionsRes.SumOfPositions/sumOfPositionsRes.SumOfQuantity, 2)
	quantity := int64(sumOfPositionsRes.SumOfQuantity)
	profitOfAllPositions := utils.RoundFloat(sumOfPositionsRes.ProfitOfAllPositions, 2)
	profitInPercentage := utils.RoundFloat((profitOfAllPositions/sumOfPositionsRes.SumOfPositions)*100, 2)
	percentOfPortfolio := utils.RoundFloat((sumOfPositionsRes.SumOfPositions/totalAmount)*100, 2)

	moexBuyDateData, err := s.GetSpecificationsFromMoex(ctx, ticker, firstBondsBuyDate)
	if err != nil {
		return domain.GeneralBondReportPosition{}, err
	}
	moexNowData, err := s.GetSpecificationsFromMoex(ctx, ticker, s.now())
	if err != nil {
		return domain.GeneralBondReportPosition{}, err
	}

	if moexNowData.ShortName.IsHasValue() {
		name = moexNowData.ShortName.Value
	} else {
		name = currentPostions[0].Name
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
		return domain.GeneralBondReportPosition{}, nil
	}

	BondReporPosition.Ticker = ticker
	BondReporPosition.Currencies = currency
	BondReporPosition.Nominal = nominal
	BondReporPosition.CurrentPrice = sellPrice
	BondReporPosition.Replaced = replaced
	BondReporPosition.PositionPrice = positionsPrice
	BondReporPosition.BuyDate = firstBondsBuyDate
	BondReporPosition.Quantity = quantity
	BondReporPosition.Profit = profitOfAllPositions
	BondReporPosition.ProfitInPercentage = profitInPercentage
	BondReporPosition.PercentOfPortfolio = percentOfPortfolio
	BondReporPosition.Name = name
	BondReporPosition.Duration = duration
	BondReporPosition.YieldToMaturity = yieldToMaturity
	BondReporPosition.MaturityDate = maturityDate

	if moexBuyDateData.YieldToOffer.IsHasValue() {
		yieldToOfferOnPurchase := moexBuyDateData.YieldToOffer.Value
		BondReporPosition.YieldToMaturityOnPurchase = yieldToOfferOnPurchase
	} else {
		if moexBuyDateData.YieldToMaturity.IsHasValue() {
			yieldToMaturityOnPurchase := moexBuyDateData.YieldToMaturity.Value
			BondReporPosition.YieldToMaturityOnPurchase = yieldToMaturityOnPurchase
		}
	}

	return BondReporPosition, nil
}

func GetSumOfPositions(positions []report.PositionByFIFO) *bondreport.SumOfPositions {
	var firstBondsBuyDate time.Time
	var sumOfPositions float64
	var sumOfQuantity float64
	var profitOfAllPositions float64
	for _, position := range positions {
		sumOfPositions += position.BuyPrice * position.Quantity
		compareDate := position.BuyDate

		if firstBondsBuyDate.After(compareDate) {
			firstBondsBuyDate = compareDate
		}
		sumOfQuantity += position.Quantity

		profitWithoutTax := getSecurityIncomeWithoutTax(position)
		totalTax := getTotalTaxFromPosition(position, profitWithoutTax)
		profitOfAllPositions += getSecurityIncome(profitWithoutTax, totalTax)
	}
	res := bondreport.NewSumOfPositons(firstBondsBuyDate, sumOfPositions, sumOfQuantity, profitOfAllPositions)

	return res
}

func (s *BondReporter) createBondReportByCurrency(ctx context.Context, position report.PositionByFIFO) (_ domain.BondReport, err error) {
	const op = "service.createBondReportByCurrency"

	defer logging.LogOperation_Debug(ctx, s.logger, op, &err)()

	var bondReport domain.BondReport
	moexBuyDateData, err := s.GetSpecificationsFromMoex(ctx, position.Ticker, position.BuyDate)
	if err != nil {
		return bondReport, errors.New("service: createBondReport" + err.Error())
	}
	date := time.Now()
	moexNowData, err := s.GetSpecificationsFromMoex(ctx, position.Ticker, date)
	if err != nil {
		return bondReport, errors.New("service: createBondReport" + err.Error())
	}

	bondReport = domain.BondReport{
		Name:         position.Name,
		Ticker:       position.Ticker,
		BuyDate:      position.BuyDate.Format(layout),
		BuyPrice:     utils.RoundFloat(position.BuyPrice, 2),
		CurrentPrice: utils.RoundFloat(position.SellPrice, 2),
		Nominal:      position.Nominal,
	}

	if moexNowData.MaturityDate.IsHasValue() {
		maturityDate := moexNowData.MaturityDate.Value
		bondReport.MaturityDate = maturityDate
	}

	if moexNowData.OfferDate.IsHasValue() {
		offerDate := moexNowData.OfferDate.Value
		bondReport.OfferDate = offerDate
	}

	if moexNowData.Duration.IsHasValue() {
		duration := moexNowData.Duration.Value
		bondReport.Duration = int64(duration)
	}

	if moexNowData.YieldToMaturity.IsHasValue() {
		yieldToMaturity := moexNowData.YieldToMaturity.Value
		bondReport.YieldToMaturity = yieldToMaturity
	}

	if moexNowData.YieldToOffer.IsHasValue() {
		yieldToOffer := moexNowData.YieldToOffer.Value
		bondReport.YieldToOffer = yieldToOffer
	}

	// -------------------

	if moexBuyDateData.YieldToMaturity.IsHasValue() {
		yieldToMaturityOnPurchase := moexBuyDateData.YieldToMaturity.Value
		bondReport.YieldToMaturityOnPurchase = utils.RoundFloat(yieldToMaturityOnPurchase, 2)
	}

	if moexBuyDateData.YieldToOffer.IsHasValue() {
		yieldToOfferOnPurchase := moexBuyDateData.YieldToOffer.Value
		bondReport.YieldToOfferOnPurchase = utils.RoundFloat(yieldToOfferOnPurchase, 2)
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

// TODO: Сделать это методом структуры PositionByFIFO
func getProfit(sharePosition report.PositionByFIFO) (_ float64, err error) {
	const op = "service.getProfit"

	profitWithoutTax := getSecurityIncomeWithoutTax(sharePosition)
	totalTax := getTotalTaxFromPosition(sharePosition, profitWithoutTax)
	profit := getSecurityIncome(profitWithoutTax, totalTax)
	profitInPercentage, err := getProfitInPercentage(profit, sharePosition.BuyPrice, sharePosition.Quantity)
	if err != nil {
		return profitInPercentage, err
	}
	return profitInPercentage, nil
}

func getAnnualizedReturnInPercentage(p report.PositionByFIFO) (_ float64, err error) {
	const op = "service.getAnnualizedReturnInPercentage"

	var totalReturn float64
	profitWithoutTax := getSecurityIncomeWithoutTax(p)
	totalTax := getTotalTaxFromPosition(p, profitWithoutTax)
	profit := getSecurityIncome(profitWithoutTax, totalTax)
	buyDate := p.BuyDate
	// Костыль. Надо переписать , так как для закрытх позиций данная функция работать не будет
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
	annualizedReturnInPercentageRound := utils.RoundFloat(annualizedReturnInPercentage, 2)

	return annualizedReturnInPercentageRound, nil
}

// Доход по позиции до вычета налога
func getSecurityIncomeWithoutTax(p report.PositionByFIFO) float64 {
	// Для незакрытых позиций необходимо посчитать еще потенциальную комиссию при продаже
	quantity := p.Quantity
	buySellDifference := (p.SellPrice-p.BuyPrice)*quantity + p.SellAccruedInt - p.BuyAccruedInt
	cashFlow := p.TotalCoupon + p.TotalDividend
	positionProfit := buySellDifference + cashFlow + p.TotalComission + p.PartialEarlyRepayment
	return positionProfit
}

// Расход полного налога по закрытой позиции
func getTotalTaxFromPosition(p report.PositionByFIFO, profit float64) float64 {
	// Рассчитываем срок владения
	buyDate := p.BuyDate
	sellDate := p.SellDate
	timeDuration := sellDate.Sub(buyDate).Hours()
	var tax float64
	// Рассчитываем налог с продажи бумаги, если сумма продажи больше суммы покупки
	if profit > 0 && timeDuration < float64(threeYearInHours) {
		tax = profit * baseTaxRate
	} else {
		tax = 0
	}
	return tax
}

// Расчет прибыли после налогообложения
func getSecurityIncome(profit, tax float64) float64 {
	profitAfterTax := profit - tax
	return profitAfterTax
}

func getProfitInPercentage(profit, buyPrice, quantity float64) (float64, error) {
	if buyPrice != 0 && quantity != 0 {
		profitInPercentage := utils.RoundFloat((profit/(buyPrice*quantity))*100, 2)
		return profitInPercentage, nil
	} else {
		return 0, errors.New("divide by zero")
	}
}

func getFirstMuturityDate(buyBackDate, offerDate, maturityDate domain.NullString) (time.Time, error) {
	var resDate time.Time

	switch {
	case offerDate.IsHasValue() && buyBackDate.IsHasValue():
		offerDateConv, err := time.Parse(layout, offerDate.Value)
		if err != nil {
			return time.Time{}, e.WrapIfErr("failed to parse offerDate", err)
		}
		buyBackDateConv, err := time.Parse(layout, buyBackDate.Value)
		if err != nil {
			return time.Time{}, e.WrapIfErr("failed to parse buyBackDate", err)
		}
		if buyBackDateConv.After(offerDateConv) {
			resDate = offerDateConv
		} else {
			resDate = buyBackDateConv
		}
	case offerDate.IsHasValue():
		offerDateConv, err := time.Parse(layout, offerDate.Value)
		if err != nil {
			return time.Time{}, e.WrapIfErr("failed to parse offerDate", err)
		}
		resDate = offerDateConv
	case buyBackDate.IsHasValue():
		buyBackDateConv, err := time.Parse(layout, buyBackDate.Value)
		if err != nil {
			return time.Time{}, e.WrapIfErr("failed to parse buyBackDate", err)
		}
		resDate = buyBackDateConv
	case maturityDate.IsHasValue():
		maturityDateConv, err := time.Parse(layout, maturityDate.Value)
		if err != nil {
			return time.Time{}, e.WrapIfErr("failed to parse maturityDate", err)
		}
		resDate = maturityDateConv
	}
	return resDate, nil
}
