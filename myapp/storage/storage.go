package storage

import (
	"context"
	"errors"
)

type Storage interface {
	Save(ctx context.Context, user_name string, chatId int, token string) error
	PickToken(ctx context.Context, chatId int) (string, error)
	Remove(ctx context.Context, p *Page) error
	IsExists(ctx context.Context, chatId int) (bool, error)
}

var ErrNoSavePages = errors.New("no saved pages")

type Page struct {
	URL      string
	UserName string
	ChatId   int
}
