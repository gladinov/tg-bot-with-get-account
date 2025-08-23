package tinkoffApi

import (
	"errors"

	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
	"main.go/lib/e"
)

func (c *Client) GetFutureBy(figi string) (*pb.Future, error) {
	instrumentService := c.Client.NewInstrumentsServiceClient()
	futuresResponse, err := instrumentService.FutureByFigi(figi)
	if err != nil {
		return nil, e.WrapIfErr("can't get futures by figi", err)
	}
	return futuresResponse.FutureResponse.Instrument, nil
}

func (c *Client) GetShareBy(figi string) (*pb.Share, error) {
	instrumentService := c.Client.NewInstrumentsServiceClient()
	shareResponse, err := instrumentService.ShareByFigi(figi)
	if err != nil {
		return nil, e.WrapIfErr("can't get share by figi", err)
	}
	return shareResponse.ShareResponse.Instrument, nil
}

func (c *Client) GetCurrencyBy(figi string) (*pb.Currency, error) {
	instrumentService := c.Client.NewInstrumentsServiceClient()
	currencyResponse, err := instrumentService.CurrencyByFigi(figi)
	if err != nil {
		return nil, e.WrapIfErr("can't get share by figi", err)
	}
	return currencyResponse.CurrencyResponse.Instrument, nil
}

func (c *Client) GetBaseShareFutureValute(positionUid string) (string, error) {
	instrumentService := c.Client.NewInstrumentsServiceClient()
	instrumentsShortResponce, err := instrumentService.FindInstrument(positionUid)
	if err != nil {
		return "", e.WrapIfErr("can't get base share future valute", err)
	}
	instrumentsShort := instrumentsShortResponce.Instruments
	if len(instrumentsShort) == 0 {
		return "", errors.New("can't get base share future valute")
	}
	instrument := instrumentsShort[0]
	if instrument.InstrumentType != "share" {
		return "", errors.New("instrument is not share")
	}
	share, err := c.GetShareBy(instrument.Figi)
	if err != nil {
		return "", e.WrapIfErr("can't get base share future valute", err)
	}

	return share.Currency, nil
}
