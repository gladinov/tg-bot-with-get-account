package tinkoffapi

import (
	"context"

	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	"main.go/lib/e"
)

type Client struct {
	ctx    context.Context
	Logg   investgo.Logger
	config *investgo.Config
	Client *investgo.Client
}

func New(ctx context.Context, logg investgo.Logger) *Client {
	return &Client{ctx: context.Background(),
		Logg: logg}
}

func (c *Client) FillClient(token string) (err error) {
	defer func() { err = e.WrapIfErr("can't create Tinkoff Client", err) }()

	if err := c.getConfig(token); err != nil {
		return err
	}

	if err = c.getClient(); err != nil {
		return err
	}
	return nil
}

func (c *Client) getConfig(token string) error {
	config, err := investgo.LoadConfig("./configs/tinkoffApiConfig.yaml")
	if err != nil {
		return e.Wrap("can't load config", err)
	}
	c.config = &config
	c.config.Token = token
	return nil
}

func (c *Client) getClient() error {
	client, err := investgo.NewClient(c.ctx, *c.config, c.Logg)
	if err != nil {
		return e.WrapIfErr("can't connect with tinkoffApi client", err)
	}
	c.Client = client
	return nil
}
