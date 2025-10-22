package tinkoffApi

import (
	"context"

	"github.com/russianinvestments/invest-api-go-sdk/investgo"

	"main.go/clients/moex"
	"main.go/lib/e"
	pathwd "main.go/lib/pathWD"
)

const (
	tinkoffApiConfigPath = "/configs/tinkoffApiConfig.yaml"
)

type Client struct {
	ctx        context.Context
	Logg       investgo.Logger
	config     *investgo.Config
	Client     *investgo.Client
	MoexClient *moex.Client
}

func New(ctx context.Context, logg investgo.Logger) *Client {
	return &Client{
		ctx:  context.Background(),
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
	tinkoffConfigAbsolutPath, err := pathwd.PathFromWD(tinkoffApiConfigPath)
	if err != nil {
		panic("can't create absolute path to tinkoffApi Config")
	}
	config, err := investgo.LoadConfig(tinkoffConfigAbsolutPath)
	if err != nil {
		c.Logg.Errorf("incorrect path by config: %s", tinkoffConfigAbsolutPath)
		panic("can't load config")
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
