package report

import (
	"bonds-report-service/internal/domain"
	"bonds-report-service/internal/utils"
	"errors"
	"math"
)

func CalculateNominal(
	nominal domain.MoneyValue,
	replaced bool,
	rate domain.Rate,
) float64 {
	base := nominal.ToFloat()

	if !replaced {
		return base
	}

	if rate.Vunit_Rate.IsSet && !rate.Vunit_Rate.IsNull {
		return base * rate.Vunit_Rate.Value
	}

	return base
}

func CalculateSellPrice(
	nominal float64,
	lastPrice domain.LastPrice,
) float64 {
	return math.Round(
		lastPrice.LastPrice.ToFloat()/100*nominal*100,
	) / 100
}

// Расчет прибыли после налогообложения
func GetSecurityIncome(profit, tax float64) float64 {
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
