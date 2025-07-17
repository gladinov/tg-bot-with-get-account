package storage

import (
	"context"
	"errors"

	"main.go/service"
)

type Storage interface {
	Save(ctx context.Context, user_name string, chatId int, token string) error
	SavePositions(ctx context.Context, chatID int, accountId string, positions []service.PortfolioPosition) error
	SaveOperations(ctx context.Context, chatID int, accountId string, operations []service.Operation) error
	SaveBondReport(ctx context.Context, chatID int, accountId string, bondReport []service.BondReport) error
	PickToken(ctx context.Context, chatId int) (string, error)
	Remove(ctx context.Context, p *Page) error
	IsExists(ctx context.Context, chatId int) (bool, error)
	GetOperations(ctx context.Context, chatId int, assetUid string, accountId string) ([]service.Operation, error)
}

var ErrNoSavePages = errors.New("no saved pages")

type Page struct {
	URL      string
	UserName string
	ChatId   int
}
