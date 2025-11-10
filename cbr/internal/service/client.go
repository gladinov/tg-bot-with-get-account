package service

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
}

type Transport struct {
	host   string
	client http.Client
}

func NewTransport(host string) *Transport {
	return &Transport{host: host,
		client: http.Client{}}
}

func NewClient(transport HTTPTransport) *Client {
	return &Client{
		transport: transport,
	}
}

func (c *Client) GetAllCurrencies(formatDate string) (CurrenciesResponce, error) {
	const op = "service.GetAllCurrencies"

	Path := path.Join("scripts", "XML_daily.asp")

	params := url.Values{}
	params.Add("date_req", formatDate)

	body, err := c.transport.DoRequest(Path, params)
	if err != nil {
		return CurrenciesResponce{}, fmt.Errorf("op: %s, error: could not do request", op)
	}
	return c.parseCurrencies(body)
}

func (c *Client) parseCurrencies(data []byte) (CurrenciesResponce, error) {
	const op = "service.parseCurrencies"
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

	return curr, nil
}

func (t *Transport) DoRequest(Path string, query url.Values) ([]byte, error) {
	const op = "service.doRequest"
	u := url.URL{
		Scheme: "https",
		Host:   t.host,
		Path:   Path,
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("op: %s, error: could not create http.NewRequest", op)
	}
	req.URL.RawQuery = query.Encode()
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; MyApp/1.0)")

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("op: %s, error: could not do request", op)
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("op: %s, error: could not read body", op)
	}

	return body, nil
}
