package usecases

import (
	"bonds-report-service/internal/application/dto"
	"bonds-report-service/internal/application/presenter"
	"bonds-report-service/internal/domain"
	"bonds-report-service/internal/domain/generalbondreport"
	"bonds-report-service/internal/utils/logging"
	"context"
	"errors"
	"log/slog"
	"sort"
	"sync"

	"github.com/gladinov/e"
)

var ErrEmptyBondPositions = errors.New("len of result bond positions is empty")

func (s *Service) GetBondReports(ctx context.Context, chatID int) (_ dto.BondReportsResponce, err error) {
	const op = "service.GetBondReports"
	defer logging.LogOperation_Debug(ctx, s.logger, op, &err)()

	reportsInByteByAccounts := make([][]*dto.MediaGroup, 0)

	accounts, err := s.Helpers.TinkoffHelper.TinkoffGetAccounts(ctx)
	if err != nil {
		return dto.BondReportsResponce{}, e.WrapIfErr("failde to get accounts from Tinkoff", err)
	}

	for _, account := range accounts {
		if account.Status != 2 {
			continue
		}
		reportsInByte, err := s.processAccount(ctx, chatID, account)
		if err != nil {
			return dto.BondReportsResponce{}, e.WrapIfErr("failed to procces account", err)
		}
		reportsInByteByAccounts = append(reportsInByteByAccounts, reportsInByte)
	}
	getBondReportsResponce := dto.BondReportsResponce{Media: reportsInByteByAccounts}
	return getBondReportsResponce, nil
}

func (s *Service) processAccount(ctx context.Context, chatID int, account domain.Account) (_ []*dto.MediaGroup, err error) {
	const op = "service.processAccount"
	defer logging.LogOperation_Debug(ctx, s.logger, op, &err)()

	err = s.Helpers.OperationsUpdater.UpdateOperations(ctx, chatID, account.ID, account.OpenedDate)
	if err != nil {
		return nil, e.WrapIfErr("failed to update operations", err)
	}
	portfolio, err := s.Helpers.TinkoffHelper.TinkoffGetPortfolio(ctx, account)
	if err != nil {
		return nil, e.WrapIfErr("failed to get portfolio from tinkoff", err)
	}

	totalAmount := portfolio.TotalAmount.ToFloat()

	portfolioPositions, err := s.Helpers.PositionProcessor.ProcessPositionsToPositionsWithAssetUid(ctx, portfolio.Positions)
	if err != nil {
		return nil, e.WrapIfErr("failed to process position to postition with asset uid", err)
	}

	allOperations, err := s.Storage.GetAllOperations(ctx, chatID, account.ID)
	if err != nil {
		return nil, e.WrapIfErr("failed to get all operations from storage", err)
	}

	operationsByAssetUid := make(map[string][]domain.OperationWithoutCustomTypes)

	for _, op := range allOperations {
		operationsByAssetUid[op.AssetUid] = append(operationsByAssetUid[op.AssetUid], op)
	}

	generalBondReports := generalbondreport.NewGeneralBondReports()

	for _, position := range portfolioPositions {
		if position.InstrumentType != "bond" {
			continue
		}
		bondReport, err := s.processBondPosition(ctx, position, totalAmount, operationsByAssetUid[position.AssetUid])
		if err != nil {
			if errors.Is(err, ErrEmptyBondPositions) {
				continue
			}
			return nil, e.WrapIfErr("failed to process bond position", err)
		}

		s.addBondReport(&generalBondReports, bondReport)

	}

	reportsInByte, err := presenter.GenerateTablePNG(ctx,
		s.logger,
		s.prepareToGenerateTablePNG(ctx, &generalBondReports),
		chatID,
		account.ID)
	if err != nil {
		return nil, e.WrapIfErr("failed to GenerateTablePNG", err)
	}

	return reportsInByte, nil
}

