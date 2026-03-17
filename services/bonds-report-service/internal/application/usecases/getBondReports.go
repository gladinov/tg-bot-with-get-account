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
	"runtime"
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
		generalBondReports, err := s.processAccount(ctx, chatID, account)
		if err != nil {
			return dto.BondReportsResponce{}, e.WrapIfErr("failed to procces account", err)
		}

		reportsInByte, err := presenter.GenerateTablePNG(ctx,
			s.logger,
			s.prepareToGenerateTablePNG(ctx, &generalBondReports),
			chatID,
			account.ID)
		if err != nil {
			return dto.BondReportsResponce{}, e.WrapIfErr("failed to GenerateTablePNG", err)
		}

		reportsInByteByAccounts = append(reportsInByteByAccounts, reportsInByte)
	}
	getBondReportsResponce := dto.BondReportsResponce{Media: reportsInByteByAccounts}
	return getBondReportsResponce, nil
}

func (s *Service) processAccount(ctx context.Context, chatID int, account domain.Account) (_ generalbondreport.GeneralBondReports, err error) {
	const op = "service.processAccount"
	defer logging.LogOperation_Debug(ctx, s.logger, op, &err)()

	err = s.Helpers.OperationsUpdater.UpdateOperations(ctx, chatID, account.ID, account.OpenedDate)
	if err != nil {
		return generalbondreport.GeneralBondReports{}, e.WrapIfErr("failed to update operations", err)
	}
	portfolio, err := s.Helpers.TinkoffHelper.TinkoffGetPortfolio(ctx, account)
	if err != nil {
		return generalbondreport.GeneralBondReports{}, e.WrapIfErr("failed to get portfolio from tinkoff", err)
	}

	totalAmount := portfolio.TotalAmount.ToFloat()

	portfolioPositions, err := s.Helpers.PositionProcessor.ProcessPositionsToPositionsWithAssetUid(ctx, portfolio.Positions)
	if err != nil {
		return generalbondreport.GeneralBondReports{}, e.WrapIfErr("failed to process position to postition with asset uid", err)
	}

	allOperations, err := s.Storage.GetAllOperations(ctx, chatID, account.ID)
	if err != nil {
		return generalbondreport.GeneralBondReports{}, e.WrapIfErr("failed to get all operations from storage", err)
	}

	operationsByAssetUid := make(map[string][]domain.OperationWithoutCustomTypes)

	for _, op := range allOperations {
		operationsByAssetUid[op.AssetUid] = append(operationsByAssetUid[op.AssetUid], op)
	}

	generalBondReports := generalbondreport.NewGeneralBondReports()

	workers := runtime.NumCPU()
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()
	workList := make(chan domain.PortfolioPositionsWithAssetUid, workers)
	bondReportCh := make(chan generalbondreport.GeneralBondReportPosition, workers)
	errCh := make(chan error, 1)
	var wg sync.WaitGroup

	// Создаю N Воркеров
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.bondWorker(ctx2, workList, bondReportCh, errCh, totalAmount, operationsByAssetUid)
		}()

	}
	// 1. Начинаем скармливать работу вокрерам
	go s.producePositions(ctx, portfolioPositions, workList)

	go func() {
		wg.Wait()
		close(bondReportCh)
	}()

loop:
	for {
		select {
		case bondReport, ok := <-bondReportCh:
			if !ok {
				break loop
			}
			s.addBondReport(&generalBondReports, bondReport)
		case er := <-errCh:
			cancel()
			return generalbondreport.GeneralBondReports{}, e.WrapIfErr("failed to process BondPosition", er)

		}
	}

	return generalBondReports, nil
}

func (s *Service) bondWorker(
	ctx context.Context,
	workList <-chan domain.PortfolioPositionsWithAssetUid,
	resultCh chan<- generalbondreport.GeneralBondReportPosition,
	errCh chan<- error,
	totalAmount float64,
	operationsByAssetUid map[string][]domain.OperationWithoutCustomTypes,
) {
	for position := range workList {
		bondReport, er := s.processBondPosition(ctx,
			position,
			totalAmount,
			operationsByAssetUid[position.AssetUid])
		if er != nil {
			if errors.Is(er, ErrEmptyBondPositions) {
				continue
			}
			select {
			case errCh <- er:
			case <-ctx.Done():
			}
			return
		}
		select {
		case resultCh <- bondReport:
		case <-ctx.Done():
			return
		}
	}
}

func (s *Service) producePositions(
	ctx context.Context,
	positions []domain.PortfolioPositionsWithAssetUid,
	out chan<- domain.PortfolioPositionsWithAssetUid,
) {
	defer close(out)

	for _, position := range positions {
		if position.InstrumentType != "bond" {
			continue
		}

		select {
		case out <- position:
		case <-ctx.Done():
			return
		}
	}
}

func (s *Service) processBondPosition(ctx context.Context,
	position domain.PortfolioPositionsWithAssetUid,
	totalAmount float64,
	operationsDb []domain.OperationWithoutCustomTypes,
) (_ generalbondreport.GeneralBondReportPosition, err error) {
	const op = "service.processBondPosition"
	defer logging.LogOperation_Debug(ctx, s.logger, op, &err)()

	select {
	case <-ctx.Done():
		return generalbondreport.GeneralBondReportPosition{}, ctx.Err()
	default:

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
