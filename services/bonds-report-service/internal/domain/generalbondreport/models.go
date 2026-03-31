package generalbondreport

import "time"

type GeneralBondReports struct {
	RubBondsReport      map[TickerTimeKey]GeneralBondReportPosition
	EuroBondsReport     map[TickerTimeKey]GeneralBondReportPosition
	ReplacedBondsReport map[TickerTimeKey]GeneralBondReportPosition
}

func NewGeneralBondReports() GeneralBondReports {
	return GeneralBondReports{
		RubBondsReport:      make(map[TickerTimeKey]GeneralBondReportPosition),
		EuroBondsReport:     make(map[TickerTimeKey]GeneralBondReportPosition),
		ReplacedBondsReport: make(map[TickerTimeKey]GeneralBondReportPosition),
	}
}

type TickerTimeKey struct {
	Ticker string
	Time   time.Time
}

func NewTickerTimeKey(ticker string, time time.Time) TickerTimeKey {
	return TickerTimeKey{
		Ticker: ticker,
		Time:   time,
	}
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
