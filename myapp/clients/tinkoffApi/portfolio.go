package tinkoffApi

import (
	"errors"

	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
	"main.go/lib/e"
)

var ErrCloseAccount = errors.New("close account havn't portffolio positions")

func (c *Client) GetPortf(accountID string, accountStatus int64) (_ []*pb.PortfolioPosition, err error) {
	portffolioPosition := make([]*pb.PortfolioPosition, 0)
	if accountStatus == 3 {
		return portffolioPosition, ErrCloseAccount
	}
	operationsService := c.Client.NewOperationsServiceClient()
	id := accountID
	portfolioResp, err := operationsService.GetPortfolio(id,
		pb.PortfolioRequest_RUB)
	if err != nil {
		return portffolioPosition, e.WrapIfErr("can't get portifolio positions from tinkoff Api", err)
	}
	portffolioPosition = portfolioResp.GetPositions()

	return portffolioPosition, nil
}
