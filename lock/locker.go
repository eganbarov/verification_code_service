package lock

import (
	"context"
	"time"

	"github.com/eganbarov/verification_code_service/config"
	"github.com/redis/go-redis/v9"
)

type Locker interface {
	Lock(phone, action string) error
	IsLocked(phone, action string) bool
	Release(phone, action string) error
}

type RedisLocker struct {
	Redis     *redis.Client
	AppConfig *config.AppConfig
}

func (l *RedisLocker) Lock(phone, action string) error {
	lockKey := generateKey(phone, action)
	lockTtl := time.Duration(l.AppConfig.RepeatSentCodeTtl) * time.Second
	if err := l.Redis.Set(context.Background(), lockKey, 1, lockTtl).Err(); err != nil {
		return err
	}

	return nil
}

func (l *RedisLocker) Release(phone, action string) error {
	lockKey := generateKey(phone, action)
	if err := l.Redis.Del(context.Background(), lockKey).Err(); err != nil {
		return err
	}

	return nil
}

func (l *RedisLocker) IsLocked(phone, action string) bool {
	lockKey := generateKey(phone, action)
	_, err := l.Redis.Get(context.Background(), lockKey).Result()
	if err != nil {
		return false
	}

	return true
}

func generateKey(phone, action string) string {
	return phone + "_" + action + "_is_sent"
}
