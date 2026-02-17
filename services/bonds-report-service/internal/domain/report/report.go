package report

import (
	"bonds-report-service/internal/domain"
	report_position "bonds-report-service/internal/domain/report_position"
	"bonds-report-service/internal/utils"
	"errors"
	"time"

	"github.com/gladinov/e"
)

func (r *Report) Add(
	position *report_position.PositionByFIFO,
	moexBuyDateData domain.ValuesMoex,
	moexNowData domain.ValuesMoex,
	now time.Time,
) (err error) {
	const op = "service.CreateBondReport"

	bondReportPosition := BondReport{}

	switch position.Currency {
	case "rub":
		err := bondReportPosition.createBondReportPosition(position, moexBuyDateData, moexNowData, now)
		if err != nil {
			return e.WrapIfErr("failed to create Bond report position", err)
		}
		r.BondsInRUB = append(r.BondsInRUB, bondReportPosition)
	case "cny":
		err := bondReportPosition.createBondReportPosition(position, moexBuyDateData, moexNowData, now)
		if err != nil {
			return e.WrapIfErr("failed to create Bond report position", err)
		}
		r.BondsInCNY = append(r.BondsInCNY, bondReportPosition)
	}

	return nil
}

func (r *BondReport) createBondReportPosition(
	position *report_position.PositionByFIFO,
	moexBuyDateData domain.ValuesMoex,
	moexNowData domain.ValuesMoex,
	now time.Time,
) (err error) {
	const op = "service.createBondReportByCurrency"

	var maturityDate string
	var offerDate string
	var duration float64
	var yieldToMaturity float64
	var yieldToOffer float64
	var yieldToMaturityOnPurchase float64
	var yieldToOfferOnPurchase float64

	name := position.Name
	ticker := position.Ticker
	buyDate := position.BuyDate.Format(layout)
	buyPrice := utils.RoundFloat(position.BuyPrice, 2)
	currentPrice := utils.RoundFloat(position.SellPrice, 2)
	nominal := position.Nominal

	if moexNowData.MaturityDate.IsHasValue() {
		maturityDate = moexNowData.MaturityDate.Value
	}

	if moexNowData.OfferDate.IsHasValue() {
		offerDate = moexNowData.OfferDate.Value
	}

	if moexNowData.Duration.IsHasValue() {
		duration = moexNowData.Duration.Value
	}

	if moexNowData.YieldToMaturity.IsHasValue() {
		yieldToMaturity = moexNowData.YieldToMaturity.Value
	}

	if moexNowData.YieldToOffer.IsHasValue() {
		yieldToOffer = moexNowData.YieldToOffer.Value
	}

	if moexBuyDateData.YieldToMaturity.IsHasValue() {
		yieldToMaturityOnPurchase = moexBuyDateData.YieldToMaturity.Value
	}

	if moexBuyDateData.YieldToOffer.IsHasValue() {
		yieldToOfferOnPurchase = moexBuyDateData.YieldToOffer.Value
	}
	profitWithoutTax := position.GetProfitBeforeTax()
	totalTax := position.GetTotalTaxFromPosition(profitWithoutTax)
	profit := getNetProfit(profitWithoutTax, totalTax)

	profitInPercentage, err := position.GetProfit(profit)
	if err != nil {
		return errors.New("service: createBondReport" + err.Error())
	}
	r.Profit = profitInPercentage

	annualizedReturn, err := position.GetAnnualizedReturnInPercentage(profit, now)
	if err != nil {
		return errors.New("service: createBondReport" + err.Error())
	}
	r.AnnualizedReturn = annualizedReturn
	r.Name = name
	r.Ticker = ticker
	r.BuyDate = buyDate
	r.BuyPrice = buyPrice
	r.CurrentPrice = currentPrice
	r.Nominal = nominal
	r.MaturityDate = maturityDate
	r.OfferDate = offerDate
	r.Duration = int64(duration)
	r.YieldToMaturity = yieldToMaturity
	r.YieldToOffer = yieldToOffer
	r.YieldToMaturityOnPurchase = utils.RoundFloat(yieldToMaturityOnPurchase, 2)
	r.YieldToOfferOnPurchase = utils.RoundFloat(yieldToOfferOnPurchase, 2)

	return nil
}
