package app

import (
	bondreport "bonds-report-service/internal/application/bondReport"
	generalbondreport "bonds-report-service/internal/application/generalBondReport"
	"bonds-report-service/internal/application/report"
	"log/slog"
)

func InitReportProcessor(logger *slog.Logger) *report.ReportProcessor {
	logger.Info("initialize report processor")
	reportProcessor := report.NewReportProcessor(logger)
	return reportProcessor
}

func InitBondReportProcessor(logger *slog.Logger) *bondreport.BondReporter {
	logger.Info("initialize bond report processor")
	bondReporter := bondreport.NewBondReporter(logger)
	return bondReporter
}

func InitGeneralReportProcessor(logger *slog.Logger) *generalbondreport.GeneralBondReporter {
	logger.Info("initialize general bond report processor")
	generalBondReporter := generalbondreport.NewGeneralBondReporter(logger)
	return generalBondReporter
}