func (s *Service) processBondPosition(ctx context.Context,
	position domain.PortfolioPositionsWithAssetUid,
	totalAmount float64,
	operationsDb []domain.OperationWithoutCustomTypes,
) (_ generalbondreport.GeneralBondReportPosition, err error) {
	const op = "service.processBondPosition"
	defer logging.LogOperation_Debug(ctx, s.logger, op, &err)()

	reporLines, err := s.Helpers.ReportLineBuilder.CreateNewReportLines(ctx, position, operationsDb)
	if err != nil {
		return generalbondreport.GeneralBondReportPosition{}, e.WrapIfErr("failed to create new report lines", err)
	}
	// Обрабатываем операции и получаем открытые позиции по данной бумаге
	resultBondPosition, err := s.Helpers.ReportProcessor.ProcessOperations(ctx, reporLines)
	if err != nil {
		return generalbondreport.GeneralBondReportPosition{}, e.WrapIfErr("failed to process operation", err)
	}
	// Общая стоимость портфеля

	if len(resultBondPosition.CurrentPositions) == 0 {
		s.logger.WarnContext(
			ctx,
			"len of result bond positions is empty. But in portfolio request to TinkoffApi position in portfolio",
			slog.Any("postion", position))
		return generalbondreport.GeneralBondReportPosition{}, ErrEmptyBondPositions
	}

	firstBuyDate := resultBondPosition.CurrentPositions[0].BuyDate

	// TODO : Кэшифрование запросов в MOEX
	moexBuyDateData, err := s.Helpers.MoexSpecificationGetter.GetSpecificationsFromMoex(ctx, reporLines.Bond.Ticker, firstBuyDate)
	if err != nil {
		return generalbondreport.GeneralBondReportPosition{}, e.WrapIfErr("failed to get specification from MOEX for buy date", err)
	}
	moexNowData, err := s.Helpers.MoexSpecificationGetter.GetSpecificationsFromMoex(ctx, reporLines.Bond.Ticker, s.now())
	if err != nil {
		return generalbondreport.GeneralBondReportPosition{}, e.WrapIfErr("failed to get specification from MOEX for current date", err)
	}

	bondReport, err := s.Helpers.GeneralBondReportProcessor.GetGeneralBondReportPosition(
		ctx,
		resultBondPosition.CurrentPositions,
		totalAmount,
		moexBuyDateData,
		moexNowData, firstBuyDate)
	if err != nil {
		return generalbondreport.GeneralBondReportPosition{}, e.WrapIfErr("failed to get general bond report position", err)
	}
	return bondReport, nil
}

func (s *Service) prepareToGenerateTablePNG(ctx context.Context,
	generalBondReports *generalbondreport.GeneralBondReports,
) (_ [][]generalbondreport.GeneralBondReportPosition) {
	const op = "service.PrepareToGenerateTablePNG"

	defer logging.LogOperation_Debug(ctx, s.logger, op, nil)()
	var wg sync.WaitGroup
	reports := make([][]generalbondreport.GeneralBondReportPosition, 3)

	wg.Add(1)
	go func() {
		defer wg.Done()
		reports[0] = s.sortGeneralBondReports(ctx, generalBondReports.RubBondsReport)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		reports[1] = s.sortGeneralBondReports(ctx, generalBondReports.ReplacedBondsReport)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		reports[2] = s.sortGeneralBondReports(ctx, generalBondReports.EuroBondsReport)
	}()

	wg.Wait()

	return reports
}

func (s *Service) sortGeneralBondReports(ctx context.Context,
	report map[generalbondreport.TickerTimeKey]generalbondreport.GeneralBondReportPosition,
) (_ []generalbondreport.GeneralBondReportPosition) {
	const op = "service.sortGeneralBondReports"

	defer logging.LogOperation_Debug(ctx, s.logger, op, nil)()

	keys := make([]generalbondreport.TickerTimeKey, 0, len(report))
	for k := range report {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Time.Equal(keys[j].Time) {
			return keys[i].Ticker < keys[j].Ticker
		}
		return keys[i].Time.Before(keys[j].Time)
	})
	result := make([]generalbondreport.GeneralBondReportPosition, len(keys))
	for i, k := range keys {
		result[i] = report[k]
	}

	return result
}

func (s *Service) addBondReport(generalBondReports *generalbondreport.GeneralBondReports, bondReport generalbondreport.GeneralBondReportPosition) {
	switch {
	case bondReport.Replaced:
		tickerTimeKey := generalbondreport.NewTickerTimeKey(bondReport.Ticker, bondReport.BuyDate)
		generalBondReports.ReplacedBondsReport[tickerTimeKey] = bondReport
	case bondReport.Currencies != "rub":
		tickerTimeKey := generalbondreport.NewTickerTimeKey(bondReport.Ticker, bondReport.BuyDate)
		generalBondReports.EuroBondsReport[tickerTimeKey] = bondReport
	default:
		tickerTimeKey := generalbondreport.NewTickerTimeKey(bondReport.Ticker, bondReport.BuyDate)
		generalBondReports.RubBondsReport[tickerTimeKey] = bondReport
	}
}
