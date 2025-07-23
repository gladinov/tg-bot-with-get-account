package tinkoffApi

import (
	"errors"
)

type BondIdentIdentifiers struct {
	Ticker          string
	ClassCode       string
	Name            string
	Nominal         float64
	NominalCurrency string
	Replaced        bool
}

func (c *Client) GetBondsActionsFromTinkoff(instrumentUid string) (BondIdentIdentifiers, error) {
	var res BondIdentIdentifiers
	instrumentService := c.Client.NewInstrumentsServiceClient()
	bondUid, err := instrumentService.BondByUid(instrumentUid)
	if err != nil {
		return res, errors.New("GetTickerFromUid: instrumentService.BondByUid" + err.Error())
	}
	res.Ticker = bondUid.BondResponse.Instrument.GetTicker()
	res.ClassCode = bondUid.BondResponse.Instrument.GetClassCode()
	res.Name = bondUid.BondResponse.Instrument.GetName()

	if bondUid.BondResponse.Instrument.GetBondType() == 1 {
		res.Replaced = true
	}
	res.Nominal = bondUid.BondResponse.Instrument.GetNominal().ToFloat()
	res.NominalCurrency = bondUid.Instrument.GetNominal().Currency
	return res, nil
}

func (c *Client) GetLastPriceFromTinkoffInPersentageToNominal(instrumentUid string) (float64, error) {
	marketDataClient := c.Client.NewMarketDataServiceClient()
	lastPriceAnswer, err := marketDataClient.GetLastPrices([]string{instrumentUid})
	if err != nil {
		return 0, errors.New("tinkoffApi:GetLastPriceFromTinkoff" + err.Error())
	}
	if len(lastPriceAnswer.LastPrices) == 0 {
		return 0, errors.New("tinkoffApi:GetLastPriceFromTinkoff: no price data")
	}

	lastPrice := lastPriceAnswer.LastPrices[0].Price.ToFloat()

	return lastPrice, nil
}
