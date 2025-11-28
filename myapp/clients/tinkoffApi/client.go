package tinkoffApi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"main.go/internal/models"
	"main.go/lib/valuefromcontext"
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

func (c *Client) CheckToken(ctx context.Context) error {
	const op = "tinkoffApi.GetAccounts"
	Path := path.Join("tinkoff", "checktoken")
	chatID, err := valuefromcontext.GetChatIDFromCtxStr(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   Path,
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("op:%s, could not create http.NewRequest", op)
	}

	req.Header.Set(models.HeaderChatID, chatID)

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("op:%s, err in client.Do", op)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s:unexpected responce status code", op)
	}

	return nil
}
