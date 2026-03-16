package generalbondreport

import (
	"bonds-report-service/internal/domain"
	report "bonds-report-service/internal/domain/report_position"
	"time"

	"github.com/gladinov/e"
)

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

func getSumOfPositions(positions []report.PositionByFIFO) *SumOfPositions {
	var sumOfPositions float64
	var sumOfQuantity float64
	var profitOfAllPositions float64
	for _, position := range positions {
		sumOfPositions += position.BuyPrice * position.Quantity

		sumOfQuantity += position.Quantity

		profitWithoutTax := position.GetProfitBeforeTax()
		totalTax := position.GetTotalTaxFromPosition(profitWithoutTax)
		profitOfAllPositions += report.GetNetProfit(profitWithoutTax, totalTax)
	}
	res := NewSumOfPositons(sumOfPositions, sumOfQuantity, profitOfAllPositions)

	return res
}
