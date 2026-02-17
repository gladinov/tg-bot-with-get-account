package handlers

import (
	"bonds-report-service/internal/domain"
	httpmodels "bonds-report-service/internal/handlers/http"
)

func MapMediaGroupToHTTP(mg *domain.MediaGroup) *httpmodels.MediaGroup {
	if mg == nil {
		return nil
	}

	httpMG := httpmodels.NewMediaGroup()
	for _, r := range mg.Reports {
		httpMG.Reports = append(httpMG.Reports, MapImageDataToHTTP(r))
	}
	return httpMG
}

func MapImageDataToHTTP(img *domain.ImageData) *httpmodels.ImageData {
	if img == nil {
		return nil
	}
	return &httpmodels.ImageData{
		Name:    img.Name,
		Data:    img.Data,
		Caption: img.Caption,
	}
}

func MapAccountListToHTTP(acc *domain.AccountListResponce) *httpmodels.AccountListResponce {
	if acc == nil {
		return nil
	}
	return &httpmodels.AccountListResponce{
		Accounts: acc.Accounts,
	}
}

func MapBondReportsToHTTP(br *domain.BondReportsResponce) *httpmodels.BondReportsResponce {
	if br == nil {
		return nil
	}

	httpBR := &httpmodels.BondReportsResponce{
		Media: make([][]*httpmodels.MediaGroup, len(br.Media)),
	}

	for i, row := range br.Media {
		httpBR.Media[i] = make([]*httpmodels.MediaGroup, len(row))
		for j, mg := range row {
			httpBR.Media[i][j] = MapMediaGroupToHTTP(mg)
		}
	}

	return httpBR
}

func MapPortfolioStructureForEachAccountToHTTP(pf *domain.PortfolioStructureForEachAccountResponce) *httpmodels.PortfolioStructureForEachAccountResponce {
	if pf == nil {
		return nil
	}
	return &httpmodels.PortfolioStructureForEachAccountResponce{
		PortfolioStructures: pf.PortfolioStructures,
	}
}

func MapUnionPortfolioStructureToHTTP(u *domain.UnionPortfolioStructureResponce) *httpmodels.UnionPortfolioStructureResponce {
	if u == nil {
		return nil
	}
	return &httpmodels.UnionPortfolioStructureResponce{
		Report: u.Report,
	}
}

func MapUnionPortfolioStructureWithSberToHTTP(u *domain.UnionPortfolioStructureWithSberResponce) *httpmodels.UnionPortfolioStructureWithSberResponce {
	if u == nil {
		return nil
	}
	return &httpmodels.UnionPortfolioStructureWithSberResponce{
		Report: u.Report,
	}
}

func MapUsdToHTTP(u *domain.UsdResponce) *httpmodels.UsdResponce {
	if u == nil {
		return nil
	}
	return &httpmodels.UsdResponce{
		Usd: u.Usd,
	}
}
