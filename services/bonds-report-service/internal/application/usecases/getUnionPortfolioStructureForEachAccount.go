package usecases

import (
	"bonds-report-service/internal/application/dto"
	unionportf "bonds-report-service/internal/application/helpers/unionPortf"
	"bonds-report-service/internal/application/presenter"
	"bonds-report-service/internal/domain"
	"bonds-report-service/internal/utils/logging"
	"context"

	"github.com/gladinov/e"
)

func (s *Service) GetUnionPortfolioStructureForEachAccount(ctx context.Context) (_ domain.UnionPortfolioStructureResponce, err error) {
	const op = "service.GetUnionPortfolioStructureForEachAccount"

	defer logging.LogOperation_Debug(ctx, s.logger, op, &err)()

	accounts, err := s.Helpers.TinkoffHelper.TinkoffGetAccounts(ctx)
	response := domain.UnionPortfolioStructureResponce{}
	if err != nil {
		return domain.UnionPortfolioStructureResponce{}, e.WrapIfErr("cant' get accounts from tinkoff", err)
	}
	unionPortfolioStructure, err := s.getUnionPortfolioStructure(ctx, accounts)
	if err != nil {
		return domain.UnionPortfolioStructureResponce{}, e.WrapIfErr("cant' get union portfolio structure", err)
	}
	response.Report = unionPortfolioStructure

	return response, nil
}

func (s *Service) getUnionPortfolioStructure(ctx context.Context, accounts map[string]domain.Account) (_ string, err error) {
	const op = "service.getUnionPortfolioStructure"

	defer logging.LogOperation_Debug(ctx, s.logger, op, &err)()

	positionsList := make([]*domain.PortfolioByTypeAndCurrency, 0)
	for _, account := range accounts {
		if account.Status != 2 {
			continue
		}
		portfolio, err := s.Helpers.TinkoffHelper.TinkoffGetPortfolio(ctx, account)
		if err != nil {
			return "", e.WrapIfErr("cant' get portfolio from Tinkoff", err)
		}
		positions := portfolio.Positions

		potfolioStructure, err := s.Helpers.DividerByAssetType.DivideByType(ctx, positions)
		if err != nil {
			return "", e.WrapIfErr("couldnot divide by type", err)
		}
		positionsList = append(positionsList, potfolioStructure)
	}

	unionPositions := unionportf.UnionPortf(positionsList)

	vizualizeUnionPositions := presenter.ResponsePortfolioStructure(ctx, s.logger, unionPositions, dto.UnionPortf, "")

	return vizualizeUnionPositions, nil
}
