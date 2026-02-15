package bondreport

import (
	"time"
)

type SumOfPositions struct {
	FirstBondsBuyDate    time.Time
	SumOfPositions       float64
	SumOfQuantity        float64
	ProfitOfAllPositions float64
}

func NewSumOfPositons(firstBondsBuyDate time.Time,
	sumOfPositions float64,
	sumOfQuantity float64,
	profitOfAllPositions float64,
) *SumOfPositions {
	return &SumOfPositions{
		FirstBondsBuyDate:    firstBondsBuyDate,
		SumOfPositions:       sumOfPositions,
		SumOfQuantity:        sumOfQuantity,
		ProfitOfAllPositions: profitOfAllPositions,
	}
}
