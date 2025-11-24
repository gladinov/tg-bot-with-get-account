package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"main.go/storage"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w by path:%s", err, path)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) Save(ctx context.Context, user_name string, chatID int, token string) error {
	q := `INSERT INTO users(user_name, chatID, token) VALUES (?,?,?)`

	if _, err := s.db.ExecContext(ctx, q, user_name, chatID, token); err != nil {
		return fmt.Errorf("can't save page: %w", err)
	}
	return nil
}

func (s *Storage) PickToken(ctx context.Context, chatId int) (string, error) {
	q := `SELECT token FROM users WHERE chatID = ? LIMIT 1`

	var token string

	err := s.db.QueryRowContext(ctx, q, chatId).Scan(&token)
	if err == sql.ErrNoRows {
		return "", storage.ErrNoSaveTokens
	}
	if err != nil {
		return "", fmt.Errorf("can't pick token: %w", err)

	}

	return token, nil
}

func (s *Storage) IsExistsToken(ctx context.Context, chatId int) (bool, error) {
	q := `SELECT token FROM users WHERE chatID = ?`

	var token sql.NullString

	err := s.db.QueryRowContext(ctx, q, chatId).Scan(&token)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("can't check user in storage: %w", err)
	}

	return token.Valid, nil
}

func (s *Storage) Init(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *Storage) CloseDB() {
	_ = s.db.Close()
}
