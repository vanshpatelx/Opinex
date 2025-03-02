// pkg/init/init.go
/**
    Service Health Checker

    This module checks the readiness of essential services:
    - PostgreSQL Database  
    - Redis Cache  

    If any service is unavailable, the system logs the failure and exits immediately.  

	Author: Vansh Patel (remotevansh@gmail.com)
	Last Updated: March 2, 2025
 */

package healthcheck

import (
	"context"
	"database/sql"
	"log"
	"os"

	"holding/pkg/config"

	_ "github.com/lib/pq"
	"holding/pkg/cache"
)

/**
    PostgreSQL Readiness Check

    Attempts to establish a connection to the PostgreSQL database.  
    Logs success or failure and returns a boolean indicating readiness.

    Returns:
    - bool → true if PostgreSQL is ready, otherwise false.  
 */
func CheckUserPostgres() bool {
	db, err := sql.Open("postgres", config.AppConfig.USER_DBURL)
	if err != nil {
		log.Println("❌ PostgreSQL is not ready:", err)
		return false
	}
	defer db.Close()

	if err = db.Ping(); err == nil {
		log.Println("✅ PostgreSQL is ready.")
		return true
	}

	log.Println("❌ PostgreSQL is not responding.")
	return false
}

func CheckHoldingPostgres() bool {
	db, err := sql.Open("postgres", config.AppConfig.HOLDING_DBURL)
	if err != nil {
		log.Println("❌ PostgreSQL is not ready:", err)
		return false
	}
	defer db.Close()

	if err = db.Ping(); err == nil {
		log.Println("✅ PostgreSQL is ready.")
		return true
	}

	log.Println("❌ PostgreSQL is not responding.")
	return false
}

/**
    Redis Readiness Check

    Sends a PING command to verify if the Redis cache is available.  
    Logs success or failure and returns a boolean indicating readiness.

    Returns:
    - bool → true if Redis is ready, otherwise false.  
 */
func CheckUserRedis() bool {
	client := cache.GetRedisClientUser() 

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err == nil {
		log.Println("✅ Redis is ready.")
		return true
	}

	log.Println("❌ Redis is not responding.")
	return false
}
func CheckHoldingRedis() bool {
	client := cache.GetRedisClientHolding()

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err == nil {
		log.Println("✅ Redis is ready.")
		return true
	}

	log.Println("❌ Redis is not responding.")
	return false
}

/**
    Service Initialization

    Verifies that all required services (PostgreSQL, Redis) are ready.  
    If any service is not available, logs the failure and exits the program.  
 */
func CheckAllServices() {
	log.Println("🚀 Checking service readiness...")

	userDbReady := CheckUserPostgres()
	holdingDbReady := CheckHoldingPostgres()
	userRedisReady := CheckUserRedis()
	holdingRedisReady := CheckHoldingRedis()


	if userDbReady && holdingDbReady &&  userRedisReady && holdingRedisReady{
		log.Println("✅ All services are ready!")
		return
	}

	log.Println("❌ Some services failed to initialize. Exiting...")
	os.Exit(1)
}
