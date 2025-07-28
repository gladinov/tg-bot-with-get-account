package tinkoffApi

import (
	"fmt"
	"time"

	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
	"main.go/lib/e"
)

func (c *Client) GetOperations(accountID string, date time.Time) (_ []*pb.OperationItem, err error) {
	defer func() { err = e.WrapIfErr("can't get opperations from tinkoffApi", err) }()
	resOpereaions := make([]*pb.OperationItem, 0)
	opereationsService := c.Client.NewOperationsServiceClient()
	operationsResp, err := opereationsService.GetOperationsByCursor(&investgo.GetOperationsByCursorRequest{
		AccountId: accountID,
		From:      date,
		To:        time.Now(),
		Limit:     1000,
	})
	if err != nil {
		return resOpereaions, err
	}
	operations := operationsResp.GetOperationsByCursorResponse.GetItems()
	resOpereaions = append(resOpereaions, operations...)
	nextCursor := operationsResp.NextCursor
	for nextCursor != "" {
		operationsResp, err := opereationsService.GetOperationsByCursor(&investgo.GetOperationsByCursorRequest{
			AccountId: accountID,
			Limit:     1000,
			Cursor:    nextCursor,
		})
		if err != nil {
			return resOpereaions, err
		} else {
			nextCursor = operationsResp.NextCursor
			operations := operationsResp.GetOperationsByCursorResponse.Items
			resOpereaions = append(resOpereaions, operations...)
		}
	}

	fmt.Printf("✓ Добавлено %v операции в Account.Operation по счету %s\n", len(resOpereaions), accountID)
	return resOpereaions, nil
}
