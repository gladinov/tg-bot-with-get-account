package tinkoffApi

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"time"

	"main.go/internal/models"
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

func (c *Client) CheckToken(ctx context.Context, tokenInBase64 string) error {
	const op = "tinkoffApi.CheckToken"

	start := time.Now()
	logg := c.logger.With(slog.String("op", op))
	logg.Debug("start")
	Path := path.Join("tinkoff", "checktoken")

	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   Path,
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("op:%s, could not create http.NewRequest", op)
	}

	req.Header.Set(models.HeaderEncryptedToken, tokenInBase64)

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("op:%s, err in client.Do", op)
	}

	body, _ := io.ReadAll(resp.Body)

	defer func() { _ = resp.Body.Close() }()
	defer func() {
		logg.Info("finished",
			slog.Duration("duration", time.Since(start)),
			slog.Int("code", resp.StatusCode),
			slog.String("body", string(body)),
		)
	}()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s:unexpected responce status code", op)
	}

	return nil
}
