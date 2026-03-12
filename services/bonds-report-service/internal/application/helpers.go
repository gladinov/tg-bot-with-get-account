package service

import (
	"bonds-report-service/internal/application/visualization"
	"bonds-report-service/internal/domain"
	"bonds-report-service/internal/domain/generalbondreport"
	"bonds-report-service/internal/domain/mapper"
	"bonds-report-service/internal/utils/logging"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"sync"
	"time"

	"github.com/gladinov/e"
)

var (
	ErrEmptyReport                = errors.New("no elements in report")
	ErrEmptyPosition              = errors.New("positions are empty")
	ErrpositionsClassCodeVariants = errors.New("positions class code variants are empty")
)

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

func (s *Service) divideByType(ctx context.Context, positions []domain.PortfolioPosition) (_ *domain.PortfolioByTypeAndCurrency, err error) {
	const op = "service.DivideByType"

	defer logging.LogOperation_Debug(ctx, s.logger, op, &err)()

	portfolio := domain.NewPortfolioByTypeAndCurrency()
	date := s.now()

	if len(positions) == 0 {
		return portfolio, ErrEmptyPosition
	}

	for _, pos := range positions {
		var positionPrice float64
		currencyOfPos := pos.CurrentPrice.Currency

		var vunit_rate float64
		if currencyOfPos != futuresPt && currencyOfPos != rub {
			vunit_rate, err = s.GetCurrencyFromCB(ctx, currencyOfPos, date)
			if err != nil {
				return nil, e.WrapIfErr("can't get currency from CB", err)
			}
		} else {
			vunit_rate = 1.0
		}
		positionPrice = pos.Quantity.ToFloat() * pos.CurrentPrice.ToFloat() * vunit_rate

		switch pos.InstrumentType {
		case bond:
			positionPrice += pos.CurrentNkd.ToFloat() * pos.Quantity.ToFloat() * vunit_rate
			portfolio.BondsAssets.SumOfAssets += positionPrice
			if _, exist := portfolio.BondsAssets.AssetsByCurrency[currencyOfPos]; !exist {
				portfolio.BondsAssets.AssetsByCurrency[currencyOfPos] = domain.NewAssetsByParam()
			}
			portfolio.BondsAssets.AssetsByCurrency[currencyOfPos].SumOfAssets += positionPrice
		case share:
			portfolio.SharesAssets.SumOfAssets += positionPrice
			if _, exist := portfolio.SharesAssets.AssetsByCurrency[currencyOfPos]; !exist {
				portfolio.SharesAssets.AssetsByCurrency[currencyOfPos] = domain.NewAssetsByParam()
			}
			portfolio.SharesAssets.AssetsByCurrency[currencyOfPos].SumOfAssets += positionPrice
		case futures:
			futures, err := s.TinkoffGetFutureBy(ctx, pos.Figi)
			if err != nil {
				return nil, e.WrapIfErr("can't get future data", err)
			}

			positionPrice = positionPrice / futures.MinPriceIncrement.ToFloat() * futures.MinPriceIncrementAmount.ToFloat()
			portfolio.FuturesAssets.SumOfAssets += positionPrice

			futureType := futures.AssetType
			switch futureType {
			case commodityType:
				if _, exist := portfolio.FuturesAssets.AssetsByType.Commodity.AssetsByCurrency[futures.Name]; !exist {
					portfolio.FuturesAssets.AssetsByType.Commodity.AssetsByCurrency[futures.Name] = domain.NewAssetsByParam()
				}
				portfolio.FuturesAssets.AssetsByType.Commodity.SumOfAssets += positionPrice
				portfolio.FuturesAssets.AssetsByType.Commodity.AssetsByCurrency[futures.Name].SumOfAssets += positionPrice
			case currencyType:
				if _, exist := portfolio.FuturesAssets.AssetsByType.Currency.AssetsByCurrency[futures.Name]; !exist {
					portfolio.FuturesAssets.AssetsByType.Currency.AssetsByCurrency[futures.Name] = domain.NewAssetsByParam()
				}
				portfolio.FuturesAssets.AssetsByType.Currency.SumOfAssets += positionPrice
				portfolio.FuturesAssets.AssetsByType.Currency.AssetsByCurrency[futures.Name].SumOfAssets += positionPrice
			case securityType:
				resp, err := s.TinkoffGetBaseShareFutureValute(ctx, futures.BasicAssetPositionUid)
				if err != nil {
					return nil, e.WrapIfErr("can't dget base share future valute from tinkoff", err)
				}
				valute := resp.Currency
				if _, exist := portfolio.FuturesAssets.AssetsByType.Security.AssetsByCurrency[valute]; !exist {
					portfolio.FuturesAssets.AssetsByType.Security.AssetsByCurrency[valute] = domain.NewAssetsByParam()
				}
				portfolio.FuturesAssets.AssetsByType.Security.SumOfAssets += positionPrice
				portfolio.FuturesAssets.AssetsByType.Security.AssetsByCurrency[valute].SumOfAssets += positionPrice
			case indexType:
				if _, exist := portfolio.FuturesAssets.AssetsByType.Index.AssetsByCurrency[futures.Name]; !exist {
					portfolio.FuturesAssets.AssetsByType.Index.AssetsByCurrency[futures.Name] = domain.NewAssetsByParam()
				}

				portfolio.FuturesAssets.AssetsByType.Index.SumOfAssets += positionPrice
				portfolio.FuturesAssets.AssetsByType.Index.AssetsByCurrency[futures.Name].SumOfAssets += positionPrice

			}
			// Чтобы сумма фьюча не сумировалась с суммой всех активов, так как фактически я за тело фьючерса не заплатил
			positionPrice = 0

		case etf:
			portfolio.EtfsAssets.SumOfAssets += positionPrice

			if _, exist := portfolio.EtfsAssets.AssetsByCurrency[currencyOfPos]; !exist {
				portfolio.EtfsAssets.AssetsByCurrency[currencyOfPos] = domain.NewAssetsByParam()
			}
			portfolio.EtfsAssets.AssetsByCurrency[currencyOfPos].SumOfAssets += positionPrice
		case currency:
			curr, err := s.TinkoffGetCurrencyBy(ctx, pos.Figi)
			if err != nil {
				return nil, e.WrapIfErr("can't get currency by figi from tinkoff", err)
			}
			currName := curr.Isin
			portfolio.CurrenciesAssets.SumOfAssets += positionPrice

			if _, exist := portfolio.CurrenciesAssets.AssetsByCurrency[currName]; !exist {
				portfolio.CurrenciesAssets.AssetsByCurrency[currName] = domain.NewAssetsByParam()
			}
			portfolio.CurrenciesAssets.AssetsByCurrency[currName].SumOfAssets += positionPrice

		default:
		}
		portfolio.AllAssets += positionPrice
	}

	return portfolio, nil
}

