package service

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path"

	"golang.org/x/text/encoding/charmap"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=HTTPClient
type HTTPClient interface {
	GetAllCurrencies(formatDate string) (CurrenciesResponce, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=HTTPTransport
type HTTPTransport interface {
	DoRequest(Path string, query url.Values) ([]byte, error)
}

type Client struct {
	transport HTTPTransport
	logger    *slog.Logger
}

type Transport struct {
	host   string
	client http.Client
	logger *slog.Logger
}

func NewTransport(host string, logger *slog.Logger) *Transport {
	return &Transport{host: host,
		client: http.Client{},
		logger: logger}
}

func NewClient(transport HTTPTransport, logger *slog.Logger) *Client {
	return &Client{
		transport: transport,
		logger:    logger,
	}
}

func (c *Client) GetAllCurrencies(formatDate string) (CurrenciesResponce, error) {
	const op = "service.GetAllCurrencies"
	logg := c.logger.With(slog.String("function", op))
	logg.Info("start " + op)
	logg.Debug("input data", slog.String("input", formatDate))

	Path := path.Join("scripts", "XML_daily.asp")

	params := url.Values{}
	params.Add("date_req", formatDate)

	body, err := c.transport.DoRequest(Path, params)
	if err != nil {
		return CurrenciesResponce{}, fmt.Errorf("op: %s, error: could not do request", op)
	}
	logg.Info("eng " + op)
	return c.parseCurrencies(body)
}

func (c *Client) parseCurrencies(data []byte) (CurrenciesResponce, error) {
	const op = "service.parseCurrencies"
	logg := c.logger.With(slog.String("function", op))
	logg.Info("start " + op)

	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.CharsetReader = func(label string, input io.Reader) (io.Reader, error) {
		if label == "windows-1251" {
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		}
		return input, nil
	}
	var curr CurrenciesResponce
	err := decoder.Decode(&curr)
	if err != nil {
		return CurrenciesResponce{}, fmt.Errorf("op: %s, could not decode Xml file", op)
	}
	logg.Debug("output", slog.Any("currencies", curr))
	logg.Info("eng " + op)
	return curr, nil
}

func (t *Transport) DoRequest(Path string, query url.Values) ([]byte, error) {
	const op = "service.doRequest"
	logg := t.logger.With(slog.String("function", op))
	logg.Info("start " + op)

	u := url.URL{
		Scheme: "https",
		Host:   t.host,
		Path:   Path,
	}
	logg.Debug("create new request")
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		logg.Error("could not create new request",
			slog.String("error", err.Error()),
			slog.String("endpoint", req.URL.String()),
			slog.String("method", req.Method),
			slog.Bool("is_timeout", os.IsTimeout(err)),
		)
		return nil, fmt.Errorf("op: %s, error: could not create http.NewRequest", op)
	}

	req.URL.RawQuery = query.Encode()
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; MyApp/1.0)")

	logg.Debug("http.Client.Do")
	resp, err := t.client.Do(req)
	if err != nil {
		logg.Error("could not do http",
			slog.String("error", err.Error()),
			slog.String("endpoint", req.URL.String()),
			slog.String("method", req.Method),
			slog.Bool("is_timeout", os.IsTimeout(err)),
		)
		return nil, fmt.Errorf("op: %s, error: could not do request", op)
	}

	defer func() { _ = resp.Body.Close() }()
	logg.Debug("io.ReadAll")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logg.Error("could not io.ReadALL",
			slog.String("error", err.Error()),
			slog.String("endpoint", req.URL.String()),
			slog.String("method", req.Method),
			slog.Bool("is_timeout", os.IsTimeout(err)),
		)
		return nil, fmt.Errorf("op: %s, error: could not read body", op)
	}
	logg.Info("eng " + op)
	return body, nil
}
