package moex

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"
)

type Client struct {
	host   string
	client *http.Client
}

func NewClient(host string) *Client {
	return &Client{
		host: host,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) GetSpecifications(ticker string, date time.Time) (*SpecificationsResponce, error) {
	var data *SpecificationsResponce
	Path := path.Join("specifications")

	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   Path,
	}

	requestData := SpecificationsRequest{
		Ticker: ticker,
		Date:   date,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