func (s *Service) divideByTypeFromSber(ctx context.Context, positions map[string]float64) (_ *domain.PortfolioByTypeAndCurrency, err error) {
	const op = "service.DivideByTypeFromSber"

	defer logging.LogOperation_Debug(ctx, s.logger, op, &err)()

	portfolio := domain.NewPortfolioByTypeAndCurrency()

	if len(positions) == 0 {
		return portfolio, ErrEmptyPosition
	}
	for ticker, quantity := range positions {
		positionsClassCodeVariants, err := s.TinkoffFindBy(ctx, ticker)
		if err != nil {
			return nil, e.WrapIfErr("can't find by ticker from tinkoff", err)
		}
		if len(positionsClassCodeVariants) == 0 {
			return nil, ErrpositionsClassCodeVariants
		}

		switch positionsClassCodeVariants[0].InstrumentType {
		case bond:
			bondUid := positionsClassCodeVariants[0].Uid
			bond, err := s.TinkoffGetBondByUid(ctx, bondUid)
			if err != nil {
				return nil, e.WrapIfErr("can't get bond by uid from tinkoff", err)
			}
			currentNkd := bond.AciValue.ToFloat()
			currency := bond.Currency
			resp, err := s.TinkoffGetLastPriceInPersentageToNominal(ctx, bondUid)
			if err != nil {
				return nil, e.WrapIfErr("can't get last price in persentage to nominal from tinkoff", err)
			}
			currentPriceInPers := resp.LastPrice.ToFloat()
			currentPrice := currentPriceInPers / 100 * bond.Nominal.ToFloat()
			currentNkdOfPosition := currentNkd * quantity
			positionPrice := currentPrice*quantity + currentNkdOfPosition

			portfolio.AllAssets += positionPrice
			portfolio.BondsAssets.SumOfAssets += positionPrice

			if existing, exist := portfolio.BondsAssets.AssetsByCurrency[currency]; !exist {
				portfolio.BondsAssets.AssetsByCurrency[currency] = &domain.AssetByParam{
					SumOfAssets: positionPrice,
				}
			} else {
				existing.SumOfAssets += positionPrice
			}

		case share:
		case futures:
		case etf:
		case currency:
		default:
		}

	}
	return portfolio, nil
}

