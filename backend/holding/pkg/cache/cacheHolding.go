package cache

import (
	"context"
	"encoding/json"
	"holding/pkg/config"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	onceHolding        sync.Once
	redisClientHolding *redis.Client
	ctxHolding         = context.Background()
)

// Init initializes the Redis client as a singleton
func InitHolding(redisURL, redisPassword string) error {
	var initErr error

	onceHolding.Do(func() {
		redisClientHolding = redis.NewClient(&redis.Options{
			Addr:     redisURL,
			Password: redisPassword,
			DB:       0,
		})

		// Retry connection with exponential backoff
		for attempt := 1; attempt <= 5; attempt++ {
			_, err := redisClientHolding.Ping(ctxHolding).Result()
			if err == nil {
				initErr = nil
				return
			}
			initErr = err
			time.Sleep(time.Duration(attempt) * 2 * time.Second) // Exponential backoff
		}
	})

	return initErr
}

// GetRedisClient returns the singleton Redis client
func GetRedisClientHolding() *redis.Client {
	if redisClientHolding == nil {
		InitHolding(config.AppConfig.HOLDING_RedisURL, config.AppConfig.HOLDING_RedisPassword)
	}
	return redisClientHolding
}

// Get retrieves a value from Redis
func GetHolding(key string) (map[string]interface{}, error) {
	val, err := GetRedisClientHolding().Get(ctxHolding, key).Result()
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(val), &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Set stores a value in Redis
func SetHolding(key string, value map[string]interface{}, expire time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return GetRedisClientHolding().Set(ctxHolding, key, jsonValue, expire).Err()
}

// Close closes the Redis connection
func CloseHolding() {
	if redisClientHolding != nil {
		_ = redisClientHolding.Close()
	}
}
