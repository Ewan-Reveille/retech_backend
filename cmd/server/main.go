package main

import (
    "log"
    "os"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/joho/godotenv"
    "github.com/Ewan-Reveille/retech/internal/db"
    "github.com/Ewan-Reveille/retech/internal/config"
	"github.com/Ewan-Reveille/retech/internal/routes"
)

func main() {
    // Load environment variables
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Setup DB connection
    DB, err := config.SetupDB()
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Run migrations
    db.Migrate(DB)
    
    // Connect to DB (assuming db.Connect() is for keeping the DB connection alive)
    db.Connect()

    // Initialize Fiber app
    app := fiber.New()

    // Enable CORS middleware
    app.Use(cors.New())

    // Define a route for testing
    app.Get("/", func(c *fiber.Ctx) error {
		log.Println("ReTech API status: launched 🚀")
        return c.SendString("ReTech API status: launched 🚀")
    })

	routes.RegisterProductRoutes(app, DB)

    // Start the server on the specified port
    port := os.Getenv("PORT")
    if port == "" {
        port = "3000" // Default to port 3000 if not provided
    }
    log.Fatal(app.Listen(":" + port))
}
