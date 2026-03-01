package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type CodeRepository struct {
	Redis *redis.Client
}

func (r *CodeRepository) StoreCode(phone, action, code string) error {
	codeKey := generateCodeKey(phone, action)
	if err := r.Redis.Set(context.Background(), codeKey, code, 300*time.Second).Err(); err != nil {
		return err
	}

	return nil
}

func (r *CodeRepository) GetCode(phone, action string) (string, error) {
	codeKey := generateCodeKey(phone, action)
	value, err := r.Redis.Get(context.Background(), codeKey).Result()
	if err != nil {
		return "", err
	}

	return value, nil
}

func (r *CodeRepository) DeleteCode(phone, action string) error {
	codeKey := generateCodeKey(phone, action)
	if err := r.Redis.Del(context.Background(), codeKey).Err(); err != nil {
		return err
	}

	return nil
}

func generateCodeKey(phone, action string) string {
	return phone + "_" + action
}
