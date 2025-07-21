package tinkoffApi

import (
	"errors"
)

func (c *Client) GetAllAssetUids() (map[string]string, error) {
	instrumentService := c.Client.NewInstrumentsServiceClient()
	AssetsResponse, err := instrumentService.GetAssets()
	if err != nil {
		return nil, errors.New("GetAllAssetUids: instrumentService.GetAssets" + err.Error())
	}
	assetUidInstrumentUidMap := make(map[string]string)
	for _, v := range AssetsResponse.AssetsResponse.Assets {
		asset_uid := v.Uid

		for _, instrument := range v.Instruments {
			instrument_uid := instrument.Uid
			assetUidInstrumentUidMap[instrument_uid] = asset_uid
		}
	}
	return assetUidInstrumentUidMap, nil
}
