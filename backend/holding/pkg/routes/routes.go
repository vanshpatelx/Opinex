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
	"github.com/gofiber/fiber/v2"
	"holding/pkg/controller"
	"holding/pkg/middleware"
)

func RegisterHoldingRoutes(app *fiber.App) {
	holdingRoutes := app.Group("/holdings")
	holdingRoutes.Use(middleware.JWTMiddleware)

	holdingRoutes.Get("/:userID", controller.GetHoldingByID)
	holdingRoutes.Get("/balance/:userID", controller.GetBalanceByID)
}
