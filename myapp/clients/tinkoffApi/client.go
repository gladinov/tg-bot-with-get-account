package tinkoffApi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"
)

type Client struct {
	host   string
	client *http.Client
	Token  string
}

func NewClient(host string) *Client {
	return &Client{
		host: host,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) IsToken(token string) (err error) {

	if len(token) == 88 { // TODO:модифицировать проверку
		c.Token = token
		_, err = c.GetAccounts()
		if err != nil {
			c.Token = ""
			return err
		}
		return nil
	}
	c.Token = ""
	return errors.New("is not token")
}

func (c *Client) GetAccounts() (map[string]Account, error) {
	const op = "tinkoffApi.GetAccounts"
	Path := path.Join("tinkoff", "accounts")
	token := c.Token
	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   Path,
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("op:%s, could not create http.NewRequest", op)
	}
	authorization := "Bearer " + token
	req.Header.Set("Authorization", authorization)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("op:%s, err in client.Do", op)
	}

	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("op:%s, err in io.ReadAll", op)
	}

	if resp.StatusCode != http.StatusOK {
		var statusErr map[string]string
		err = json.Unmarshal(body, &statusErr)
		if err != nil {
			return nil, fmt.Errorf("op:%s, could not unmarshall json, delete this block. err : %s", op, err.Error())
		}
		return nil, fmt.Errorf("op:%s, err:"+statusErr["error"], op)
	}

	data := make(map[string]Account, 0)
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("op:%s, could not unmarshall json, delete this block. err : %s", op, err.Error())
	}
	return data, nil
}

func (c *Client) GetPortfolio(requestBody PortfolioRequest) (Portfolio, error) {
	const op = "tinkoffApi.GetPortfolio"
	token := c.Token
	Path := path.Join("tinkoff", "portfolio")

	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   Path,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return Portfolio{}, fmt.Errorf("op:%s, could not marshall JSON", op)
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return Portfolio{}, fmt.Errorf("op:%s, could not create http.NewRequest", op)
	}

	authorization := "Bearer " + token
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authorization)

	resp, err := c.client.Do(req)
	if err != nil {
		return Portfolio{}, fmt.Errorf("op:%s, err in client.Do", op)
	}

	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Portfolio{}, fmt.Errorf("op:%s, err in io.ReadAll", op)
	}

	if resp.StatusCode != http.StatusOK {
		var statusErr map[string]string
		err = json.Unmarshal(body, &statusErr)
		if err != nil {
			return Portfolio{}, fmt.Errorf("op:%s, could not unmarshall json", op)
		}
		return Portfolio{}, fmt.Errorf("op:%s, err:"+statusErr["error"], op)
	}

	var data Portfolio
	err = json.Unmarshal(body, &data)
	if err != nil {
		return Portfolio{}, fmt.Errorf("op:%s, could not unmarshall json", op)
	}
	return data, nil

}

func (c *Client) GetOperations(requestBody OperationsRequest) (_ []Operation, err error) {
	const op = "tinkoffApi.GetOperations"
	Path := path.Join("tinkoff", "operations")
	token := c.Token
	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   Path,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("op:%s, could not marshall JSON", op)
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("op:%s, could not create http.NewRequest", op)
	}

	authorization := "Bearer " + token
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authorization)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("op:%s, err in client.Do", op)
	}

	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("op:%s, err in io.ReadAll", op)
	}
	if resp.StatusCode != http.StatusOK {
		var statusErr map[string]string
		err = json.Unmarshal(body, &statusErr)
		if err != nil {
			return nil, fmt.Errorf("op:%s, could not unmarshall json, delete this block. err : %s", op, err.Error())
		}
		return nil, fmt.Errorf("op:%s, err:"+statusErr["error"], op)
	}

	var data []Operation
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("op:%s, could not unmarshall json", op)
	}
	return data, nil

}

func (c *Client) GetAllAssetUids() (map[string]string, error) {
	const op = "tinkoffApi.GetAllAssetUids"
	Path := path.Join("tinkoff", "allassetsuid")
	token := c.Token
	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   Path,
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("op:%s, could not create http.NewRequest", op)
	}
	authorization := "Bearer " + token
	req.Header.Set("Authorization", authorization)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("op:%s, err in client.Do", op)
	}

	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("op:%s, err in io.ReadAll", op)
	}
	if resp.StatusCode != http.StatusOK {
		var statusErr map[string]string
		err = json.Unmarshal(body, &statusErr)
		if err != nil {
			return nil, fmt.Errorf("op:%s, could not unmarshall json, delete this block. err : %s", op, err.Error())
		}
		return nil, fmt.Errorf("op:%s, err:"+statusErr["error"], op)
	}

	var data map[string]string
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("op:%s, could not unmarshall json", op)
	}
	return data, nil
}

