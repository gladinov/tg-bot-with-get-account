package instrumentsclient

import (
	httperrors "bonds-report-service/internal/clients/http"
	"bonds-report-service/internal/models/domain"
	tinkoffDto "bonds-report-service/internal/models/dto/tinkoffApi"
	"bonds-report-service/internal/utils/logging"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"path"

	"github.com/gladinov/e"
)

func (c *InstrumentsTinkoffClient) GetFutureBy(ctx context.Context, figi string) (_ domain.Future, err error) {
	const op = "tinkoffApi.GetFutureBy"

	logg := c.logger.With()
	defer logging.LogOperation_Debug(ctx, logg, op, &err)()

	path := path.Join("tinkoff", "future")
	query := url.Values{}

	requestBody := tinkoffDto.NewFutureReq(figi)

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return domain.Future{}, e.WrapIfErr("failed json.Marshal", err)
	}
	formatRequestBody := bytes.NewBuffer(jsonData)

	httpResponse, err := c.transport.DoRequest(ctx, path, query, formatRequestBody)
	if err != nil {
		return domain.Future{}, e.WrapIfErr("failed transport DoRequest", err)
	}

	if httpResponse.StatusCode != http.StatusOK {
		return domain.Future{}, httperrors.MapHTTPError(
			httpResponse.StatusCode,
			httpResponse.Body,
		)
	}

	var data tinkoffDto.Future
	err = json.Unmarshal(httpResponse.Body, &data)
	if err != nil {
		return domain.Future{}, e.WrapIfErr("failed to unmarshal response", err)
	}

	domainData := MapFutureToDomain(data)
	return domainData, nil
}

func (c *InstrumentsTinkoffClient) GetBondByUid(ctx context.Context, uid string) (_ domain.Bond, err error) {
	const op = "tinkoffApi.GetBondByUid"

	logg := c.logger.With()
	defer logging.LogOperation_Debug(ctx, logg, op, &err)()

	path := path.Join("tinkoff", "bond")
	query := url.Values{}

	requestBody := tinkoffDto.NewBondReq(uid)

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return domain.Bond{}, e.WrapIfErr("failed json.Marshal", err)
	}

	formatRequestBody := bytes.NewBuffer(jsonData)

	httpResponse, err := c.transport.DoRequest(ctx, path, query, formatRequestBody)
	if err != nil {
		return domain.Bond{}, e.WrapIfErr("failed transport DoRequest", err)
	}

	if httpResponse.StatusCode != http.StatusOK {
		return domain.Bond{}, httperrors.MapHTTPError(
			httpResponse.StatusCode,
			httpResponse.Body,
		)
	}

	var data tinkoffDto.Bond
	err = json.Unmarshal(httpResponse.Body, &data)
	if err != nil {
		return domain.Bond{}, e.WrapIfErr("failed to unmarshal response", err)
	}

	domainData := MapBondToDomain(data)
	return domainData, nil
}

func (c *InstrumentsTinkoffClient) GetCurrencyBy(ctx context.Context, figi string) (_ domain.Currency, err error) {
	const op = "tinkoffApi.GetCurrencyBy"

	logg := c.logger.With()
	defer logging.LogOperation_Debug(ctx, logg, op, &err)()

	path := path.Join("tinkoff", "currency")
	query := url.Values{}

	requestBody := tinkoffDto.NewCurrencyReq(figi)

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return domain.Currency{}, e.WrapIfErr("failed json.Marshal", err)
	}
	formatRequestBody := bytes.NewBuffer(jsonData)

	httpResponse, err := c.transport.DoRequest(ctx, path, query, formatRequestBody)
	if err != nil {
		return domain.Currency{}, e.WrapIfErr("failed transport DoRequest", err)
	}

	if httpResponse.StatusCode != http.StatusOK {
		return domain.Currency{}, httperrors.MapHTTPError(
			httpResponse.StatusCode,
			httpResponse.Body,
		)
	}

	var data tinkoffDto.Currency
	err = json.Unmarshal(httpResponse.Body, &data)
	if err != nil {
		return domain.Currency{}, e.WrapIfErr("failed to unmarshal response", err)
	}

	domainData := MapCurrencyToDomain(data)

	return domainData, nil
}

func (c *InstrumentsTinkoffClient) FindBy(ctx context.Context, findQuery string) (_ domain.InstrumentShortList, err error) {
	const op = "tinkoffApi.FindBy"

	logg := c.logger.With()
	defer logging.LogOperation_Debug(ctx, logg, op, &err)()

	path := path.Join("tinkoff", "findby")

	requestBody := tinkoffDto.NewFindByReq(findQuery)
	query := url.Values{}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, e.WrapIfErr("failed json.Marshal", err)
	}

	formatRequestBody := bytes.NewBuffer(jsonData)

	httpResponse, err := c.transport.DoRequest(ctx, path, query, formatRequestBody)
	if err != nil {
		return nil, e.WrapIfErr("failed transport DoRequest", err)
	}

	if httpResponse.StatusCode != http.StatusOK {
		return nil, httperrors.MapHTTPError(
			httpResponse.StatusCode,
			httpResponse.Body,
		)
	}

	var data []tinkoffDto.InstrumentShort
	err = json.Unmarshal(httpResponse.Body, &data)
	if err != nil {
		return nil, e.WrapIfErr("failed to unmarshal response", err)
	}

	domainData := MapSliceInstrumentShortToDomain(data)
	return domainData, nil
}

func (c *InstrumentsTinkoffClient) GetShareCurrencyBy(ctx context.Context, figi string) (_ domain.ShareCurrency, err error) {
	const op = "tinkoffApi.GetShareCurrencyBy"

	logg := c.logger.With()
	defer logging.LogOperation_Debug(ctx, logg, op, &err)()

	path := path.Join("tinkoff", "share", "currency")
	query := url.Values{}

	requestBody := tinkoffDto.NewShareCurrencyByReq(figi)

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return domain.ShareCurrency{}, e.WrapIfErr("failed json.Marshal", err)
	}
	formatRequestBody := bytes.NewBuffer(jsonData)

	httpResponse, err := c.transport.DoRequest(ctx, path, query, formatRequestBody)
	if err != nil {
		return domain.ShareCurrency{}, e.WrapIfErr("failed transport DoRequest", err)
	}

	if httpResponse.StatusCode != http.StatusOK {
		return domain.ShareCurrency{}, httperrors.MapHTTPError(
			httpResponse.StatusCode,
			httpResponse.Body,
		)
	}

	var data tinkoffDto.ShareCurrencyByResponse
	err = json.Unmarshal(httpResponse.Body, &data)
	if err != nil {
		return domain.ShareCurrency{}, e.WrapIfErr("failed to unmarshal response", err)
	}

	domainData := MapShareCurrencyByResponseToDomain(data)
	return domainData, nil
}
