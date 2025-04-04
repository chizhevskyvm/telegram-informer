package redis

import (
	"encoding/json"
	"time"
)

// Cache интерфейс для кэша
type Cache interface {
	Set(key string, value string, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
}

func GetTyped[T any](cache Cache, key string) (T, error) {
	str, err := cache.Get(key)
	if err != nil {
		var zero T
		return zero, err
	}
	var result T
	err = json.Unmarshal([]byte(str), &result)
	return result, err
}

func SetTyped[T any](cache Cache, key string, value T, expiration time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return cache.Set(key, string(bytes), expiration)
}
