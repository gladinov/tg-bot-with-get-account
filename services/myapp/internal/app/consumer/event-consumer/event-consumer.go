package event_consumer

import (
	"context"
	"log/slog"
	"time"

	"github.com/gladinov/contracts/trace"
	"github.com/gladinov/e"
	"github.com/gladinov/traceidgenerator"
	"main.go/internal/app/events"
	"main.go/internal/app/events/telegram"
)

type Consumer struct {
	logger    *slog.Logger
	fetcher   events.Fetcher
	processor events.Processor
	bathcSize int
}

func New(logger *slog.Logger, fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		logger:    logger,
		fetcher:   fetcher,
		processor: processor,
		bathcSize: batchSize,
	}
}

func (c Consumer) Start() error {
	const op = "event_consumer.Start"
	logg := c.logger.With(
		slog.String("op", op),
		slog.Any("batchSize", c.bathcSize),
	)
	logg.Info("consumer started")
	for {
		gotEvents, err := c.fetcher.Fetch(c.bathcSize)
		if err != nil {
			logg.Warn("fetch events failed",
				slog.Any("error", err),
			)
			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			logg.Error("handle events failed",
				slog.Any("error", err),
				slog.Any("events count", len(gotEvents)))

			continue
		}
	}
}

func (c *Consumer) handleEvents(events []events.Event) error {
	const op = "event_consumer.handleEvents"

	start := time.Now()
	logg := c.logger.With(slog.String("op", op))
	logg.Debug("start")
	defer func() {
		logg.Info("finished",
			slog.Duration("duration", time.Since(start)),
		)
	}()

	for _, event := range events {
		traceID, err := traceidgenerator.New()
		if err != nil {
			return e.WrapIfErr(op, err)
		}
		ctx := trace.WithTraceID(context.Background(), traceID)

		eventLog := logg.With(
			slog.Any("event_type", event.Type),
		)

		// TODO: Подумать насколько это корректно? Разные логи для предусмотренных и не предусмотренных комманд
		if telegram.ContainsInConstantCommands(event.Text) {
			eventLog.DebugContext(ctx, "got new event", slog.String("event", event.Text))
		} else {
			eventLog.DebugContext(ctx, "got new other event")
		}

		if err := c.processor.Process(ctx, event); err != nil {
			if telegram.ContainsInConstantCommands(event.Text) {
				eventLog.Error("process event failed",
					slog.String("event", event.Text),
					slog.Any("error", err))
			} else {
				eventLog.Error("process other event failed",
					slog.Any("error", err))
			}
			continue
		}
	}

	return nil
}
