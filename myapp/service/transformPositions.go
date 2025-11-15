package service

import (
	"context"
	"errors"
	"time"

	"main.go/clients/tinkoffApi"
	"main.go/lib/e"
	"main.go/service/service_models"
)

const (
	hoursToUpdate = 12.0
)

func (c *Client) TransformPositions(accountID string, portffolioPositions []tinkoffApi.PortfolioPositions) (_ []service_models.PortfolioPosition, err error) {
	defer func() { err = e.WrapIfErr("transPositions err ", err) }()
	portfolio := make([]service_models.PortfolioPosition, 0)
	for _, v := range portffolioPositions {
		assetUid, err := c.GetUidByInstrUid(v.InstrumentUid)
		if err != nil {
			return portfolio, err
		}
		transPosionRet := service_models.PortfolioPosition{
			InstrumentType: v.InstrumentType,
		}
		transPosionRet.AssetUid = assetUid
		portfolio = append(portfolio, transPosionRet)
	}

	return portfolio, nil
}

func (c *Client) GetUidByInstrUid(instrumentUid string) (asset_uid string, err error) {
	defer func() { err = e.WrapIfErr("can't get uid", err) }()
	date, err := c.Storage.IsUpdatedUids(context.Background())
	if err != nil && !errors.Is(err, service_models.ErrEmptyUids) {
		return "", err
	}

	if time.Since(date).Hours() < hoursToUpdate {
		assetUid, err := c.Storage.GetUid(context.Background(), instrumentUid)
		if !errors.Is(err, service_models.ErrEmptyUids) {
			return "", err
		}
		return assetUid, nil
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
