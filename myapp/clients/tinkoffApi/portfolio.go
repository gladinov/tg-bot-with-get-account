package tinkoffApi

import (
	"errors"

	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
	"main.go/lib/e"
)

var ErrCloseAccount = errors.New("close account havn't portffolio positions")

type Portfolio struct {
	Positions   []*pb.PortfolioPosition
	TotalAmount float64
}

func (c *Client) GetPortf(accountID string, accountStatus int64) (_ Portfolio, err error) {
	portfolio := Portfolio{}
	if accountStatus == 3 {
		return portfolio, ErrCloseAccount
	}
	operationsService := c.Client.NewOperationsServiceClient()
	id := accountID
	portfolioResp, err := operationsService.GetPortfolio(id,
		pb.PortfolioRequest_RUB)
	if err != nil {
		return portfolio, e.WrapIfErr("can't get portifolio positions from tinkoff Api", err)
	}
	portfolio.Positions = portfolioResp.GetPositions()
	portfolio.TotalAmount = portfolioResp.GetTotalAmountPortfolio().ToFloat()

	return portfolio, nil
}
