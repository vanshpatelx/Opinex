// holding/pkg/controller/controller.go

/***
Controller Package

Handles HTTP request processing for user holdings and balance.

Features:
1. **GetHoldingByID**:
   - Retrieves holdings for a given user ID.
   - Authorizes requests based on JWT and user role.
   - Caches holdings data for optimized performance.
   - Converts price values safely from string to integer.

2. **GetBalanceByID**:
   - Retrieves user balance by user ID.
   - Implements authorization and cache retrieval.
   - Ensures balance data integrity before returning a response.

Dependencies:
- Fiber for HTTP handling.
- Custom cache package for Redis-based caching.
- Custom db package for database queries.

Author: Vansh Patel (remotevansh@gmail.com)
Last Updated: February 26, 2025
***/

package controller

import (
	"holding/pkg/cache"
	"holding/pkg/db"
	"log"
	"strconv"
	"time"

	"holding/pkg/common"
	"github.com/gofiber/fiber/v2"
	"holding/pkg/rabbitmqQueue"
)

var rabbitMQProducer *rabbitmqQueue.RabbitMQConsumer

func SetRabbitMQProducer(producer *rabbitmqQueue.RabbitMQConsumer) {
    rabbitMQProducer = producer
}

/*
**
GetHoldingByID retrieves holding details by user ID.

Steps:
1. Extract userID from request parameters and verify JWT authorization.
2. Check if holdings data is available in the cache.
3. If cached data exists, return it immediately.
4. Query the database for the user's holdings if not found in the cache.
5. Convert price values from string to integer safely.
6. Store the fetched holdings in the cache for future requests.
7. Return the processed holdings data as JSON response.

Returns:
- 200 OK with user holdings if found.
- 401 Unauthorized if user lacks permission.
- 500 Internal Server Error if DB query fails.

Last Updated: February 26, 2025
**
*/
func GetHoldingByID(c *fiber.Ctx) error {
	userID := c.Params("userID")
	userIDJWT, ok := c.Locals("Id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Authorization check
	if userID != userIDJWT && c.Locals("Role") != "ADMIN" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized access"})
	}

	// Try fetching from cache
	cacheKey := "Holding:UserID:" + userID
	if cachedData, err := cache.GetHolding(cacheKey); err == nil {
		if holdings, exists := cachedData["holdings"]; exists {
			return c.JSON(fiber.Map{"message": "Fetched Holdings successfully", "userID": userID, "holdings": holdings})
		}
	}

	// Fetch from DB
	results, err := db.FetchAllHolding("SELECT user_id, balance, lb, FROM holdings WHERE userID = $1;", userID)
	if err != nil {
		log.Println("DB Error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch holdings from DB"})
	}

	// Check if there are results
	if len(results) == 0 {
		return c.JSON(fiber.Map{"message": "No holdings found", "userID": userID, "holdings": []interface{}{}})
	}

	// Convert price from string to int safely
	for i, row := range results {
		priceStr, ok := row["price"].(string)
		if !ok {
			log.Println("Warning: Unexpected price type", row["price"])
			results[i]["price"] = 0
			continue
		}
		priceFloat, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			log.Println("Error parsing price:", err)
			results[i]["price"] = 0
			continue
		}
		results[i]["price"] = int(priceFloat)
	}

	// Update cache
	cache.SetHolding(cacheKey, map[string]interface{}{"holdings": results}, 24*time.Hour)

	return c.JSON(fiber.Map{"message": "Fetched Holdings successfully", "userID": userID, "holdings": results})
}

/***
GetBalanceByID retrieves a user's balance by their user ID.

Workflow:
1. Extract userID from request parameters and verify JWT authorization.
2. Validate if the requesting user has permission (self or admin access).
3. Check if balance data is available in the cache.
   - If found, return the cached balance immediately.
4. If not found in cache, query the database for the balance.
5. Validate that the retrieved balance and locked balance are valid numeric values.
6. Store the fetched balance in the cache for future requests.
7. Return the balance data as a JSON response.

Returns:
- 200 OK with balance details if successful.
- 401 Unauthorized if the user lacks permission.
- 500 Internal Server Error for database or data validation errors.

Last Updated: March 2, 2025
***/

