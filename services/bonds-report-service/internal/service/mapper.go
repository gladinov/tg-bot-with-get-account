package service

import (
	"bonds-report-service/internal/models/domain"
	"context"
	"errors"

	"github.com/gladinov/e"
)

func (c *Client) MapPositionsToPositionsWithAssetUid(ctx context.Context, portffolioPositions []domain.PortfolioPosition) (_ []domain.PortfolioPositionsWithAssetUid, err error) {
	const op = "service.TransformPositions"

	portfolio := make([]domain.PortfolioPositionsWithAssetUid, 0, len(portffolioPositions))
	for _, v := range portffolioPositions {
		assetUid, err := c.UidProvider.GetUid(ctx, v.InstrumentUid)
		if err != nil && !errors.Is(err, domain.ErrEmptyUidAfterUpdate) {
			return nil, e.WrapIfErr("failed to get uid by instrument uid", err)
		}
		if errors.Is(err, domain.ErrEmptyUidAfterUpdate) {
			assetUid = ""
		}
		transformPosition := domain.PortfolioPositionsWithAssetUid{
			InstrumentType: v.InstrumentType,
			AssetUid:       assetUid,
		}
		portfolio = append(portfolio, transformPosition)
	}

	return portfolio, nil
}
