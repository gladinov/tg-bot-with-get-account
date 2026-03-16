package report

import (
	"bonds-report-service/internal/domain"
	report "bonds-report-service/internal/domain/report_position"
	"bonds-report-service/internal/utils/logging"
	"context"
	"errors"
	"log/slog"

	"github.com/gladinov/e"
)

func (p *ReportProcessor) ProcessOperations(ctx context.Context, reportLine *domain.ReportLine) (_ *report.ReportPositions, err error) {
	const op = "report.ProcessOperations"

	defer logging.LogOperation_Debug(ctx, p.logger, op, &err)()

	processPosition := report.NewReportPositons()

	for _, operation := range reportLine.Operation {
		if err := processPosition.Apply(
			operation,
			reportLine.Bond,
			reportLine.LastPrice,
			reportLine.Vunit_rate); err != nil {
			if errors.Is(err, report.ErrUnknownOpp) {
				p.logger.WarnContext(ctx, "unkown opperation type", slog.String("op", op))
				continue
			}
			if errors.Is(err, report.ErrZeroQuantity) {
				continue
			}
			return nil, e.WrapIfErr("failed to apply operation to processPosition", err)
		}
	}
	return processPosition, nil
}
