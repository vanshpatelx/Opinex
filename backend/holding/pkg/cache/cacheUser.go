// holding/pkg/cache/cache.go

/***
Redis Cache Connection

Handles connection and operations for Redis caching.

Features:
1. Manages a Redis connection.
2. Supports GET and SET operations.
3. Handles auto-reconnect in case of failures.

Author: Vansh Patel (remotevansh@gmail.com)
Last Updated: February 26, 2025
**/

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
	onceUser        sync.Once
	redisClientUser *redis.Client
	ctxUser         = context.Background()
)

// Init initializes the Redis client as a singleton
func InitUser(redisURL, redisPassword string) error {
	var initErr error

	onceUser.Do(func() {
		redisClientUser = redis.NewClient(&redis.Options{
			Addr:     redisURL,
			Password: redisPassword,
			DB:       0,
		})

		// Retry connection with exponential backoff
		for attempt := 1; attempt <= 5; attempt++ {
			_, err := redisClientUser.Ping(ctxUser).Result()
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
func GetRedisClientUser() *redis.Client {
	if redisClientUser == nil {
		InitUser(config.AppConfig.USER_RedisURL, config.AppConfig.USER_RedisPassword)
	}
	return redisClientUser
}

// Get retrieves a value from Redis
func GetUser(key string) (map[string]interface{}, error) {
	val, err := GetRedisClientUser().Get(ctxUser, key).Result()
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
func SetUser(key string, value map[string]interface{}, expire time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return GetRedisClientUser().Set(ctxUser, key, jsonValue, expire).Err()
}

// Close closes the Redis connection
func CloseUser() {
	if redisClientUser != nil {
		_ = redisClientUser.Close()
	}
}
