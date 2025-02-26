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
	"log"
	"strconv"
	"time"

	"holding/pkg/cache"
	"holding/pkg/db"

	"github.com/gofiber/fiber/v2"
)

/***
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
***/
func GetHoldingByID(c *fiber.Ctx) error {
	userID := c.Params("userID")
	userIDJWT, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Authorization check
	if userID != userIDJWT && c.Locals("role") != "ADMIN" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized access"})
	}

	// Try fetching from cache
	cacheKey := "Holding:UserID:" + userID
	if cachedData, err := cache.Get(cacheKey); err == nil {
		if holdings, exists := cachedData["holdings"]; exists {
			return c.JSON(fiber.Map{"message": "Fetched Holdings successfully", "userID": userID, "holdings": holdings})
		}
	}

	// Fetch from DB
	results, err := db.FetchAll("SELECT eventID, price, quantity, locked FROM holdings WHERE userID = $1;", userID)
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
	cache.Set(cacheKey, map[string]interface{}{"holdings": results}, 24*time.Hour)

	return c.JSON(fiber.Map{"message": "Fetched Holdings successfully", "userID": userID, "holdings": results})
}

/*** 
GetBalanceByID retrieves the balance of a user by user ID.
Steps:
1. Extract userID from request parameters and verify JWT authorization.
2. Check if balance data is available in the cache.
3. If cached data exists, return it immediately.
4. Query the database for the user's balance if not found in the cache.
5. Ensure the balance value is a valid string before returning.
6. Store the fetched balance in the cache for future requests.
7. Return the processed balance data as JSON response.

Returns:
- 200 OK with user balance if found.
- 401 Unauthorized if user lacks permission.
- 500 Internal Server Error if DB query fails.

Last Updated: February 26, 2025
***/
func GetBalanceByID(c *fiber.Ctx) error {
	userID := c.Params("userID")
	userIDJWT, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Authorization check
	if userID != userIDJWT && c.Locals("role") != "ADMIN" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized access"})
	}

	// Try fetching from cache
	cacheKey := "Balance:UserID:" + userID
	if data, err := cache.Get(cacheKey); err == nil {
		if balance, exists := data["balance"].(string); exists {
			return c.JSON(fiber.Map{"message": "Fetched Balance successfully", "userID": userID, "balance": balance})
		}
	}

	// Fetch from DB
	result, err := db.FetchOne("SELECT balance FROM users WHERE id = $1;", userID)
	if err != nil {
		log.Println("DB Error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch balance from DB"})
	}

	// Ensure balance is a string
	balance, ok := result["balance"].(string)
	if !ok {
		log.Println("Invalid balance data type:", result["balance"])
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid balance data type"})
	}

	// Update cache
	cache.Set(cacheKey, fiber.Map{"balance": balance}, 24*time.Hour)

	return c.JSON(fiber.Map{"message": "Fetched Balance successfully", "userID": userID, "balance": balance})
}
