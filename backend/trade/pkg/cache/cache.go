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
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	ctx         = context.Background()
)

func Init(redisURL, redisPassword string) error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: redisPassword,
		DB:       0,
	})

	for attempt := 1; attempt <= 5; attempt++ {
		_, err := redisClient.Ping(ctx).Result()
		if err == nil {
			return nil
		}
		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("could not connect to Redis after multiple attempts")
}

func Get(key string) (map[string]interface{}, error) {
	val, err := redisClient.Get(ctx, key).Result()
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

func Set(key string, value map[string]interface{}, expire time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return redisClient.Set(ctx, key, jsonValue, expire).Err()
}

func Close() {
	if redisClient != nil {
		_ = redisClient.Close()
	}
}
