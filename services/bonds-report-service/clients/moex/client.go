package moex

import (
	"bytes"
	"encoding/json"
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

func (c *Client) GetSpecifications(ticker string, date time.Time) (data Values, err error) {
	const op = "moex.GetSpecifications"
	start := time.Now()
	logg := c.logger.With(slog.String("op", op))
	logg.Debug("start", slog.String("ticker", ticker), slog.Time("date", date))

	defer func() {
		logg.Info("finished",
			slog.Duration("duration", time.Since(start)),
			slog.Any("err", err),
		)
	}()

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
		logg.Debug("failed to marshal request", slog.Any("err", err))
		return Values{}, err
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		logg.Debug("failed to create HTTP request", slog.Any("err", err))
		return Values{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		logg.Debug("HTTP request failed", slog.Any("err", err))
		return Values{}, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logg.Debug("failed to read response body", slog.Any("err", err))
		return Values{}, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		logg.Debug("failed to unmarshal response", slog.Any("err", err))
		return Values{}, err
	}

	logg.Debug("successfully fetched specifications")
	return data, nil
}
