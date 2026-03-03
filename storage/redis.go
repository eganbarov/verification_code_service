package storage

import (
	"context"
	"fmt"

	"github.com/eganbarov/verification_code_service/config"
	"github.com/redis/go-redis/v9"
)

func NewClient(ctx context.Context, cfg *config.StorageConfig) (*redis.Client, error) {
	db := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		DB:           cfg.DB,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	})

	if err := db.Ping(ctx).Err(); err != nil {
		fmt.Printf("Failed to connect to redis server %s\n", err.Error())
		return nil, err
	}

	return db, nil
}
