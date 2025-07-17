package service

import (
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
)

func (s *Client) TransOperations(operations []*pb.OperationItem) []Operation {
	transformOperations := make([]Operation, 0)
	for _, v := range operations {
		transformOperation := Operation{
			Currency:          v.GetPrice().Currency,
			BrokerAccountId:   v.GetBrokerAccountId(),
			Operation_Id:      v.GetId(),
			ParentOperationId: v.GetParentOperationId(),
			Name:              v.GetName(),
			Date:              v.Date.AsTime(),
			Type:              int64(v.GetType()),
			Description:       v.GetDescription(),
			InstrumentUid:     v.GetInstrumentUid(),
			Figi:              v.GetFigi(),
			InstrumentType:    v.GetInstrumentType(),
			InstrumentKind:    string(v.GetInstrumentKind()),
			PositionUid:       v.GetPositionUid(),
			Payment:           v.GetPayment().ToFloat(),
			Price:             v.GetPrice().ToFloat(),
			Commission:        v.GetCommission().ToFloat(),
			Yield:             v.GetYield().ToFloat(),
			YieldRelative:     v.GetYieldRelative().ToFloat(),
			AccruedInt:        v.GetAccruedInt().ToFloat(),
			QuantityDone:      float64(v.GetQuantityDone()),
			AssetUid:          v.GetAssetUid(),
		}
		transformOperations = append(transformOperations, transformOperation)
	}
	return transformOperations
}