func (s *Service) createNewReportLines(
	ctx context.Context,
	position domain.PortfolioPositionsWithAssetUid,
	operationsDb []domain.OperationWithoutCustomTypes,
) (_ *domain.ReportLine, err error) {
	const op = "service.CreateNewReportLines"
	defer logging.LogOperation_Debug(ctx, s.logger, op, &err)()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	var lastPrice domain.LastPrice
	var bondActions domain.BondIdentIdentifiers
	var vunitRate domain.Rate

	wg.Add(1)
	go func() {
		defer wg.Done()
		price, e := s.TinkoffGetLastPriceInPersentageToNominal(ctx, position.InstrumentUid)
		if e != nil {
			select {
			case errCh <- e:
				cancel()
			default:
			}
			return
		}

		lastPrice = price
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		actions, e := s.TinkoffGetBondActions(ctx, position.InstrumentUid)
		if e != nil {
			select {
			case errCh <- e:
				cancel()
			default:
			}
			return
		}

		bondActions = actions
	}()

	wg.Wait()

	select {
	case err := <-errCh:
		return nil, err
	default:
	}

	vunitRate, err = s.buildVunitRate(ctx, bondActions)
	if err != nil {
		return nil, err
	}

	reportLine := domain.NewReportLine(
		operationsDb,
		bondActions,
		lastPrice,
		vunitRate,
	)

	return &reportLine, nil
}

func (s *Service) buildVunitRate(ctx context.Context, bondActions domain.BondIdentIdentifiers) (domain.Rate, error) {
	// TODO: Не очевидная логика! Потом есть проверка,
	if !bondActions.Replaced {
		return domain.Rate{}, nil
	}

	isoCurrName := bondActions.NominalCurrency

	rate, err := s.GetCurrencyFromCB(ctx, isoCurrName, s.now())
	if err != nil {
		return domain.Rate{}, e.WrapIfErr("failed to get currency rate", err)
	}

	return domain.Rate{
		IsoCurrencyName: isoCurrName,
		Vunit_Rate:      domain.NewNullFloat64(rate, true, false),
	}, nil
}

