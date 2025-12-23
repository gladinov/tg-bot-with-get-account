package cbr

import (
	"bonds-report-service/lib/e"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"time"
)

const (
	layout = "02/01/2006"
)

type Client struct {
	logger *slog.Logger
	host   string
	client http.Client
}

func New(logger *slog.Logger, host string) *Client {
	return &Client{
		logger: logger,
		host:   host,
		client: http.Client{},
	}
}

func (c *Client) GetAllCurrencies(date time.Time) (res CurrenciesResponce, err error) {
	const op = "cbr.GetAllCurrencies"
	start := time.Now()
	logg := c.logger.With(slog.String("op", op))
	logg.Debug("start", slog.Time("date", date))

	defer func() {
		logg.Info("finished",
			slog.Duration("duration", time.Since(start)),
			slog.Any("err", err),
		)
		err = e.WrapIfErr(op, err)
	}()

	request := CurrencyRequest{Date: date}
	Path := path.Join("cbr", "currencies")
	params := url.Values{}

	requestBody, err := json.Marshal(request)
	if err != nil {
		logg.Debug("failed to marshal request", slog.Any("err", err))
		return CurrenciesResponce{}, err
	}
	formatRequestBody := bytes.NewBuffer(requestBody)

	httpResponse, err := c.doRequest(Path, params, formatRequestBody)
	if err != nil {
		logg.Debug("http request failed", slog.Any("err", err))
		return CurrenciesResponce{}, err
	}

	switch httpResponse.StatusCode {
	case http.StatusBadRequest:
		err = fmt.Errorf("op:%s, statusCode:%v, error: Invalid request", op, httpResponse.StatusCode)
		logg.Debug("bad request", slog.Any("err", err))
		return CurrenciesResponce{}, err
	case http.StatusInternalServerError:
		err = fmt.Errorf("op:%s, statusCode:%v, error: could not get currencies from cbr", op, httpResponse.StatusCode)
		logg.Debug("internal server error", slog.Any("err", err))
		return CurrenciesResponce{}, err
	}

	err = json.Unmarshal(httpResponse.Body, &res)
	if err != nil {
		logg.Debug("failed to unmarshal response", slog.Any("err", err))
		return CurrenciesResponce{}, err
	}

	logg.Debug("successfully fetched currencies", slog.Int("count", len(res.Currencies)))
	return res, nil
}

func (c *Client) doRequest(Path string, query url.Values, requestBody io.Reader) (resp HTTPResponse, err error) {
	const op = "cbr.doRequest"
	start := time.Now()
	logg := c.logger.With(slog.String("op", op))
	logg.Debug("start", slog.String("path", Path))

	defer func() {
		logg.Info("finished",
			slog.Duration("duration", time.Since(start)),
			slog.Any("err", err),
		)
		err = e.WrapIfErr("could not do request", err)
	}()

	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   Path,
	}
	req, err := http.NewRequest(http.MethodPost, u.String(), requestBody)
	if err != nil {
		logg.Debug("failed to create http request", slog.Any("err", err))
		return HTTPResponse{}, err
	}
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Content-Type", "application/json")

	response, err := c.client.Do(req)
	if err != nil {
		logg.Debug("http client error", slog.Any("err", err))
		return HTTPResponse{}, err
	}
	defer func() { _ = response.Body.Close() }()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logg.Debug("failed to read response body", slog.Any("err", err))
		return HTTPResponse{}, err
	}

	resp = HTTPResponse{
		StatusCode: response.StatusCode,
		Body:       body,
	}
	logg.Debug("http request successful", slog.Int("statusCode", resp.StatusCode))
	return resp, nil
}
