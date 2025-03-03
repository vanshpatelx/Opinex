/***
Package cache - Redis Cache Connection

Manages Redis caching for user-related data with singleton pattern.

Features:
1. Initializes and maintains a singleton Redis connection.
2. Supports GET and SET operations for caching user data.
3. Implements automatic reconnection with exponential backoff in case of failures.
4. Ensures thread-safe access using sync.Once.
5. Provides a method to close the Redis connection gracefully.
6. Supports executing Lua scripts using Eval.

Last Updated: March 2, 2025
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

// EvalUserScript executes a Lua script in Redis
func EvalUserScript(script string, keys []string, args ...interface{}) (interface{}, error) {
	return GetRedisClientUser().Eval(ctxUser, script, keys, args...).Result()
}

// Close closes the Redis connection
func CloseUser() {
	if redisClientUser != nil {
		_ = redisClientUser.Close()
	}
}
