package usecases

import (
	"bonds-report-service/internal/application/dto"
	unionportf "bonds-report-service/internal/application/helpers/unionPortf"
	"bonds-report-service/internal/application/presenter"
	"bonds-report-service/internal/domain"
	"bonds-report-service/internal/utils/logging"
	"context"
	"errors"

	"github.com/gladinov/e"
)

var (
	ErrEmptyReport                = errors.New("no elements in report")
	ErrEmptyPosition              = errors.New("positions are empty")
	ErrpositionsClassCodeVariants = errors.New("positions class code variants are empty")
)

func (s *Service) GetUnionPortfolioStructureWithSber(ctx context.Context) (_ domain.UnionPortfolioStructureWithSberResponce, err error) {
	const op = "service.GetUnionPortfolioStructureWithSber"

	defer logging.LogOperation_Debug(ctx, s.logger, op, &err)()

	responce := domain.UnionPortfolioStructureWithSberResponce{}
	accounts, err := s.Helpers.TinkoffHelper.TinkoffGetAccounts(ctx)
	if err != nil {
		return domain.UnionPortfolioStructureWithSberResponce{}, e.WrapIfErr("cant' get accounts from tinkoff", err)
	}
	positionsList := make([]*domain.PortfolioByTypeAndCurrency, 0)
	for _, account := range accounts {
		if account.Status != 2 {
			continue
		}
		portfolio, err := s.Helpers.TinkoffHelper.TinkoffGetPortfolio(ctx, account)
		if err != nil {
			return domain.UnionPortfolioStructureWithSberResponce{}, e.WrapIfErr("cant' get portfolio from Tinkoff", err)
		}
		positions := portfolio.Positions

		potfolioStructure, err := s.Helpers.DividerByAssetType.DivideByType(ctx, positions)
		if err != nil {
			return domain.UnionPortfolioStructureWithSberResponce{}, e.WrapIfErr("couldnot divide by type", err)
		}
		positionsList = append(positionsList, potfolioStructure)
	}

	sberPortfolio, err := s.divideByTypeFromSber(ctx, s.External.Sber.Portfolio)
	if err != nil {
		return domain.UnionPortfolioStructureWithSberResponce{}, e.WrapIfErr("couldnot divide by type from sber", err)
	}

	positionsList = append(positionsList, sberPortfolio)

	unionPositions := unionportf.UnionPortf(positionsList)

	vizualizeUnionPositions := presenter.ResponsePortfolioStructure(ctx, s.logger, unionPositions, dto.UnionPortfWithSber, "")

	responce.Report = vizualizeUnionPositions
	return responce, nil
}

func (s *Service) divideByTypeFromSber(ctx context.Context, positions map[string]float64) (_ *domain.PortfolioByTypeAndCurrency, err error) {
	const op = "service.DivideByTypeFromSber"

	defer logging.LogOperation_Debug(ctx, s.logger, op, &err)()

	portfolio := domain.NewPortfolioByTypeAndCurrency()

	if len(positions) == 0 {
		return portfolio, ErrEmptyPosition
	}
	for ticker, quantity := range positions {
		positionsClassCodeVariants, err := s.Helpers.TinkoffHelper.TinkoffFindBy(ctx, ticker)
		if err != nil {
			return nil, e.WrapIfErr("can't find by ticker from tinkoff", err)
		}
		if len(positionsClassCodeVariants) == 0 {
			return nil, ErrpositionsClassCodeVariants
		}

		switch positionsClassCodeVariants[0].InstrumentType {
		case bond:
			bondUid := positionsClassCodeVariants[0].Uid
			bond, err := s.Helpers.TinkoffHelper.TinkoffGetBondByUid(ctx, bondUid)
			if err != nil {
				return nil, e.WrapIfErr("can't get bond by uid from tinkoff", err)
			}
			currentNkd := bond.AciValue.ToFloat()
			currency := bond.Currency
			resp, err := s.Helpers.TinkoffHelper.TinkoffGetLastPriceInPersentageToNominal(ctx, bondUid)
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
