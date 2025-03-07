// package main

// import (

// 	// "github.com/gofiber/fiber/v3"
// 	// "holding/pkg/cache"
// 	// "time"
// 	// "fmt"

// )

// func main(){
// 		app := fiber.New(fiber.Config{
// 	JSONEncoder: json.Marshal,
// 	JSONDecoder: json.Unmarshal,
// })

// 	app.Get("/", func(c fiber.Ctx) error {
// 		return c.SendString("Hello World")
// 	});

// 	log.Fatal(app.Listen(":3000"))
// }

// Example usage
// func main() {
// 	// Initialize Redis
// 	if err := cache.Init("localhost:6379", ""); err != nil {
// 		fmt.Println("{\"error\": \"", err.Error(), "\"}")
// 		return
// 	}

// 	sampleData := map[string]interface{}{"name": "John Doe", "age": 30}

// 	// Set data in Redis
// 	if err :=cache.Set("user:1001", sampleData, 72000*time.Second); err != nil {
// 		fmt.Println("{\"error\": \"Error setting data:\", \"details\": \"", err.Error(), "\"}")
// 		return
// 	}

// 	// Get data from Redis
// 	data, err := cache.Get("user:1001")
// 	if err != nil {
// 		fmt.Println("{\"error\": \"Error getting data:\", \"details\": \"", err.Error(), "\"}")
// 		return
// 	}

// 	fmt.Print(data)

// 	// Close Redis connection
// 	cache.Close()
// }

// package main

// import (
// 	"fmt"
// 	"log"
// 	"holding/pkg/db"
// )

// func main() {
// 	databaseURL := "postgres://admin:admin123@localhost:5432/mydatabase?sslmode=disable"

// 	// Initialize the database connection
// 	if err := db.Init(databaseURL); err != nil {
// 		log.Fatalf("Database initialization failed: %v", err)
// 	}
// 	defer db.Close()

// 	fmt.Println("Database connection successful!")

// 	// Test FetchOne
// 	query := "SELECT id, name FROM users LIMIT 1;" // Replace with an actual query
// 	result, err := db.FetchOne(query)
// 	if err != nil {
// 		log.Printf("FetchOne error: %v", err)
// 	} else {
// 		fmt.Println("FetchOne result:", result)
// 	}

// 	// Test FetchAll
// 	queryAll := "SELECT id, name FROM users;" // Replace with an actual query
// 	results, err := db.FetchAll(queryAll)
// 	if err != nil {
// 		log.Printf("FetchAll error: %v", err)
// 	} else {
// 		fmt.Println("FetchAll results:", results)
// 	}
// }

package main

import (
	"log"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"holding/pkg/cache"
	"holding/pkg/db"
	"holding/pkg/middleware"
	"holding/pkg/routes"
)

func main() {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

    cache.Init("localhost:6379", "")
    db.Init("postgres://admin:admin123@localhost:5432/mydatabase?sslmode=disable")

	// Middleware
	app.Use(logger.New()) // Log requests
	app.Use(cors.New())   // Enable CORS

    app.Post("/auth/login", func(c *fiber.Ctx) error {
		// Dummy user data (in real-world, verify username & password)
		id := "12334512412334512412"
		role := "admin"

		token, err := middleware.GenerateJWT(id, role)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
		}
		return c.JSON(fiber.Map{"token": token})
	})


	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "🚀 Trade Service is Running!"})
	})

	app.Get("/trade/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true, "message": "Server is running."})
	})

	routes.RegisterTradeRoutes(app)

	log.Fatal(app.Listen(":3005"))
}
