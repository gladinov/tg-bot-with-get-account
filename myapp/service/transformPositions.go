package service

import (
	"errors"
	"fmt"

	"main.go/clients/tinkoffApi"
	"main.go/service/service_models"
)

// Обрабатываем в нормальный формат портфеля
func (s *Client) TransPositions(account *tinkoffApi.Account, assetUidInstrumentUidMap map[string]string) (service_models.Portfolio, error) {
	Portfolio := service_models.Portfolio{}
	for _, v := range account.Portfolio {
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
			BondPosition.Identifiers.AssetUid = assetUidInstrumentUidMap[BondPosition.Identifiers.InstrumentUid]

			// Получаем Тикер, Режим торгов и Короткое имя инструмента
			// BondPosition.GetBondsActionsFromPortfolio(client)
			resFromTinkoff, err := s.Tinkoffapi.GetBondsActionsFromTinkoff(BondPosition.Identifiers.InstrumentUid)
			if err != nil {
				return Portfolio, errors.New("TransPositions:GetBondsActionsFromTinkoff" + err.Error())
			}
			BondPosition.Identifiers.Ticker = resFromTinkoff.Ticker
			BondPosition.Identifiers.ClassCode = resFromTinkoff.ClassCode
			BondPosition.Name = resFromTinkoff.Name

			//  Получение данных с московской биржи
			// BondPosition.GetActionFromMoex()

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
				AssetUid:                 assetUidInstrumentUidMap[v.GetInstrumentUid()],
				VarMargin:                v.GetVarMargin().ToFloat(),
				ExpectedYieldFifo:        v.GetExpectedYieldFifo().ToFloat(),
				DailyYield:               v.GetDailyYield().ToFloat(),
			}
			Portfolio.PortfolioPositions = append(Portfolio.PortfolioPositions, transPosionRet)
		}
	}
	fmt.Printf("✓ Добавлено %v позиций в Account.PortfolioPositions по счету %s\n", len(Portfolio.PortfolioPositions), account.Id)
	fmt.Printf("✓ Добавлено %v позиций в Account.PortfolioBondPositions по счету %s\n", len(Portfolio.BondPositions), account.Id)
	return Portfolio, nil
}
