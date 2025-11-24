package telegram

import (
	"errors"

	"main.go/clients/telegram"
	"main.go/events"
	"main.go/lib/cryptoToken"
	"main.go/lib/e"
	"main.go/service"
	"main.go/storage"
)

type Processor struct {
	tokenCrypter *cryptoToken.TokenCrypter
	tg           *telegram.Client
	storage      storage.Storage
	service      *service.Client
}

type Meta struct {
	ChatID   int
	Username string
}

var ErrUnknownEventType = errors.New("unknown event type")
var ErrUnknownMetaType = errors.New("unknown meta type")

func NewProccesor(tokenCrypter *cryptoToken.TokenCrypter, client *telegram.Client, storage storage.Storage, serviceCLient *service.Client) *Processor {
	return &Processor{
		tokenCrypter: tokenCrypter,
		tg:           client,
		storage:      storage,
		service:      serviceCLient,
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