func (c *Client) GetFutureBy(figi string) (Future, error) {
	const op = "tinkoffApi.GetFutureBy"
	Path := path.Join("tinkoff", "future")
	token := c.Token
	requestBody := FutureReq{
		Figi: figi,
	}

	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   Path,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return Future{}, fmt.Errorf("op:%s, could not marshall JSON", op)
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return Future{}, fmt.Errorf("op:%s, could not create http.NewRequest", op)
	}

	authorization := "Bearer " + token
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authorization)

	resp, err := c.client.Do(req)
	if err != nil {
		return Future{}, fmt.Errorf("op:%s, err in client.Do", op)
	}

	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Future{}, fmt.Errorf("op:%s, err in io.ReadAll", op)
	}
	if resp.StatusCode != http.StatusOK {
		var statusErr map[string]string
		err = json.Unmarshal(body, &statusErr)
		if err != nil {
			return Future{}, fmt.Errorf("op:%s, could not unmarshall json, delete this block. err : %s", op, err.Error())
		}
		return Future{}, fmt.Errorf("op:%s, err:"+statusErr["error"], op)
	}

	var data Future
	err = json.Unmarshal(body, &data)
	if err != nil {
		return Future{}, fmt.Errorf("op:%s, could not unmarshall json", op)
	}
	return data, nil
}

func (c *Client) GetBondByUid(uid string) (Bond, error) {
	const op = "tinkoffApi.GetBondByUid"
	Path := path.Join("tinkoff", "bond")
	token := c.Token
	requestBody := BondReq{
		Uid: uid,
	}

	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   Path,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return Bond{}, fmt.Errorf("op:%s, could not marshall JSON", op)
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return Bond{}, fmt.Errorf("op:%s, could not create http.NewRequest", op)
	}

	authorization := "Bearer " + token
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authorization)

	resp, err := c.client.Do(req)
	if err != nil {
		return Bond{}, fmt.Errorf("op:%s, err in client.Do", op)
	}

	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Bond{}, fmt.Errorf("op:%s, err in io.ReadAll", op)
	}
	if resp.StatusCode != http.StatusOK {
		var statusErr map[string]string
		err = json.Unmarshal(body, &statusErr)
		if err != nil {
			return Bond{}, fmt.Errorf("op:%s, could not unmarshall json, delete this block. err : %s", op, err.Error())
		}
		return Bond{}, fmt.Errorf("op:%s, err:"+statusErr["error"], op)
	}

	var data Bond
	err = json.Unmarshal(body, &data)
	if err != nil {
		return Bond{}, fmt.Errorf("op:%s, could not unmarshall json", op)
	}
	return data, nil
}

func (c *Client) GetCurrencyBy(figi string) (Currency, error) {
	const op = "tinkoffApi.GetCurrencyBy"
	Path := path.Join("tinkoff", "currency")
	token := c.Token
	requestBody := CurrencyReq{
		Figi: figi,
	}

	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   Path,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return Currency{}, fmt.Errorf("op:%s, could not marshall JSON", op)
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return Currency{}, fmt.Errorf("op:%s, could not create http.NewRequest", op)
	}

	authorization := "Bearer " + token
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authorization)

	resp, err := c.client.Do(req)
	if err != nil {
		return Currency{}, fmt.Errorf("op:%s, err in client.Do", op)
	}

	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Currency{}, fmt.Errorf("op:%s, err in io.ReadAll", op)
	}
	if resp.StatusCode != http.StatusOK {
		var statusErr map[string]string
		err = json.Unmarshal(body, &statusErr)
		if err != nil {
			return Currency{}, fmt.Errorf("op:%s, could not unmarshall json, delete this block. err : %s", op, err.Error())
		}
		return Currency{}, fmt.Errorf("op:%s, err:"+statusErr["error"], op)
	}

	var data Currency
	err = json.Unmarshal(body, &data)
	if err != nil {
		return Currency{}, fmt.Errorf("op:%s, could not unmarshall json", op)
	}
	return data, nil
}

func (c *Client) FindBy(query string) ([]InstrumentShort, error) {
	const op = "tinkoffApi.FindBy"
	Path := path.Join("tinkoff", "findby")
	token := c.Token
	requestBody := FindByReq{
		Query: query,
	}

	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   Path,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("op:%s, could not marshall JSON", op)
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("op:%s, could not create http.NewRequest", op)
	}

	authorization := "Bearer " + token
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authorization)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("op:%s, err in client.Do", op)
	}

	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("op:%s, err in io.ReadAll", op)
	}
	if resp.StatusCode != http.StatusOK {
		var statusErr map[string]string
		err = json.Unmarshal(body, &statusErr)
		if err != nil {
			return nil, fmt.Errorf("op:%s, could not unmarshall json, delete this block. err : %s", op, err.Error())
		}
		return nil, fmt.Errorf("op:%s, err:"+statusErr["error"], op)
	}

	var data []InstrumentShort
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("op:%s, could not unmarshall json", op)
	}
	return data, nil
}

