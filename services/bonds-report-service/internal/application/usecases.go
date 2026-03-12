package service

import (
	"bonds-report-service/internal/application/visualization"
	"bonds-report-service/internal/domain"
	"bonds-report-service/internal/domain/generalbondreport"
	"bonds-report-service/internal/domain/report"
	"bonds-report-service/internal/utils/logging"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gladinov/e"
)

func (s *Service) GetBondReportsByFifo(ctx context.Context, chatID int) (err error) {
	const op = "service.GetBondReportsByFifo"
	logg := slog.With(s.logger)
	defer logging.LogOperation_Debug(ctx, logg, op, &err)()

	accounts, err := s.Tinkoff.Portfolio.GetAccounts(ctx)
	if err != nil {
		return e.WrapIfErr("get accounts error", err)
	}

	for _, account := range accounts {
		err = s.updateOperations(ctx, chatID, account.ID, account.OpenedDate)
		if err != nil {
			return e.WrapIfErr("update operation error", err)
		}
		if account.Status != 2 {
			continue
		}
		portfolio, err := s.TinkoffGetPortfolio(ctx, account)
		if err != nil {
			return e.WrapIfErr("tinkoffGetPortfolio err", err)
		}

		portfolioPositions, err := s.MapPositionsToPositionsWithAssetUid(ctx, portfolio.Positions)
		if err != nil {
			return e.WrapIfErr("transformPositions err", err)
		}
		err = s.Storage.DeleteBondReport(ctx, chatID, account.ID)
		if err != nil {
			return e.WrapIfErr("deleteBondReport err", err)
		}
		bondsInRub := make([]report.BondReport, 0)

		for _, position := range portfolioPositions {
			if position.InstrumentType == "bond" {
				operationsDb, err := s.Storage.GetOperations(ctx, chatID, position.AssetUid, account.ID)
				if err != nil {
					return e.WrapIfErr("storage.GetOperation error", err)
				}

				reporLines, err := s.createNewReportLines(ctx, position, operationsDb)
				if err != nil {
					return e.WrapIfErr("failed to create new report lines", err)
				}

				resultBondPosition, err := s.ReportProcessor.ProcessOperations(ctx, reporLines)
				if err != nil {
					return e.WrapIfErr("failed to process operation", err)
				}

				firstBuyDate := resultBondPosition.CurrentPositions[0].BuyDate

				moexBuyDateData, err := s.GetSpecificationsFromMoex(ctx, reporLines.Bond.Ticker, firstBuyDate)
				if err != nil {
					return e.WrapIfErr("failed to get specifications from moex to buy date", err)
				}
				moexNowData, err := s.GetSpecificationsFromMoex(ctx, reporLines.Bond.Ticker, s.now())
				if err != nil {
					return e.WrapIfErr("failed to get specifications from moex to buy now", err)
				}

				bondReport, err := s.BondReportProcessor.CreateBondReport(
					ctx,
					resultBondPosition.CurrentPositions,
					moexBuyDateData,
					moexNowData,
				)
				if err != nil {
					return e.WrapIfErr("CreateBondReport error", err)
				}
				bondsInRub = append(bondsInRub, bondReport.BondsInRUB...)
			}
		}
		// TODO: Сдлеать это асинхронно, не дожидаясь завершения, выйти из функции
		err = s.Storage.SaveBondReport(ctx, chatID, account.ID, bondsInRub)
		if err != nil {
			return e.WrapIfErr("Storage.SaveBondReport error", err)
		}
	}
	return nil
}

