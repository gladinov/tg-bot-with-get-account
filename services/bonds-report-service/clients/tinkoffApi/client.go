package tinkoffApi

import (
	"bonds-report-service/lib/valuefromcontext"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"time"
)

type Client struct {
	logger *slog.Logger
	host   string
	client *http.Client
}

func NewClient(logger *slog.Logger, host string) *Client {
	return &Client{
		logger: logger,
		host:   host,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) GetAccounts(ctx context.Context) (map[string]Account, error) {
	const op = "tinkoffApi.GetAccounts"
	Path := path.Join("tinkoff", "accounts")
	chatID, err := valuefromcontext.GetChatIDFromCtxStr(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	start := time.Now()
	logg := c.logger.With(slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("finished",
			slog.Duration("duration", time.Since(start)),
			slog.Any("err", err),
		)
	}()

	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   Path,
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("op:%s, could not create http.NewRequest", op)
	}

	req.Header.Set(valuefromcontext.HeaderChatID, chatID)

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

func (c *Client) GetPortfolio(ctx context.Context, requestBody PortfolioRequest) (Portfolio, error) {
	const op = "tinkoffApi.GetPortfolio"

	start := time.Now()
	logg := c.logger.With(slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("finished",
			slog.Duration("duration", time.Since(start)),
		)
	}()

	chatID, err := valuefromcontext.GetChatIDFromCtxStr(ctx)
	if err != nil {
		return Portfolio{}, fmt.Errorf("%s:%w", op, err)
	}
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

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(valuefromcontext.HeaderChatID, chatID)

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

func (c *Client) GetOperations(ctx context.Context, requestBody OperationsRequest) (_ []Operation, err error) {
	const op = "tinkoffApi.GetOperations"

	start := time.Now()
	logg := c.logger.With(slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("finished",
			slog.Duration("duration", time.Since(start)),
		)
	}()

	Path := path.Join("tinkoff", "operations")
	chatID, err := valuefromcontext.GetChatIDFromCtxStr(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
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

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(valuefromcontext.HeaderChatID, chatID)

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

func (c *Client) GetAllAssetUids(ctx context.Context) (map[string]string, error) {
	const op = "tinkoffApi.GetAllAssetUids"

	start := time.Now()
	logg := c.logger.With(slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("finished",
			slog.Duration("duration", time.Since(start)),
		)
	}()

	Path := path.Join("tinkoff", "allassetsuid")
	chatID, err := valuefromcontext.GetChatIDFromCtxStr(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   Path,
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("op:%s, could not create http.NewRequest", op)
	}
	req.Header.Set(valuefromcontext.HeaderChatID, chatID)

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

func (c *Client) GetFutureBy(ctx context.Context, figi string) (Future, error) {
	const op = "tinkoffApi.GetFutureBy"

	start := time.Now()
	logg := c.logger.With(slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("finished",
			slog.Duration("duration", time.Since(start)),
		)
	}()

	Path := path.Join("tinkoff", "future")
	chatID, err := valuefromcontext.GetChatIDFromCtxStr(ctx)
	if err != nil {
		return Future{}, fmt.Errorf("%s:%w", op, err)
	}
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

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(valuefromcontext.HeaderChatID, chatID)

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

func (c *Client) GetBondByUid(ctx context.Context, uid string) (Bond, error) {
	const op = "tinkoffApi.GetBondByUid"

	start := time.Now()
	logg := c.logger.With(slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("finished",
			slog.Duration("duration", time.Since(start)),
		)
	}()

	Path := path.Join("tinkoff", "bond")
	chatID, err := valuefromcontext.GetChatIDFromCtxStr(ctx)
	if err != nil {
		return Bond{}, fmt.Errorf("%s:%w", op, err)
	}
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

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(valuefromcontext.HeaderChatID, chatID)

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

func (c *Client) GetCurrencyBy(ctx context.Context, figi string) (Currency, error) {
	const op = "tinkoffApi.GetCurrencyBy"

	start := time.Now()
	logg := c.logger.With(slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("finished",
			slog.Duration("duration", time.Since(start)),
		)
	}()

	Path := path.Join("tinkoff", "currency")
	chatID, err := valuefromcontext.GetChatIDFromCtxStr(ctx)
	if err != nil {
		return Currency{}, fmt.Errorf("%s:%w", op, err)
	}
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

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(valuefromcontext.HeaderChatID, chatID)

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

func (c *Client) FindBy(ctx context.Context, query string) ([]InstrumentShort, error) {
	const op = "tinkoffApi.FindBy"

	start := time.Now()
	logg := c.logger.With(slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("finished",
			slog.Duration("duration", time.Since(start)),
		)
	}()

	Path := path.Join("tinkoff", "findby")
	chatID, err := valuefromcontext.GetChatIDFromCtxStr(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
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

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(valuefromcontext.HeaderChatID, chatID)

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

func (c *Client) GetBondsActions(ctx context.Context, instrumentUid string) (BondIdentIdentifiers, error) {
	const op = "tinkoffApi.GetBondsActions"

	start := time.Now()
	logg := c.logger.With(slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("finished",
			slog.Duration("duration", time.Since(start)),
		)
	}()

	Path := path.Join("tinkoff", "bondactions")
	chatID, err := valuefromcontext.GetChatIDFromCtxStr(ctx)
	if err != nil {
		return BondIdentIdentifiers{}, fmt.Errorf("%s:%w", op, err)
	}
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

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(valuefromcontext.HeaderChatID, chatID)

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

func (c *Client) GetLastPriceInPersentageToNominal(ctx context.Context, instrumentUid string) (LastPriceResponse, error) {
	const op = "tinkoffApi.GetLastPriceInPersentageToNominal"

	start := time.Now()
	logg := c.logger.With(slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("finished",
			slog.Duration("duration", time.Since(start)),
		)
	}()

	Path := path.Join("tinkoff", "lastprice")
	chatID, err := valuefromcontext.GetChatIDFromCtxStr(ctx)
	if err != nil {
		return LastPriceResponse{}, fmt.Errorf("%s:%w", op, err)
	}
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

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(valuefromcontext.HeaderChatID, chatID)

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

func (c *Client) GetShareCurrencyBy(ctx context.Context, figi string) (ShareCurrencyByResponse, error) {
	const op = "tinkoffApi.GetShareCurrencyBy"

	start := time.Now()
	logg := c.logger.With(slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("finished",
			slog.Duration("duration", time.Since(start)),
		)
	}()

	Path := path.Join("tinkoff", "share", "currency")

	chatID, err := valuefromcontext.GetChatIDFromCtxStr(ctx)
	if err != nil {
		return ShareCurrencyByResponse{}, fmt.Errorf("%s:%w", op, err)
	}
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

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(valuefromcontext.HeaderChatID, chatID)

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
