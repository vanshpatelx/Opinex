package main

import (
	"log"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"holding/pkg/config"
	"holding/pkg/healthcheck"
	"holding/pkg/routes"
)

func main() {
    config.LoadConfig()

    // Check service readiness before proceeding
    healthcheck.CheckAllServices()

    app := fiber.New(fiber.Config{
        JSONEncoder: json.Marshal,
        JSONDecoder: json.Unmarshal,
    })

    // Middleware
    app.Use(logger.New()) // Log requests
    app.Use(cors.New())   // Enable CORS


    app.Get("/", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "message": "🚀 Holding Service is Running!",
        })
    })

        
    app.Get("/holding/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"success": true, "message": " 🚀 🚀Server is running."})
    })

    routes.GetBalanceRoutes(app)
    routes.UpdateBalanceRoutes(app)

    log.Println("✅ Server is starting on port", config.AppConfig.Port)
    app.Listen(":" + config.AppConfig.Port)
}
