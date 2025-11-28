package redisClient

import (
	"context"
	"tinkoffApi/internal/configs"

	"github.com/redis/go-redis/v9"
)

func NewClient(ctx context.Context, cfg *configs.Config) (*redis.Client, error) {
	db := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHTTPServer.Address,
		Password: cfg.RedisHTTPServer.Password,
		DB:       cfg.RedisHTTPServer.DB,
		// Username:     cfg.RedisHTTPServer.User,
		MaxRetries:   cfg.RedisHTTPServer.MaxRetries,
		DialTimeout:  cfg.RedisHTTPServer.DialTimeout,
		ReadTimeout:  cfg.RedisHTTPServer.Timeout,
		WriteTimeout: cfg.RedisHTTPServer.Timeout,
	})

	if err := db.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return db, nil
}
