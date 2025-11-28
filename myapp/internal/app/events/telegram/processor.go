package telegram

import (
	"errors"

	"github.com/redis/go-redis/v9"
	bondreportservice "main.go/clients/bondReportService"
	"main.go/clients/telegram"
	"main.go/clients/tinkoffApi"
	"main.go/internal/app/events"
	"main.go/lib/cryptoToken"
	"main.go/lib/e"

	storage "main.go/internal/repository"
)

type Processor struct {
	tokenCrypter      *cryptoToken.TokenCrypter
	tg                *telegram.Client
	tinkoffApi        *tinkoffApi.Client
	bondReportService *bondreportservice.Client
	redis             *redis.Client
	storage           storage.Storage
}

type Meta struct {
	ChatID   int
	Username string
}

var ErrUnknownEventType = errors.New("unknown event type")
var ErrUnknownMetaType = errors.New("unknown meta type")

func NewProccesor(
	tokenCrypter *cryptoToken.TokenCrypter,
	client *telegram.Client,
	tinkoffApiClient *tinkoffApi.Client,
	bondReportServiceClient *bondreportservice.Client,
	redisClient *redis.Client,
	storage storage.Storage) *Processor {
	return &Processor{
		tokenCrypter:      tokenCrypter,
		tg:                client,
		tinkoffApi:        tinkoffApiClient,
		bondReportService: bondReportServiceClient,
		redis:             redisClient,
		storage:           storage,
	}
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return e.Wrap("can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
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
