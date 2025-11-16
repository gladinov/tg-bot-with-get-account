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

func (c *Client) GetSpecifications(ticker string, date time.Time) (Values, error) {
	var data Values
	Path := path.Join("moex", "specifications")

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
		return Values{}, err
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return Values{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return Values{}, err
	}

	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Values{}, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return Values{}, err
	}
	
	return data, nil
}
