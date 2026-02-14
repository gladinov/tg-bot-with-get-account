package app

import (
	"bonds-report-service/internal/service/report"
	"log/slog"
)

func InitReportProcessor(logger *slog.Logger) *report.ReportProcessor {
	logger.Info("initialize report processor")
	reportProcessor := report.NewReportProcessor(logger)
	return reportProcessor
}
