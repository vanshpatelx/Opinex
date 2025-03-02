/***
 * Routes Package
 *
 * Registers API endpoints for holdings and applies authentication middleware.
 *
 * Features:
 * - Defines and registers holdings-related routes.
 * - Secures endpoints with JWT authentication middleware.
 *
 * Security Considerations:
 * - Ensures only authenticated users can access holdings data.
 *
 * Endpoints:
 * - GET /holdings/:userID → Fetch user holdings.
 * - GET /holdings/balance/ → Fetch user balance (userID as query param).
 *
 * Author: Vansh Patel (remotevansh@gmail.com)
 * Last Updated: February 26, 2025
 ***/


package routes

import (
	"holding/pkg/controller"
	"holding/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterHoldingRoutes(app *fiber.App) {
	holdingRoutes := app.Group("/holdings")

	// holdingRoutes.Get("/:userID", controller.GetHoldingByID)
	holdingRoutes.Get("/balance",middleware.JWTMiddleware,  controller.GetBalanceByID)
}
