package kafka

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/gladinov/e"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Producer struct {
	logg  *slog.Logger
	kafka *kgo.Client
}

func NewProducer(logg *slog.Logger, kafka *kgo.Client) *Producer {
	return &Producer{
		logg:  logg,
		kafka: kafka,
	}
}

func (p *Producer) PublishRequest(
	ctx context.Context,
	reportKind string,
	chatID string,
	traceID string,
) error {
	// TODO: Доставать traceID и chatID из контекста
	resp := NewRequest(reportKind, traceID, chatID)

	body, err := json.Marshal(resp)
	if err != nil {
		return e.WrapIfErr("failed to marshall data", err)
	}

	record := kgo.Record{
		Topic: ReportRequested,
		Value: body,
	}
	// TODO: ПОчитать как проверять и обрабатывать ошибку при отправке
	p.kafka.Produce(ctx, &record, nil)
	return nil
}
