package service

import (
	"context"
	"errors"

	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
	"main.go/lib/e"
	"main.go/service/service_models"
)

func (c *Client) TransformPositions(accountID string, portffolioPositions []*pb.PortfolioPosition) (_ []service_models.PortfolioPosition, err error) {
	defer func() { err = e.WrapIfErr("transPositions err ", err) }()
	portfolio := make([]service_models.PortfolioPosition, 0)
	for _, v := range portffolioPositions {
		assetUid, err := c.GetUidByInstrUid(v.GetInstrumentUid())
		if err != nil {
			return portfolio, err
		}
		transPosionRet := service_models.PortfolioPosition{
			InstrumentType: v.GetInstrumentType(),
		}
		transPosionRet.AssetUid = assetUid
		portfolio = append(portfolio, transPosionRet)
	}

	c.Tinkoffapi.Client.Logger.Infof("✓ Добавлено %v позиций в portfolio по счету %s\n", len(portfolio), accountID)
	// fmt.Printf("✓ Добавлено %v позиций в portfolio по счету %s\n", len(portfolio), accountID)
	return portfolio, nil
}

func (c *Client) GetUidByInstrUid(instrumentUid string) (asset_uid string, err error) {
	defer func() { err = e.WrapIfErr("can't get uid", err) }()
	exist, err := c.Storage.IsUpdatedUids(context.Background())
	if err != nil && !errors.Is(err, service_models.ErrEmptyUids) {
		return "", err
	}

	if exist {
		assetUid, err := c.Storage.GetUid(context.Background(), instrumentUid)
		if err == nil {
			return assetUid, nil
		}
		if !errors.Is(err, service_models.ErrEmptyUids) {
			return "", err
		}
	}

	assetUid, err := c.updateAndGetUid(instrumentUid)
	if err != nil {
		return "", err
	}

	return assetUid, nil

}

func (c *Client) updateAndGetUid(instrumentUid string) (asset_uid string, err error) {
	defer func() { err = e.WrapIfErr("can't update uids", err) }()
	allAssetUids, err := c.Tinkoffapi.GetAllAssetUids()
	if err != nil {
		return "", err
	}
	asset_uid, exist := allAssetUids[instrumentUid]
	if !exist {
		return "", err
	}
	err = c.Storage.SaveUids(context.Background(), allAssetUids)
	if err != nil {
		return "", err
	}
	return asset_uid, nil
}