func (s *Service) unionPortf(portfolios []*domain.PortfolioByTypeAndCurrency) (_ *domain.PortfolioByTypeAndCurrency, err error) {
	const op = "service.UnionPortf"

	start := time.Now()
	logg := s.logger.With(
		slog.String("op", op))
	logg.Info(fmt.Sprintf("start %s", op))
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("can't union portfolios ", err)
	}()

	unionPortf := domain.NewPortfolioByTypeAndCurrency()
	for _, portf := range portfolios {
		unionPortf.AllAssets += portf.AllAssets

		unionPortf.BondsAssets.SumOfAssets += portf.BondsAssets.SumOfAssets
		for k, v := range portf.BondsAssets.AssetsByCurrency {
			if existing, exist := unionPortf.BondsAssets.AssetsByCurrency[k]; !exist {
				unionPortf.BondsAssets.AssetsByCurrency[k] = domain.NewAssetsByParam()
				unionPortf.BondsAssets.AssetsByCurrency[k] = v
			} else {
				existing.SumOfAssets += v.SumOfAssets
			}
		}

		unionPortf.SharesAssets.SumOfAssets += portf.SharesAssets.SumOfAssets
		for currency, asset := range portf.SharesAssets.AssetsByCurrency {
			if existing, exists := unionPortf.SharesAssets.AssetsByCurrency[currency]; !exists {
				unionPortf.SharesAssets.AssetsByCurrency[currency] = &domain.AssetByParam{
					SumOfAssets: asset.SumOfAssets,
				}
			} else {
				existing.SumOfAssets += asset.SumOfAssets
			}
		}

		unionPortf.EtfsAssets.SumOfAssets += portf.EtfsAssets.SumOfAssets
		for currency, asset := range portf.EtfsAssets.AssetsByCurrency {
			if existing, exists := unionPortf.EtfsAssets.AssetsByCurrency[currency]; !exists {
				unionPortf.EtfsAssets.AssetsByCurrency[currency] = &domain.AssetByParam{
					SumOfAssets: asset.SumOfAssets,
				}
			} else {
				existing.SumOfAssets += asset.SumOfAssets
			}
		}

		unionPortf.CurrenciesAssets.SumOfAssets += portf.CurrenciesAssets.SumOfAssets
		for currency, asset := range portf.CurrenciesAssets.AssetsByCurrency {
			if existing, exists := unionPortf.CurrenciesAssets.AssetsByCurrency[currency]; !exists {
				unionPortf.CurrenciesAssets.AssetsByCurrency[currency] = &domain.AssetByParam{
					SumOfAssets: asset.SumOfAssets,
				}
			} else {
				existing.SumOfAssets += asset.SumOfAssets
			}
		}

		unionPortf.FuturesAssets.SumOfAssets += portf.FuturesAssets.SumOfAssets
		unionPortf.FuturesAssets.AssetsByType.Commodity.SumOfAssets += portf.FuturesAssets.AssetsByType.Commodity.SumOfAssets
		for currency, asset := range portf.FuturesAssets.AssetsByType.Commodity.AssetsByCurrency {
			if existing, exist := unionPortf.FuturesAssets.AssetsByType.Commodity.AssetsByCurrency[currency]; !exist {
				unionPortf.FuturesAssets.AssetsByType.Commodity.AssetsByCurrency[currency] = &domain.AssetByParam{
					SumOfAssets: asset.SumOfAssets,
				}
			} else {
				existing.SumOfAssets += asset.SumOfAssets
			}
		}

		unionPortf.FuturesAssets.AssetsByType.Currency.SumOfAssets += portf.FuturesAssets.AssetsByType.Currency.SumOfAssets

		for currency, asset := range portf.FuturesAssets.AssetsByType.Currency.AssetsByCurrency {
			if existing, exists := unionPortf.FuturesAssets.AssetsByType.Currency.AssetsByCurrency[currency]; !exists {
				unionPortf.FuturesAssets.AssetsByType.Currency.AssetsByCurrency[currency] = &domain.AssetByParam{
					SumOfAssets: asset.SumOfAssets,
				}
			} else {
				existing.SumOfAssets += asset.SumOfAssets
			}
		}

		unionPortf.FuturesAssets.AssetsByType.Security.SumOfAssets += portf.FuturesAssets.AssetsByType.Security.SumOfAssets

		for currency, asset := range portf.FuturesAssets.AssetsByType.Security.AssetsByCurrency {
			if existing, exists := unionPortf.FuturesAssets.AssetsByType.Security.AssetsByCurrency[currency]; !exists {
				unionPortf.FuturesAssets.AssetsByType.Security.AssetsByCurrency[currency] = &domain.AssetByParam{
					SumOfAssets: asset.SumOfAssets,
				}
			} else {
				existing.SumOfAssets += asset.SumOfAssets
			}
		}

		unionPortf.FuturesAssets.AssetsByType.Index.SumOfAssets += portf.FuturesAssets.AssetsByType.Index.SumOfAssets

		for currency, asset := range portf.FuturesAssets.AssetsByType.Index.AssetsByCurrency {
			if existing, exists := unionPortf.FuturesAssets.AssetsByType.Index.AssetsByCurrency[currency]; !exists {
				unionPortf.FuturesAssets.AssetsByType.Index.AssetsByCurrency[currency] = &domain.AssetByParam{
					SumOfAssets: asset.SumOfAssets,
				}
			} else {
				existing.SumOfAssets += asset.SumOfAssets
			}
		}
	}
	return unionPortf, nil
}

