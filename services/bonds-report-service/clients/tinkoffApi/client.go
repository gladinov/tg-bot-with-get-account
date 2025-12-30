package tinkoffApi

import (
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

	"github.com/gladinov/valuefromcontext"

	httpheaders "github.com/gladinov/contracts/http"
	"github.com/gladinov/contracts/trace"
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

func (c *Client) GetAccounts(ctx context.Context) (_ map[string]Account, err error) {
	const op = "tinkoffApi.GetAccounts"
	Path := path.Join("tinkoff", "accounts")

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

	reqWithHeaders, err := c.setHeaders(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("err set hearders", err)
	}

	resp, err := c.client.Do(reqWithHeaders)
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
	reqWithHeaders, err := c.setHeaders(ctx, req)
	if err != nil {
		return Portfolio{}, fmt.Errorf("err set hearders", err)
	}

	resp, err := c.client.Do(reqWithHeaders)
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
	reqWithHeaders, err := c.setHeaders(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("err set hearders", err)
	}

	resp, err := c.client.Do(reqWithHeaders)
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

	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   Path,
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("op:%s, could not create http.NewRequest", op)
	}
	reqWithHeaders, err := c.setHeaders(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("err set hearders", err)
	}

	resp, err := c.client.Do(reqWithHeaders)
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
	reqWithHeaders, err := c.setHeaders(ctx, req)
	if err != nil {
		return Future{}, fmt.Errorf("err set hearders", err)
	}

	resp, err := c.client.Do(reqWithHeaders)
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
	reqWithHeaders, err := c.setHeaders(ctx, req)
	if err != nil {
		return Bond{}, fmt.Errorf("err set hearders", err)
	}

	resp, err := c.client.Do(reqWithHeaders)
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
	reqWithHeaders, err := c.setHeaders(ctx, req)
	if err != nil {
		return Currency{}, fmt.Errorf("err set hearders", err)
	}

	resp, err := c.client.Do(reqWithHeaders)
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
	reqWithHeaders, err := c.setHeaders(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("err set hearders", err)
	}

	resp, err := c.client.Do(reqWithHeaders)
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
	reqWithHeaders, err := c.setHeaders(ctx, req)
	if err != nil {
		return BondIdentIdentifiers{}, fmt.Errorf("err set hearders", err)
	}

	resp, err := c.client.Do(reqWithHeaders)
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
	reqWithHeaders, err := c.setHeaders(ctx, req)
	if err != nil {
		return LastPriceResponse{}, fmt.Errorf("err set hearders", err)
	}

	resp, err := c.client.Do(reqWithHeaders)
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
	reqWithHeaders, err := c.setHeaders(ctx, req)
	if err != nil {
		return ShareCurrencyByResponse{}, fmt.Errorf("err set hearders", err)
	}

	resp, err := c.client.Do(reqWithHeaders)
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

func (c *Client) setHeaders(ctx context.Context, req *http.Request) (*http.Request, error) {
	const op = "bondreportservice.SetHeaders"

	logg := c.logger.With(slog.String("op", op))
	logg.DebugContext(ctx, "start")
	defer func() {
		logg.InfoContext(ctx, "finished")
	}()
	chatID, err := valuefromcontext.GetChatIDFromCtxStr(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	req.Header.Set(httpheaders.HeaderChatID, chatID)
	traceID, ok := trace.TraceIDFromContext(ctx)
	if !ok {
		logg.WarnContext(ctx, "hasn't traceID in ctx")
	}
	req.Header.Set(httpheaders.HeaderTraceID, traceID)

	return req, nil
}
