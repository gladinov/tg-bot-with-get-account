package cbr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"
)

const (
	layout = "02/01/2006"
)

type Client struct {
	host   string
	client http.Client
}

func New(host string) *Client {
	return &Client{
		host:   host,
		client: http.Client{},
	}
}

func (c *Client) GetAllCurrencies(date time.Time) (CurrenciesResponce, error) {
	const op = "cbr.GetAllCurrencies"
	request := CurrencyRequest{Date: date}
	Path := path.Join("cbr", "currencies")
	params := url.Values{}
	requestBody, err := json.Marshal(request)
	if err != nil {
		return CurrenciesResponce{}, err
	}
	formatRequestBody := bytes.NewBuffer(requestBody)

	httpResponse, err := c.doRequest(Path, params, formatRequestBody)
	if err != nil {
		return CurrenciesResponce{}, err
	}
	switch httpResponse.StatusCode {
	case http.StatusBadRequest:
		return CurrenciesResponce{}, fmt.Errorf("op:%s, statusCode:%v, error: Invalid request", op, httpResponse.StatusCode)
	case http.StatusInternalServerError:
		return CurrenciesResponce{}, fmt.Errorf("op:%s, statusCode:%v, error: could not get currencies from cbr", op, httpResponse.StatusCode)
	}
	var currencies CurrenciesResponce
	err = json.Unmarshal(httpResponse.Body, &currencies)
	if err != nil {
		return CurrenciesResponce{}, err
	}

	return currencies, nil
}

func (c *Client) doRequest(Path string, query url.Values, requestBody io.Reader) (HTTPResponse, error) {
	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   Path,
	}
	req, err := http.NewRequest(http.MethodPost, u.String(), requestBody)
	if err != nil {
		return HTTPResponse{}, err
	}
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return HTTPResponse{}, err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return HTTPResponse{}, err
	}
	httpResponse := HTTPResponse{
		StatusCode: resp.StatusCode,
		Body:       body,
	}
	return httpResponse, nil
}