func (c *Client) GetBondsActions(instrumentUid string) (BondIdentIdentifiers, error) {
	const op = "tinkoffApi.GetBondsActions"
	Path := path.Join("tinkoff", "bondactions")
	token := c.Token
	requestBody := BondsActionsReq{
		InstrumentUid: instrumentUid,
	}

	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   Path,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return BondIdentIdentifiers{}, fmt.Errorf("op:%s, could not marshall JSON", op)
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return BondIdentIdentifiers{}, fmt.Errorf("op:%s, could not create http.NewRequest", op)
	}

	authorization := "Bearer " + token
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authorization)

	resp, err := c.client.Do(req)
	if err != nil {
		return BondIdentIdentifiers{}, fmt.Errorf("op:%s, err in client.Do", op)
	}

	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return BondIdentIdentifiers{}, fmt.Errorf("op:%s, err in io.ReadAll", op)
	}
	if resp.StatusCode != http.StatusOK {
		var statusErr map[string]string
		err = json.Unmarshal(body, &statusErr)
		if err != nil {
			return BondIdentIdentifiers{}, fmt.Errorf("op:%s, could not unmarshall json, delete this block. err : %s", op, err.Error())
		}
		return BondIdentIdentifiers{}, fmt.Errorf("op:%s, err:"+statusErr["error"], op)
	}

	var data BondIdentIdentifiers
	err = json.Unmarshal(body, &data)
	if err != nil {
		return BondIdentIdentifiers{}, fmt.Errorf("op:%s, could not unmarshall json", op)
	}
	return data, nil
}

func (c *Client) GetLastPriceInPersentageToNominal(instrumentUid string) (LastPriceResponse, error) {
	const op = "tinkoffApi.GetLastPriceInPersentageToNominal"
	Path := path.Join("tinkoff", "lastprice")
	token := c.Token
	requestBody := LastPriceReq{
		InstrumentUid: instrumentUid,
	}

	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   Path,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return LastPriceResponse{}, fmt.Errorf("op:%s, could not marshall JSON", op)
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return LastPriceResponse{}, fmt.Errorf("op:%s, could not create http.NewRequest", op)
	}

	authorization := "Bearer " + token
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authorization)

	resp, err := c.client.Do(req)
	if err != nil {
		return LastPriceResponse{}, fmt.Errorf("op:%s, err in client.Do", op)
	}

	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LastPriceResponse{}, fmt.Errorf("op:%s, err in io.ReadAll", op)
	}
	if resp.StatusCode != http.StatusOK {
		var statusErr map[string]string
		err = json.Unmarshal(body, &statusErr)
		if err != nil {
			return LastPriceResponse{}, fmt.Errorf("op:%s, could not unmarshall json, delete this block. err : %s", op, err.Error())
		}
		return LastPriceResponse{}, fmt.Errorf("op:%s, err:"+statusErr["error"], op)
	}

	var data LastPriceResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return LastPriceResponse{}, fmt.Errorf("op:%s, could not unmarshall json", op)
	}
	return data, nil
}

func (c *Client) GetShareCurrencyBy(figi string) (ShareCurrencyByResponse, error) {
	const op = "tinkoffApi.GetShareCurrencyBy"
	Path := path.Join("tinkoff", "share", "currency")
	token := c.Token
	requestBody := ShareCurrencyByRequest{
		Figi: figi,
	}

	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   Path,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return ShareCurrencyByResponse{}, fmt.Errorf("op:%s, could not marshall JSON", op)
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return ShareCurrencyByResponse{}, fmt.Errorf("op:%s, could not create http.NewRequest", op)
	}

	authorization := "Bearer " + token
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authorization)

	resp, err := c.client.Do(req)
	if err != nil {
		return ShareCurrencyByResponse{}, fmt.Errorf("op:%s, err in client.Do", op)
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ShareCurrencyByResponse{}, fmt.Errorf("op:%s, err in io.ReadAll", op)
	}
	if resp.StatusCode != http.StatusOK {
		var statusErr map[string]string
		err = json.Unmarshal(body, &statusErr)
		if err != nil {
			return ShareCurrencyByResponse{}, fmt.Errorf("op:%s, could not unmarshall json, delete this block. err : %s", op, err.Error())
		}
		return ShareCurrencyByResponse{}, fmt.Errorf("op:%s, err:"+statusErr["error"], op)
	}

	var data ShareCurrencyByResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return ShareCurrencyByResponse{}, fmt.Errorf("op:%s, could not unmarshall json", op)
	}
	return data, nil
}
