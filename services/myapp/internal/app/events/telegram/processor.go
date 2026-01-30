package telegram

import (
	"context"
	"errors"
	"log/slog"

	"github.com/gladinov/e"
	bondreportservice "main.go/clients/bondReportService"
	"main.go/clients/telegram"
	"main.go/clients/tinkoffApi"
	"main.go/internal/app/events"
	tokenauth "main.go/internal/tokenAuth"
)

type Processor struct {
	logger            *slog.Logger
	tg                *telegram.Client
	tinkoffApi        *tinkoffApi.Client
	bondReportService *bondreportservice.Client
	tokenAuthService  *tokenauth.TokenAuthService
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func NewProccesor(
	logger *slog.Logger,
	client *telegram.Client,
	tinkoffApiClient *tinkoffApi.Client,
	bondReportServiceClient *bondreportservice.Client,
	tokenAuthService *tokenauth.TokenAuthService,
) *Processor {
	return &Processor{
		logger:            logger,
		tg:                client,
		tinkoffApi:        tinkoffApiClient,
		bondReportService: bondReportServiceClient,
		tokenAuthService:  tokenAuthService,
	}
}

func (p *Processor) Process(ctx context.Context, event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(ctx, event)
	default:
		return e.Wrap("can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(ctx context.Context, event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	if err := p.doCmd(ctx, event.Text, meta.ChatID, meta.Username); err != nil {
		return e.Wrap("can't process message", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}
	return res, nil
}
