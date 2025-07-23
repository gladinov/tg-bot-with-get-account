package tinkoffApi

import (
	"errors"
	"fmt"
	"time"

	"github.com/russianinvestments/invest-api-go-sdk/investgo"
)

func (c *Client) GetOpp(account *Account, date time.Time) error {
	opereationsService := c.Client.NewOperationsServiceClient()
	operationsResp, err := opereationsService.GetOperationsByCursor(&investgo.GetOperationsByCursorRequest{
		AccountId: account.Id,
		From:      date,
		To:        time.Now(),
		Limit:     1000,
	})
	if err != nil {
		return errors.New("GetOpp: operationsService.GetOperationsByCursor" + err.Error())
	}
	operations := operationsResp.GetOperationsByCursorResponse.GetItems()
	account.Operations = append(account.Operations, operations...)
	nextCursor := operationsResp.NextCursor
	for nextCursor != "" {
		operationsResp, err := opereationsService.GetOperationsByCursor(&investgo.GetOperationsByCursorRequest{
			AccountId: account.Id,
			Limit:     1000,
			Cursor:    nextCursor,
		})
		if err != nil {
			return errors.New("GetOpp: operationsService.GetOperationsByCursor" + err.Error())
		} else {
			nextCursor = operationsResp.NextCursor
			operations := operationsResp.GetOperationsByCursorResponse.Items
			account.Operations = append(account.Operations, operations...)
		}
	}

	fmt.Printf("✓ Добавлено %v операции в Account.Operation по счету %s\n", len(account.Operations), account.Id)
	return nil
}
