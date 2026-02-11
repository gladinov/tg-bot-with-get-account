package app

import (
	"bonds-report-service/internal/repository"
	"bonds-report-service/internal/service"
	"bonds-report-service/internal/service/uidprovider"
	"log/slog"
)

func InitUidProvider(logger *slog.Logger, repo repository.Storage, analyticService service.TinkoffAnalyticsClient) *uidprovider.UidProvider {
	logger.Info("initialize uid provider")
	uidProvider := uidprovider.NewUidProvider(repo, analyticService)
	return uidProvider
}
