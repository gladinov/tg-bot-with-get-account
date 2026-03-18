package usecases

import (
	"bonds-report-service/internal/domain"
	"bonds-report-service/internal/domain/report"
	"bonds-report-service/internal/utils/logging"
	"context"
	"log/slog"

	"github.com/gladinov/e"
)

func (s *Service) GetBondReportsByFifo(ctx context.Context, chatID int) (err error) {
	const op = "service.GetBondReportsByFifo"
	logg := slog.With(s.logger)
	defer logging.LogOperation_Debug(ctx, logg, op, &err)()

	accounts, err := s.Helpers.TinkoffHelper.TinkoffGetAccounts(ctx)
	if err != nil {
		return e.WrapIfErr("get accounts error", err)
	}

	for _, account := range accounts {
		if !isActiveAccounts(account) {
			continue
		}
		err := s.processAccountForBondReportByFifo(ctx, chatID, account)
		if err != nil {
			e.WrapIfErr("failed to process account", err)
		}

	}
	return nil
}

func (s *Service) processAccountForBondReportByFifo(ctx context.Context, chatID int, account domain.Account) error {
	err := s.Helpers.OperationsUpdater.UpdateOperations(ctx, chatID, account.ID, account.OpenedDate)
	if err != nil {
		return e.WrapIfErr("update operation error", err)
	}

	portfolio, err := s.Helpers.TinkoffHelper.TinkoffGetPortfolio(ctx, account)
	if err != nil {
		return e.WrapIfErr("tinkoffGetPortfolio err", err)
	}

	portfolioPositions, err := s.Helpers.PositionProcessor.ProcessPositionsToPositionsWithAssetUid(ctx, portfolio.Positions)
	if err != nil {
		return e.WrapIfErr("transformPositions err", err)
	}
	err = s.Storage.DeleteBondReport(ctx, chatID, account.ID)
	if err != nil {
		return e.WrapIfErr("deleteBondReport err", err)
	}

	operationsDb, err := s.Storage.GetAllOperations(ctx, chatID, account.ID)
	if err != nil {
		return e.WrapIfErr("failed to get all operations ", err)
	}

	operationsByAssetUid := make(map[string][]domain.OperationWithoutCustomTypes)

	for _, operations := range operationsDb {
		operationsByAssetUid[operations.AssetUid] = append(operationsByAssetUid[operations.AssetUid], operations)
	}

	bondsInRub := make([]report.BondReport, 0)

	var bondReport report.Report

	for _, position := range portfolioPositions {
		if !isBondType(position) {
			continue
		}
		operationsOfPosition := operationsByAssetUid[position.AssetUid]
		bondReport, err = s.processPositionsForBondReportByFifo(ctx, chatID, position, account, operationsOfPosition)

		bondsInRub = append(bondsInRub, bondReport.BondsInRUB...)

	}

	err = s.Storage.SaveBondReport(ctx, chatID, account.ID, bondsInRub)
	if err != nil {
		return e.WrapIfErr("Storage.SaveBondReport error", err)
	}
	return nil
}

func (s *Service) processPositionsForBondReportByFifo(ctx context.Context, chatID int, position domain.PortfolioPositionsWithAssetUid, account domain.Account, operationDbByAssetUid []domain.OperationWithoutCustomTypes) (report.Report, error) {
	reporLines, err := s.Helpers.ReportLineBuilder.CreateNewReportLines(ctx, position, operationDbByAssetUid)
	if err != nil {
		return report.Report{}, e.WrapIfErr("failed to create new report lines", err)
	}

	resultBondPosition, err := s.Helpers.ReportProcessor.ProcessOperations(ctx, reporLines)
	if err != nil {
		return report.Report{}, e.WrapIfErr("failed to process operation", err)
	}

	firstBuyDate := resultBondPosition.CurrentPositions[0].BuyDate

	moexBuyDateData, err := s.Helpers.MoexSpecificationGetter.GetSpecificationsFromMoex(ctx, reporLines.Bond.Ticker, firstBuyDate)
	if err != nil {
		return report.Report{}, e.WrapIfErr("failed to get specifications from moex to buy date", err)
	}
	moexNowData, err := s.Helpers.MoexSpecificationGetter.GetSpecificationsFromMoex(ctx, reporLines.Bond.Ticker, firstBuyDate)
	if err != nil {
		return report.Report{}, e.WrapIfErr("failed to get specifications from moex to buy now", err)
	}

	bondReport, err := s.Helpers.BondReportProcessor.CreateBondReport(
		ctx,
		resultBondPosition.CurrentPositions,
		moexBuyDateData,
		moexNowData,
	)
	if err != nil {
		return report.Report{}, e.WrapIfErr("CreateBondReport error", err)
	}
	return bondReport, nil
}
