package generalbondreport

import "time"

type GeneralBondReports struct {
	RubBondsReport      map[TickerTimeKey]GeneralBondReportPosition
	EuroBondsReport     map[TickerTimeKey]GeneralBondReportPosition
	ReplacedBondsReport map[TickerTimeKey]GeneralBondReportPosition
}

type TickerTimeKey struct {
	Ticker string
	Time   time.Time
}

type SumOfPositions struct {
	SumOfPositions       float64
	SumOfQuantity        float64
	ProfitOfAllPositions float64
}

func NewSumOfPositons(
	sumOfPositions float64,
	sumOfQuantity float64,
	profitOfAllPositions float64,
) *SumOfPositions {
	return &SumOfPositions{
		SumOfPositions:       sumOfPositions,
		SumOfQuantity:        sumOfQuantity,
		ProfitOfAllPositions: profitOfAllPositions,
	}
}
