package repository

import (
	"context"
	"time"

	"github.com/eganbarov/verification_code_service/config"
	"github.com/redis/go-redis/v9"
)

type CodeRepo interface {
	GetCode(phone, action string) (string, error)
	StoreCode(phone, action, code string) error
	DeleteCode(phone, action string) error
}

type CodeRepository struct {
	Redis     *redis.Client
	AppConfig *config.AppConfig
}

func (r *CodeRepository) GetCode(phone, action string) (string, error) {
	codeKey := generateCodeKey(phone, action)
	value, err := r.Redis.Get(context.Background(), codeKey).Result()
	if err != nil {
		return "", err
	}

	return value, nil
}

func (r *CodeRepository) StoreCode(phone, action, code string) error {
	codeKey := generateCodeKey(phone, action)
	codeTtl := time.Duration(r.AppConfig.CodeTtl) * time.Second
	if err := r.Redis.Set(context.Background(), codeKey, code, codeTtl).Err(); err != nil {
		return err
	}

	return nil
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
