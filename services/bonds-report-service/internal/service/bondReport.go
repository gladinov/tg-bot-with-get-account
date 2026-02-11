package service

import (
	"bonds-report-service/internal/service/service_models"
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

func (c *Client) CreateBondReport(ctx context.Context, reportPostions service_models.ReportPositions) (_ service_models.Report, err error) {
	const op = "service.CreateBondReport"

	var resultReports service_models.Report
	for i := range reportPostions.CurrentPositions {
		position := reportPostions.CurrentPositions[i]
		switch position.Currency {
		case "rub":
			bondReport, err := c.createBondReportByCurrency(ctx, position)
			if err != nil {
				return resultReports, errors.New("service: GetBondReport" + err.Error())
			}
			resultReports.BondsInRUB = append(resultReports.BondsInRUB, bondReport)
		case "cny":
			bondReport, err := c.createBondReportByCurrency(ctx, position)
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

func (c *Client) CreateGeneralBondReport(ctx context.Context, resultBondPosition *service_models.ReportPositions, totalAmount float64) (_ service_models.GeneralBondReportPosition, err error) {
	const op = "service.CreateGeneralBondReport"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("can't create general bond report", err)
	}()

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

	moexBuyDateData, err := c.External.Moex.GetSpecifications(ctx, BondReporPosition.Ticker, buyDate)
	if err != nil {
		return BondReporPosition, err
	}
	date := time.Now()
	moexNowData, err := c.External.Moex.GetSpecifications(ctx, BondReporPosition.Ticker, date)
	if err != nil {
		return BondReporPosition, err
	}

	if moexNowData.ShortName.IsHasValue() {
		BondReporPosition.Name = moexNowData.ShortName.Value
	} else {
		BondReporPosition.Name = currentPostions[0].Name
	}

	if moexNowData.Duration.IsHasValue() {
		BondReporPosition.Duration = int64(moexNowData.Duration.Value)
	}

	if moexNowData.YieldToOffer.IsHasValue() {
		BondReporPosition.YieldToMaturity = moexNowData.YieldToOffer.Value
	} else {
		if moexNowData.YieldToMaturity.IsHasValue() {
			BondReporPosition.YieldToMaturity = moexNowData.YieldToMaturity.Value
		}
	}
	var maturityDate string
	if moexNowData.MaturityDate.IsHasValue() {
		maturityDate = moexNowData.MaturityDate.Value
	}

	var offerDate string
	if moexNowData.OfferDate.IsHasValue() {
		offerDate = moexNowData.OfferDate.Value
	}

	var buyBackDate string
	if moexNowData.BuybackDate.IsHasValue() {
		buyBackDate = moexNowData.BuybackDate.Value
	}

	switch {
	case moexNowData.OfferDate.IsHasValue() && moexNowData.BuybackDate.IsHasValue():
		offerDateConv, err := time.Parse(layout, offerDate)
		if err != nil {
			return BondReporPosition, err
		}
		buyBackDateConv, err := time.Parse(layout, buyBackDate)
		if err != nil {
			return BondReporPosition, err
		}
		if offerDateConv.Compare(buyBackDateConv) == -1 {
			BondReporPosition.MaturityDate = offerDateConv
		} else {
			BondReporPosition.MaturityDate = buyBackDateConv
		}
	case moexNowData.OfferDate.IsHasValue():
		offerDateConv, err := time.Parse(layout, offerDate)
		if err != nil {
			return BondReporPosition, err
		}
		BondReporPosition.MaturityDate = offerDateConv
	case moexNowData.BuybackDate.IsHasValue():
		buyBackDateConv, err := time.Parse(layout, buyBackDate)
		if err != nil {
			return BondReporPosition, err
		}
		BondReporPosition.MaturityDate = buyBackDateConv
	case moexNowData.MaturityDate.IsHasValue():
		maturityDateConv, err := time.Parse(layout, maturityDate)
		if err != nil {
			return BondReporPosition, err
		}
		BondReporPosition.MaturityDate = maturityDateConv
	}

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

func (c *Client) createBondReportByCurrency(ctx context.Context, position service_models.PositionByFIFO) (_ service_models.BondReport, err error) {
	const op = "service.createBondReportByCurrency"

	start := time.Now()
	logg := c.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("can't create bond report by currency", err)
	}()

	var bondReport service_models.BondReport
	moexBuyDateData, err := c.External.Moex.GetSpecifications(ctx, position.Ticker, position.BuyDate)
	if err != nil {
		return bondReport, errors.New("service: createBondReport" + err.Error())
	}
	date := time.Now()
	moexNowData, err := c.External.Moex.GetSpecifications(ctx, position.Ticker, date)
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
		bondReport.YieldToMaturityOnPurchase = RoundFloat(yieldToMaturityOnPurchase, 2)
	}

	if moexBuyDateData.YieldToOffer.IsHasValue() {
		yieldToOfferOnPurchase := moexBuyDateData.YieldToOffer.Value
		bondReport.YieldToOfferOnPurchase = RoundFloat(yieldToOfferOnPurchase, 2)
	}

	profitInPercentage, err := getProfit(logg, position)
	if err != nil {
		return bondReport, errors.New("service: createBondReport" + err.Error())
	}
	bondReport.Profit = profitInPercentage

	annualizedReturn, err := getAnnualizedReturnInPercentage(logg, position)
	if err != nil {
		return bondReport, errors.New("service: createBondReport" + err.Error())
	}
	bondReport.AnnualizedReturn = annualizedReturn

	return bondReport, nil
}

func getProfit(logger *slog.Logger, sharePosition service_models.PositionByFIFO) (_ float64, err error) {
	const op = "service.getProfit"

	start := time.Now()
	logg := logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("can't get prorfit ", err)
	}()
	profitWithoutTax := getSecurityIncomeWithoutTax(sharePosition)
	totalTax := getTotalTaxFromPosition(sharePosition, profitWithoutTax)
	profit := getSecurityIncome(profitWithoutTax, totalTax)
	profitInPercentage, err := getProfitInPercentage(profit, sharePosition.BuyPrice, sharePosition.Quantity)
	if err != nil {
		return profitInPercentage, err
	}
	return profitInPercentage, nil
}

func getAnnualizedReturnInPercentage(logger *slog.Logger, p service_models.PositionByFIFO) (_ float64, err error) {
	const op = "service.getAnnualizedReturnInPercentage"

	start := time.Now()
	logg := logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("can't get Annualized Return In Percentage", err)
	}()

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
	annualizedReturnInPercentageRound := RoundFloat(annualizedReturnInPercentage, 2)

	return annualizedReturnInPercentageRound, nil
}
