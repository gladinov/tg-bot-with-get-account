package service

import (
	"bonds-report-service/internal/models/domain"
	"context"
	"errors"

	"github.com/gladinov/e"
)

func (c *Client) TransformPositions(ctx context.Context, portffolioPositions []domain.PortfolioPosition) (_ []domain.PortfolioPositionWithoutCustomTypes, err error) {
	const op = "service.TransformPositions"

	portfolio := make([]domain.PortfolioPositionWithoutCustomTypes, 0, len(portffolioPositions))
	for _, v := range portffolioPositions {
		assetUid, err := c.GetUidByInstrUid(ctx, v.InstrumentUid)
		if err != nil && !errors.Is(err, domain.ErrEmptyUidAfterUpdate) {
			return nil, e.WrapIfErr("failed to get uid by instrument uid", err)
		}
		if errors.Is(err, domain.ErrEmptyUidAfterUpdate) {
			assetUid = ""
		}
		transformPosition := domain.PortfolioPositionWithoutCustomTypes{
			InstrumentType: v.InstrumentType,
			AssetUid:       assetUid,
		}
		portfolio = append(portfolio, transformPosition)
	}

	return portfolio, nil
}