func (s *Service) GetBondReports(ctx context.Context, chatID int) (_ BondReportsResponce, err error) {
	const op = "service.GetBondReports"
	defer logging.LogOperation_Debug(ctx, s.logger, op, &err)()

	reportsInByteByAccounts := make([][]*visualization.MediaGroup, 0)

	accounts, err := s.Tinkoff.Portfolio.GetAccounts(ctx)
	if err != nil {
		return BondReportsResponce{}, err
	}

	for _, account := range accounts {
		err = s.updateOperations(ctx, chatID, account.ID, account.OpenedDate)
		if err != nil {
			return BondReportsResponce{}, err
		}
		if account.Status != 2 {
			continue
		}
		portfolio, err := s.TinkoffGetPortfolio(ctx, account)
		if err != nil {
			return BondReportsResponce{}, err
		}

		portfolioPositions, err := s.MapPositionsToPositionsWithAssetUid(ctx, portfolio.Positions)
		if err != nil {
			return BondReportsResponce{}, err
		}

		err = s.Storage.DeleteGeneralBondReport(ctx, chatID, account.ID)
		if err != nil {
			return BondReportsResponce{}, err
		}
		generalBondReports := generalbondreport.GeneralBondReports{
			RubBondsReport:      make(map[generalbondreport.TickerTimeKey]generalbondreport.GeneralBondReportPosition),
			EuroBondsReport:     make(map[generalbondreport.TickerTimeKey]generalbondreport.GeneralBondReportPosition),
			ReplacedBondsReport: make(map[generalbondreport.TickerTimeKey]generalbondreport.GeneralBondReportPosition),
		}
		// Итерируемся по позициям портфеля
		// TODO: Возможно лучше итерироваться по указктелю
		for _, position := range portfolioPositions {
			// Проверяем, что бумага типа БОНД
			if position.InstrumentType == "bond" {
				// Получаем все операции по данной бумаге
				operationsDb, err := s.Storage.GetOperations(ctx, chatID, position.AssetUid, account.ID)
				if err != nil {
					return BondReportsResponce{}, e.WrapIfErr("failed to get operations from strorge", err)
				}
				// Создаем струкуру ReportLine, которая вбирает в себя:
				// список операций по бумаге в портфеле,
				// идендификаторф бумаги, полученные с ТинькоффАпи
				// последняя цена данной бумаги
				// Курс валюты
				reporLines, err := s.createNewReportLines(ctx, position, operationsDb)
				if err != nil {
					return BondReportsResponce{}, e.WrapIfErr("failed to create new report lines", err)
				}
				// Обрабатываем операции и получаем открытые позиции по данной бумаге
				resultBondPosition, err := s.ReportProcessor.ProcessOperations(ctx, reporLines)
				if err != nil {
					return BondReportsResponce{}, e.WrapIfErr("failed to process operation", err)
				}
				// Общая стоимость портфеля
				totalAmount := portfolio.TotalAmount.ToFloat()

				if len(resultBondPosition.CurrentPositions) == 0 {
					s.logger.WarnContext(
						ctx,
						"len of result bond positions is empty. But in portfolio request to TinkoffApi position in portfolio",
						slog.Any("postion", position))
					continue
				}

				firstBuyDate := resultBondPosition.CurrentPositions[0].BuyDate

				moexBuyDateData, err := s.GetSpecificationsFromMoex(ctx, reporLines.Bond.Ticker, firstBuyDate)
				if err != nil {
					return BondReportsResponce{}, err
				}
				moexNowData, err := s.GetSpecificationsFromMoex(ctx, reporLines.Bond.Ticker, s.now())
				if err != nil {
					return BondReportsResponce{}, err
				}

				bondReport, err := s.GeneralBondReportProcessor.GetGeneralBondReportPosition(
					ctx,
					resultBondPosition.CurrentPositions,
					totalAmount,
					moexBuyDateData,
					moexNowData, firstBuyDate)
				if err != nil {
					return BondReportsResponce{}, err
				}
				switch {
				case bondReport.Replaced:
					tickerTimeKey := generalbondreport.TickerTimeKey{
						Ticker: bondReport.Ticker,
						Time:   bondReport.BuyDate,
					}
					generalBondReports.ReplacedBondsReport[tickerTimeKey] = bondReport
				case bondReport.Currencies != "rub":
					tickerTimeKey := generalbondreport.TickerTimeKey{
						Ticker: bondReport.Ticker,
						Time:   bondReport.BuyDate,
					}
					generalBondReports.EuroBondsReport[tickerTimeKey] = bondReport
				default:
					tickerTimeKey := generalbondreport.TickerTimeKey{
						Ticker: bondReport.Ticker,
						Time:   bondReport.BuyDate,
					}
					generalBondReports.RubBondsReport[tickerTimeKey] = bondReport
				}

			}
		}

		reportsInByte, err := visualization.GenerateTablePNG(ctx,
			s.logger,
			s.prepareToGenerateTablePNG(ctx, &generalBondReports),
			chatID,
			account.ID)
		if err != nil {
			return BondReportsResponce{}, err
		}
		reportsInByteByAccounts = append(reportsInByteByAccounts, reportsInByte)

	}
	getBondReportsResponce := BondReportsResponce{Media: reportsInByteByAccounts}
	return getBondReportsResponce, nil
}

