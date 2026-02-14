package report

import (
	"log/slog"
)

type ReportProcessor struct {
	logger *slog.Logger
}

func NewReportProcessor(logger *slog.Logger) *ReportProcessor {
	return &ReportProcessor{logger: logger}
}
