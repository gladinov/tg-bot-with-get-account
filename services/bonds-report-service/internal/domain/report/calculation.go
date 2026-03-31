package report

func getNetProfit(profit, tax float64) float64 {
	profitAfterTax := profit - tax
	return profitAfterTax
}
