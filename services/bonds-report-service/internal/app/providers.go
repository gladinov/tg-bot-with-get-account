package app

import (
	service "bonds-report-service/internal/application"
	"bonds-report-service/internal/application/uidprovider"
	"log/slog"
)

func InitUidProvider(logger *slog.Logger, repo service.Storage, analyticService service.TinkoffAnalyticsClient) *uidprovider.UidProvider {
	logger.Info("initialize uid provider")
	uidProvider := uidprovider.NewUidProvider(repo, analyticService)
	return uidProvider
}
