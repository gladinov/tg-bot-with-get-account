package usecases

import (
	"bonds-report-service/internal/application/dto"
	"bonds-report-service/internal/application/presenter"
	"bonds-report-service/internal/domain"
	"bonds-report-service/internal/utils/logging"
	"context"

	"github.com/gladinov/e"
)

func (s *Service) GetPortfolioStructureForEachAccount(ctx context.Context) (_ domain.PortfolioStructureForEachAccountResponce, err error) {
	const op = "service.GetPortfolioStructureForEachAccount"

	defer logging.LogOperation_Debug(ctx, s.logger, op, &err)()

	accounts, err := s.Helpers.TinkoffHelper.TinkoffGetAccounts(ctx)
	response := domain.PortfolioStructureForEachAccountResponce{}
	if err != nil {
		return domain.PortfolioStructureForEachAccountResponce{}, e.WrapIfErr("cant' get accounts from tinkoff", err)
	}
	for _, account := range accounts {
		if account.Status == 3 {
			continue
		}
		report, err := s.getPortfolioStructure(ctx, account)
		if err != nil {
			return domain.PortfolioStructureForEachAccountResponce{}, e.WrapIfErr("cant' get portfolio structure", err)
		}
		response.PortfolioStructures = append(response.PortfolioStructures, report)
	}
	return response, nil
}

func (s *Service) getPortfolioStructure(ctx context.Context, account domain.Account) (_ string, err error) {
	const op = "service.getPortfolioStructure"

	defer logging.LogOperation_Debug(ctx, s.logger, op, &err)()

	portfolio, err := s.Helpers.TinkoffHelper.TinkoffGetPortfolio(ctx, account)
	if err != nil {
		return "", e.WrapIfErr("cant' get portfolio from Tinkoff", err)
	}
	positions := portfolio.Positions

	potfolioStructure, err := s.Helpers.DividerByAssetType.DivideByType(ctx, positions)
	if err != nil {
		return "", e.WrapIfErr("couldnot divide by type", err)
	}
	response := presenter.ResponsePortfolioStructure(ctx, s.logger, potfolioStructure, dto.EachPortf, account.Name)

	return response, nil
}