func (s *Service) getUnionPortfolioStructure(ctx context.Context, accounts map[string]domain.Account) (_ string, err error) {
	const op = "service.getUnionPortfolioStructure"

	start := time.Now()
	logg := s.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("can't get union portfolio structuret", err)
	}()

	positionsList := make([]*domain.PortfolioByTypeAndCurrency, 0)
	for _, account := range accounts {
		if account.Status != 2 {
			continue
		}
		portfolio, err := s.TinkoffGetPortfolio(ctx, account)
		if err != nil {
			return "", err
		}
		positions := portfolio.Positions

		potfolioStructure, err := s.divideByType(ctx, positions)
		if err != nil {
			return "", err
		}
		positionsList = append(positionsList, potfolioStructure)
	}
	accountTitle := "Струтура всех открытых счетов в Тинькофф Инвестициях\n"
	unionPositions, err := s.unionPortf(positionsList)
	if err != nil {
		return "", err
	}
	vizualizeUnionPositions := visualization.ResponsePortfolioStructure(ctx, s.logger, unionPositions)
	out := accountTitle + vizualizeUnionPositions
	return out, nil
}

func (s *Service) getPortfolioStructure(ctx context.Context, account domain.Account) (_ string, err error) {
	const op = "service.getPortfolioStructure"

	start := time.Now()
	logg := s.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("cant' get portfolio structure", err)
	}()

	portfolio, err := s.TinkoffGetPortfolio(ctx, account)
	if err != nil {
		return "", err
	}
	positions := portfolio.Positions

	accountTitle := fmt.Sprintf("Струтура брокерского счета: %s\n", account.Name)
	potfolioStructure, err := s.divideByType(ctx, positions)
	if err != nil {
		return "", err
	}
	respPotfolioStructure := visualization.ResponsePortfolioStructure(ctx, s.logger, potfolioStructure)

	response := accountTitle + respPotfolioStructure
	return response, nil
}

func (s *Service) updateOperations(ctx context.Context, chatID int, accountID string, openDate time.Time) (err error) {
	const op = "service.updateOperations"

	start := time.Now()
	logg := s.logger.With(
		slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("fineshed",
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		err = e.WrapIfErr("can't updateOperations", err)
	}()

	fromDate, err := s.Storage.LastOperationTime(ctx, chatID, accountID)
	// TODO: Если fromDate будет больше time.Now, то будет ошибка.
	// Есть вероятность такой ошибки, поэтому при тестировании функции нужно придумать другой способ вызова функции по последней операции
	fromDate = fromDate.Add(time.Microsecond * 1)

	if err != nil {
		if errors.Is(err, domain.ErrNoOpperations) {
			fromDate = openDate
		} else {
			return err
		}
	}

	opRequest := domain.NewOperationsRequest(accountID, fromDate)

	tinkoffOperations, err := s.TinkoffGetOperations(ctx, opRequest)
	if err != nil {
		return err
	}

	operations := mapper.MapOperationToOperationWithoutCustomTypes(tinkoffOperations)

	err = s.Storage.SaveOperations(ctx, chatID, accountID, operations)
	if err != nil {
		return err
	}
	return nil
}
