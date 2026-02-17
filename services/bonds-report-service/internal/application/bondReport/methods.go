package bondreport

import (
	"bonds-report-service/internal/domain"
	report "bonds-report-service/internal/domain/report"
	report_position "bonds-report-service/internal/domain/report_position"
	"bonds-report-service/internal/utils/logging"
	"context"
)

func (s *BondReporter) CreateBondReport(
	ctx context.Context,
	currentPositions []report_position.PositionByFIFO,
	moexBuyDateData domain.ValuesMoex,
	moexNowData domain.ValuesMoex,
) (_ report.Report, err error) {
	const op = "service.CreateBondReport"

	defer logging.LogOperation_Debug(ctx, s.logger, op, &err)()

	var resultReports report.Report

	for i := range currentPositions {
		position := &currentPositions[i]
		err := resultReports.Add(position, moexBuyDateData, moexNowData, s.now())
		if err != nil {
			return report.Report{}, err
		}
	}
	return resultReports, nil
}
