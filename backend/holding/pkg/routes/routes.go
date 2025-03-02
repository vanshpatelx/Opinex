// holding/pkg/routes/routes.go
/***
Routes Package

Registers API endpoints and applies authentication middleware.

Features:
1. **RegisterHoldingRoutes**:
   - Defines and registers routes related to holdings.
   - Secures endpoints using JWT authentication middleware.

Security Considerations:
- Middleware ensures only authenticated users can access holdings data.

Endpoints:
   - GET /holdings/:userID → Fetch user holdings.
   - GET /holdings/balance/:userID → Fetch user balance.

Author: Vansh Patel (remotevansh@gmail.com)
Last Updated: February 26, 2025
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
