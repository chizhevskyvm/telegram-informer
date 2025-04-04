package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type CacheStorage struct {
	Client *redis.Client
	Ctx    context.Context
}

func New(ctx context.Context, addr, password string, db int) (*CacheStorage, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Проверим подключение
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping error: %w", err)
	}

	return &CacheStorage{Client: rdb, Ctx: ctx}, nil
}

func (r *CacheStorage) GetOrSet(key string, value any, expiration time.Duration) (any, error) {
	result, err := r.Client.Get(r.Ctx, key).Result()
	if err == nil {
		// Ключ уже есть — возвращаем значение
		return result, nil
	}

	if err != redis.Nil {
		// Какая-то другая ошибка
		return "", err
	}

	// Ключа нет — пробуем установить (SetNX = "Set if Not Exists")
	ok, err := r.Client.SetNX(r.Ctx, key, value, expiration).Result()
	if err != nil {
		return "", err
	}

	if ok {
		// Установили только что
		return value, nil
	}

	// Кто-то уже установил — читаем снова
	return r.Client.Get(r.Ctx, key).Result()
}

func (r *CacheStorage) Set(key string, value string, expiration time.Duration) error {
	return r.Client.Set(r.Ctx, key, value, expiration).Err()
}

func (r *CacheStorage) Get(key string) (string, error) {
	return r.Client.Get(r.Ctx, key).Result()
}

func (r *CacheStorage) Delete(key string) error {
	return r.Client.Del(r.Ctx, key).Err()
}
