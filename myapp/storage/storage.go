package storage

import (
	"context"
	"errors"
)

type Storage interface {
	Save(ctx context.Context, user_name string, chatId int, token string) error
	PickToken(ctx context.Context, chatId int) (string, error)
	IsExistsToken(ctx context.Context, chatId int) (bool, error)
	CloseDB()
}

var ErrNoSaveTokens = errors.New("no saved tokens")
