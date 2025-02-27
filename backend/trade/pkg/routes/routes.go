// holding/pkg/routes/routes.go

/***
Routes Package

Registers API endpoints and applies authentication middleware.

Features:
1. **RegisterTradeRoutes**:
   - Defines and registers routes related to trades.
   - Secures endpoints using JWT authentication middleware.

Security Considerations:
- Middleware ensures only authenticated users can access trade data.

Endpoints:
- **GET /trade/{tradeID}** → Retrieve a specific trade by its ID.
- **GET /trade/user/{userID}** → Retrieve all trades for a user (supports cursor and limit).
- **GET /trade/event/{eventID}** → Retrieve all trades for a specific event (supports cursor and limit).
- **GET /trade** → Retrieve trades filtered by eventID and userID (supports cursor and limit, requires query params).

Author: Vansh Patel (remotevansh@gmail.com)  
Last Updated: February 27, 2025  
***/

package routes

import (
	"github.com/gofiber/fiber/v2"
	"holding/pkg/controller"
	"holding/pkg/middleware"
)

func RegisterTradeRoutes(app *fiber.App) {
	tradeRoutes := app.Group("/trade")
	tradeRoutes.Use(middleware.JWTMiddleware)

	tradeRoutes.Get("/:tradeID", controller.GetTradeByID)
	tradeRoutes.Get("/user/:userID", controller.GetTradesByUserID)
	tradeRoutes.Get("/event/:eventID", controller.GetTradesByEventID)
	tradeRoutes.Get("/", controller.GetTradesByEventIDUserID)
}
