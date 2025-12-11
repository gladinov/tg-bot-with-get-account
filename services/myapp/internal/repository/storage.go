package storage

import (
	"context"
	"errors"

	"main.go/internal/config"
	"main.go/internal/repository/postgres"
	"main.go/internal/repository/sqlite"
	pathwd "main.go/lib/pathWD"
)

// myapp\internal\repository\postgres\postrges.go

const (
	postreSQL = "postgreSQL"
	SQLite    = "SQLite"
)

type Storage interface {
	Save(ctx context.Context, user_name string, token string) error
	PickToken(ctx context.Context) (string, error)
	IsExistsToken(ctx context.Context) (bool, error)
	CloseDB()
}

func NewStorage(ctx context.Context, config config.Config) (Storage, error) {
	switch config.DbType {
	case postreSQL:
		storage, err := postgres.NewStorage(config)
		if err != nil {
			return nil, err
		}
		err = storage.Init(ctx)
		if err != nil {
			return nil, err
		}

		return storage, nil
	case SQLite:
		storageAbsolutPath, err := pathwd.PathFromWD(config.RootPath, config.StorageSQLLitePath)
		if err != nil {
			return nil, err
		}
		storage, err := sqlite.New(storageAbsolutPath)
		if err != nil {
			return nil, err
		}
		err = storage.Init(ctx)
		if err != nil {
			return nil, err
		}

		return storage, nil
	default:
		return nil, errors.New("possible init only SQLite or PostgreSQL databases")
	}
}