func GetBalanceByID(c *fiber.Ctx) error {
	userIDJWT, ok1 := c.Locals("Id").(string)
	if !ok1 {
		log.Printf("Invalid user ID: Queried by %v", c.Locals("Id"))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	userID := c.Query("userID")
	if userID == "" {
		log.Printf("No userID provided in query, using JWT userID: %v", userIDJWT)
		userID = userIDJWT
	}

	log.Printf("GetBalance by ID Query for userID: %v by user: %v", userIDJWT, userID)

	// Authorization check
	if userID != userIDJWT && c.Locals("Role") != "ADMIN" {
		log.Printf("Unauthorized access by userID: %v for user: %v", userIDJWT, userID)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized access"})
	}

	// Try fetching from cache
	cacheKey := "Balance:UserID:" + userIDJWT
	log.Print("check in cache")
	if data, err := cache.GetUser(cacheKey); err == nil {
		if balance, exists1 := data["balance"].(float64); exists1 {
			log.Print("check in cache 3")
			if lb, exists2 := data["lb"].(float64); exists2 {
				log.Printf("Fetched from cache for userID %v - balance: %v, locked balance: %v", userIDJWT, balance, lb)
				return c.JSON(fiber.Map{
					"message":        "Fetched Balance successfully",
					"userID":         userIDJWT,
					"balance":        balance,
					"locked_balance": lb,
				})
			}
		}
	}

	// Fetch from DB
	result, err := db.FetchOneUser(`SELECT balance, lb FROM "User" WHERE user_id = $1;`, userID)
	if err != nil {
		log.Printf("DB Error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	// Ensure balance and lb are valid
	balance, ok1 := result["balance"].(int64)
	lb, ok2 := result["lb"].(int64)

	if !ok1 || !ok2 {
		log.Printf("Invalid balance or lb data type: %v", result)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	// Update cache
	cache.SetUser(cacheKey, fiber.Map{"balance": balance, "lb": lb}, 60*time.Minute)

	return c.JSON(fiber.Map{
		"message":        "Fetched Balance successfully",
		"user_id":        userID,
		"balance":        balance,
		"locked_balance": lb,
	})
}




func UpdateBalanceByID(c *fiber.Ctx) error {
	userIDJWT, ok1 := c.Locals("Id").(string)
	if !ok1 {
		log.Printf("Invalid user ID: Queried by %v", c.Locals("Id"))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	userID := c.Query("userID")
	if userID == "" {
		log.Printf("No userID provided in query, using JWT userID: %v", userIDJWT)
		userID = userIDJWT
	}

	// Parse balance from request body
	var request struct {
		Balance float64 `json:"balance"`
	}
	if err := c.BodyParser(&request); err != nil {
		log.Printf("Invalid request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	balance := request.Balance

	log.Printf("UpdateBalance request for userID: %v by user: %v with new balance: %.2f", userID, userIDJWT, balance)

	// Authorization check
	if userID != userIDJWT && c.Locals("Role") != "ADMIN" {
		log.Printf("Unauthorized access by userID: %v for user: %v", userIDJWT, userID)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized access"})
	}

	if err := common.UpdateBalance("addBalance", userID, balance); err != nil {
		log.Printf("Failed to update balance for userID: %s, Error: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	// Send to pubsub ✅ Use RabbitMQ Publisher
	if rabbitMQProducer != nil {
		success := rabbitMQProducer.PublishMessage("user_balance_exchange", "balance.updated", map[string]interface{}{
			"userID":  userID,
			"balance": balance,
		})

		if !success {
			// Rollback balance update in case of failed message publish
			if err := common.UpdateBalance("decreaseBalance", userID, balance); err != nil {
				log.Printf("Failed to rollback balance for userID: %s, Error: %v", userID, err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
		}
	}

	return c.JSON(fiber.Map{"message": "Balance updated successfully"})
}
