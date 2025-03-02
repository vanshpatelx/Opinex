// pkg/init/init.go
/**
 * Service Health Checker
 *
 * Ensures the readiness of essential services before startup:
 * - PostgreSQL Databases (User & Holding)  
 * - Redis Caches (User & Holding)  
 *
 * If any service is unavailable, the system logs the failure and exits immediately.  
 *
 * Components:
 * - **CheckUserPostgres**: Verifies connectivity to the PostgreSQL database for user-related data.  
 * - **CheckHoldingPostgres**: Verifies connectivity to the PostgreSQL database for holding-related data.  
 * - **CheckUserRedis**: Sends a PING command to ensure the Redis cache for user data is available.  
 * - **CheckHoldingRedis**: Sends a PING command to ensure the Redis cache for holding data is available.  
 * - **CheckAllServices**: Aggregates all readiness checks and determines system startup eligibility.  
 *
 * Security & Reliability Considerations:
 * - Ensures critical services are available before application execution.
 * - Logs failures explicitly for easier debugging.
 * - Uses graceful exit (`os.Exit(1)`) to prevent running in a degraded state.
 *
 * Author: Vansh Patel (remotevansh@gmail.com)  
 * Last Updated: March 2, 2025  
 */


package healthcheck

import (
	"context"
	"log"
	"os"

	_ "github.com/lib/pq"
	"holding/pkg/cache"
	"holding/pkg/db"
)

/**
    PostgreSQL Readiness Check

    Attempts to establish a connection to the PostgreSQL database.  
    Logs success or failure and returns a boolean indicating readiness.

    Returns:
    - bool → true if PostgreSQL is ready, otherwise false.  
 */
 func CheckUserPostgres() bool {
	db, err := db.GetDBUser()
	if err != nil {
		log.Println("❌ PostgreSQL User is not ready:", err)
		return false
	}

	if err = db.Ping(); err == nil {
		log.Println("✅ PostgreSQL User is ready.")
		return true
	}

	log.Println("❌ PostgreSQL User is not responding.")
	return false
}

func CheckHoldingPostgres() bool {
	db, err := db.GetDBHolding()
	if err != nil {
		log.Println("❌ PostgreSQL Holding is not ready:", err)
		return false
	}

	if err = db.Ping(); err == nil {
		log.Println("✅ PostgreSQL Holding is ready.")
		return true
	}

	log.Println("❌ PostgreSQL Holding is not responding.")
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
		log.Println("✅ Redis User is ready.")
		return true
	}

	log.Println("❌ Redis User is not responding.")
	return false
}
func CheckHoldingRedis() bool {
	client := cache.GetRedisClientHolding()

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err == nil {
		log.Println("✅ Redis User is ready.")
		return true
	}

	log.Println("❌ Redis User is not responding.")
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
