package tinkoffApi

import (
	"errors"

	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
)

func (c *Client) GetPortf(account *Account) error {
	// По закрытым счетам не забираем данные
	if account.Status == 3 {
		return nil
	}
	operationsService := c.Client.NewOperationsServiceClient()
	id := account.Id
	portfolioResp, err := operationsService.GetPortfolio(id,
		pb.PortfolioRequest_RUB)
	if err != nil {
		return errors.New("GetPortf: operationsService.GetPortfolio" + err.Error())
	}
	account.Portfolio = portfolioResp.GetPositions()

	return nil
}
