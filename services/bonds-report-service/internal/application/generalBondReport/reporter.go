package generalbondreport

import "log/slog"

type GeneralBondReporter struct {
	logger *slog.Logger
}

func NewGeneralBondReporter(logg *slog.Logger) *GeneralBondReporter {
	return &GeneralBondReporter{logger: logg}
}
