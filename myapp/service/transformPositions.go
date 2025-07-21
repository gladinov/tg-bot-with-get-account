package service

import (
	"context"
	"errors"
	"fmt"

	"main.go/clients/tinkoffApi"
	"main.go/lib/e"
	"main.go/service/service_models"
)

// Обрабатываем в нормальный формат портфеля
func (c *Client) TransPositions(account *tinkoffApi.Account) (portf service_models.Portfolio, err error) {
	defer func() { err = e.WrapIfErr("TransPositions err ", err) }()
	Portfolio := service_models.Portfolio{}
	for _, v := range account.Portfolio {
		assetUid, err := c.GetUidByInstrUid(v.GetInstrumentUid())
		if err != nil {
			return Portfolio, err
		}
		if v.InstrumentType == "bond" {
			BondPosition := service_models.Bond{
				Identifiers: service_models.Identifiers{
					Figi:          v.GetFigi(),
					InstrumentUid: v.GetInstrumentUid(),
					PositionUid:   v.GetPositionUid(),
				},
				InstrumentType:           v.GetInstrumentType(),
				Currency:                 v.GetAveragePositionPrice().Currency,
				Quantity:                 v.GetQuantity().ToFloat(),
				AveragePositionPrice:     v.GetAveragePositionPrice().ToFloat(),
				ExpectedYield:            v.GetExpectedYield().ToFloat(),
				CurrentNkd:               v.GetCurrentNkd().ToFloat(),
				CurrentPrice:             v.GetCurrentPrice().ToFloat(),
				AveragePositionPriceFifo: v.GetAveragePositionPriceFifo().ToFloat(),
				Blocked:                  v.GetBlocked(),
				ExpectedYieldFifo:        v.GetExpectedYieldFifo().ToFloat(),
				DailyYield:               v.GetDailyYield().ToFloat(),
			}
			// Получаем AssetUid с помощью МАПЫ assetUidInstrumentUidMap
			BondPosition.Identifiers.AssetUid = assetUid

			// Получаем Тикер, Режим торгов и Короткое имя инструмента
			// BondPosition.GetBondsActionsFromPortfolio(client)
			resFromTinkoff, err := c.Tinkoffapi.GetBondsActionsFromTinkoff(BondPosition.Identifiers.InstrumentUid)
			if err != nil {
				return Portfolio, err
			}
			BondPosition.Identifiers.Ticker = resFromTinkoff.Ticker
			BondPosition.Identifiers.ClassCode = resFromTinkoff.ClassCode
			BondPosition.Name = resFromTinkoff.Name

			// Добавляем позицию в срез позиций
			Portfolio.BondPositions = append(Portfolio.BondPositions, BondPosition)
		} else {
			transPosionRet := service_models.PortfolioPosition{
				Figi:                     v.GetFigi(),
				InstrumentType:           v.GetInstrumentType(),
				Currency:                 v.GetAveragePositionPrice().Currency,
				Quantity:                 v.GetQuantity().ToFloat(),
				AveragePositionPrice:     v.GetAveragePositionPrice().ToFloat(),
				ExpectedYield:            v.GetExpectedYield().ToFloat(),
				CurrentNkd:               v.GetCurrentNkd().ToFloat(),
				CurrentPrice:             v.GetCurrentPrice().ToFloat(),
				AveragePositionPriceFifo: v.GetAveragePositionPriceFifo().ToFloat(),
				Blocked:                  v.GetBlocked(),
				BlockedLots:              v.GetBlockedLots().ToFloat(),
				PositionUid:              v.GetPositionUid(),
				InstrumentUid:            v.GetInstrumentUid(),
				VarMargin:                v.GetVarMargin().ToFloat(),
				ExpectedYieldFifo:        v.GetExpectedYieldFifo().ToFloat(),
				DailyYield:               v.GetDailyYield().ToFloat(),
			}
			transPosionRet.AssetUid = assetUid
			Portfolio.PortfolioPositions = append(Portfolio.PortfolioPositions, transPosionRet)
		}
	}
	fmt.Printf("✓ Добавлено %v позиций в Account.PortfolioPositions по счету %s\n", len(Portfolio.PortfolioPositions), account.Id)
	fmt.Printf("✓ Добавлено %v позиций в Account.PortfolioBondPositions по счету %s\n", len(Portfolio.BondPositions), account.Id)
	return Portfolio, nil
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
