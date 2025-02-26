// holding/pkg/routes/routes.go
package routes

import (
	"holding/pkg/middleware"
	"log"
	"strconv"
	"time"

	"holding/pkg/cache"
	"holding/pkg/db"

	"github.com/gofiber/fiber/v2"
)

// RegisterHoldingRoutes sets up the routes for holdings.
func RegisterHoldingRoutes(app *fiber.App) {
	holdingRoutes := app.Group("/holdings")
	holdingRoutes.Use(middleware.JWTMiddleware)

	holdingRoutes.Get("/:userID", getHoldingByID)
	holdingRoutes.Get("/balance/:userID", getBalanceByID)
}

// getHoldingByID retrieves holding details by user ID.
func getHoldingByID(c *fiber.Ctx) error {
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch holdings from DB"})
	}

	// Check if there are results
	if len(results) == 0 {
		return c.JSON(fiber.Map{"message": "No holdings found", "userID": userID, "holdings": []interface{}{}})
	}

	// Convert price from string to int
	for i, row := range results {
		if priceStr, ok := row["price"].([]byte); ok {
			priceFloat, err := strconv.ParseFloat(string(priceStr), 64)
			if err == nil {
				results[i]["price"] = int(priceFloat) // Convert to int
			} else {
				results[i]["price"] = 0 // Fallback if conversion fails
			}
		}
	}

	// Update cache (Store array directly)
	cache.Set(cacheKey, map[string]interface{}{"holdings": results}, 24*time.Hour)

	return c.JSON(fiber.Map{"message": "Fetched Holdings successfully", "userID": userID, "holdings": results})
}


// getBalanceByID retrieves balance by user ID.
func getBalanceByID(c *fiber.Ctx) error {
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
		return c.JSON(fiber.Map{"message": "Fetched Balance successfully", "userID": userID, "balance": data["balance"]})
	}

	// Fetch from DB
	result, err := db.FetchOne("SELECT balance FROM users WHERE id = $1;", userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch balance from DB"})
	}

	balance, ok := result["balance"].([]byte)
	if !ok {
		log.Println("Invalid balance data type")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid balance data type"})
	}

	// Convert balance and update cache
	balanceStr := string(balance)
	cache.Set(cacheKey, fiber.Map{"balance": balanceStr}, 24*time.Hour)

	return c.JSON(fiber.Map{"message": "Fetched Balance successfully", "userID": userID, "balance": balanceStr})
}