func (s *Service) GetAccountsList(ctx context.Context) (answ AccountListResponce, err error) {
	const op = "service.GetAccountsList"

	start := time.Now()
	logg := s.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("can't get accounts", err)
	}()

	var accStr string = "По данному аккаунту доступны следующие счета:"

	accs, err := s.Tinkoff.Portfolio.GetAccounts(ctx)
	if err != nil {
		return AccountListResponce{}, err
	}
	for _, account := range accs {
		accStr += fmt.Sprintf("\n ID:%s, Type: %s, Name: %s, Status: %v \n", account.ID, account.Type, account.Name, account.Status)
	}
	accountResponce := AccountListResponce{Accounts: accStr}
	return accountResponce, nil
}

func (s *Service) GetUsd(ctx context.Context) (_ domain.UsdResponce, err error) {
	const op = "service.GetUsd"

	start := time.Now()
	logg := s.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("usd get error", err)
	}()

	usd, err := s.GetCurrencyFromCB(ctx, "usd", time.Now())
	if err != nil {
		return domain.UsdResponce{}, err
	}
	usdResponce := domain.UsdResponce{Usd: usd}

	return usdResponce, nil
}

func (s *Service) GetPortfolioStructureForEachAccount(ctx context.Context) (_ domain.PortfolioStructureForEachAccountResponce, err error) {
	const op = "service.GetPortfolioStructureForEachAccount"

	start := time.Now()
	logg := s.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("cant' Get Portfolio Structure For Each Account", err)
	}()

	accounts, err := s.Tinkoff.Portfolio.GetAccounts(ctx)
	response := domain.PortfolioStructureForEachAccountResponce{}
	if err != nil {
		return domain.PortfolioStructureForEachAccountResponce{}, err
	}
	for _, account := range accounts {
		if account.Status == 3 {
			continue
		}
		report, err := s.getPortfolioStructure(ctx, account)
		if err != nil {
			return domain.PortfolioStructureForEachAccountResponce{}, err
		}
		response.PortfolioStructures = append(response.PortfolioStructures, report)
	}
	return response, nil
}

func (s *Service) GetUnionPortfolioStructureForEachAccount(ctx context.Context) (_ domain.UnionPortfolioStructureResponce, err error) {
	const op = "service.GetUnionPortfolioStructureForEachAccount"

	start := time.Now()
	logg := s.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("cant' Get Union Portfolio Structure For Each Account", err)
	}()

	accounts, err := s.Tinkoff.Portfolio.GetAccounts(ctx)
	response := domain.UnionPortfolioStructureResponce{}
	if err != nil {
		return domain.UnionPortfolioStructureResponce{}, err
	}
	unionPortfolioStructure, err := s.getUnionPortfolioStructure(ctx, accounts)
	if err != nil {
		return domain.UnionPortfolioStructureResponce{}, err
	}
	response.Report = unionPortfolioStructure

	return response, nil
}

func (s *Service) GetUnionPortfolioStructureWithSber(ctx context.Context) (_ domain.UnionPortfolioStructureWithSberResponce, err error) {
	const op = "service.GetUnionPortfolioStructureWithSber"

	start := time.Now()
	logg := s.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("can't get union portfolio structure with Sber", err)
	}()

	responce := domain.UnionPortfolioStructureWithSberResponce{}
	accounts, err := s.Tinkoff.Portfolio.GetAccounts(ctx)
	if err != nil {
		return domain.UnionPortfolioStructureWithSberResponce{}, err
	}
	positionsList := make([]*domain.PortfolioByTypeAndCurrency, 0)
	for _, account := range accounts {
		if account.Status != 2 {
			continue
		}
		portfolio, err := s.TinkoffGetPortfolio(ctx, account)
		if err != nil {
			return domain.UnionPortfolioStructureWithSberResponce{}, err
		}
		positions := portfolio.Positions

		potfolioStructure, err := s.divideByType(ctx, positions)
		if err != nil {
			return domain.UnionPortfolioStructureWithSberResponce{}, err
		}
		positionsList = append(positionsList, potfolioStructure)
	}

	sberPortfolio, err := s.divideByTypeFromSber(ctx, s.External.Sber.Portfolio)
	if err != nil {
		return domain.UnionPortfolioStructureWithSberResponce{}, err
	}

	positionsList = append(positionsList, sberPortfolio)

	accountTitle := "Струтура всех инвестиций\n"
	unionPositions, err := s.unionPortf(positionsList)
	if err != nil {
		return domain.UnionPortfolioStructureWithSberResponce{}, err
	}
	vizualizeUnionPositions := visualization.ResponsePortfolioStructure(ctx, s.logger, unionPositions)
	out := accountTitle + vizualizeUnionPositions
	responce.Report = out
	return responce, nil
}
