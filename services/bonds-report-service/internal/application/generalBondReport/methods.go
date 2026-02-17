package generalbondreport

import (
	"bonds-report-service/internal/domain"
	"bonds-report-service/internal/domain/generalbondreport"
	report_position "bonds-report-service/internal/domain/report_position"
	"bonds-report-service/internal/utils/logging"
	"context"
	"time"

	"github.com/gladinov/e"
)

func (r *GeneralBondReporter) GetGeneralBondReportPosition(
	ctx context.Context,
	currentPositions []report_position.PositionByFIFO,
	totalAmount float64,
	moexBuyDateData domain.ValuesMoex,
	moexNowData domain.ValuesMoex,
	firstBondsBuyDate time.Time,
) (_ generalbondreport.GeneralBondReportPosition, err error) {
	const op = "generalbondreport.GetGeneralBondReportPosition"
	defer logging.LogOperation_Debug(ctx, r.logger, op, &err)()

	var reportPosition generalbondreport.GeneralBondReportPosition
	if err := reportPosition.CreateGeneralBondReportPosition(
		currentPositions,
		totalAmount,
		moexBuyDateData,
		moexNowData,
		firstBondsBuyDate,
	); err != nil {
		return generalbondreport.GeneralBondReportPosition{}, e.WrapIfErr("failed to create general bond report position", err)
	}
	return reportPosition, nil
}
