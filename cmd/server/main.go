package main

import (
	"log"
	"os"

	"github.com/Ewan-Reveille/retech/config"
	"github.com/Ewan-Reveille/retech/db"
	"github.com/Ewan-Reveille/retech/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v81"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	stripeKey := os.Getenv("STRIPE_SECRET_KEY")
	if stripeKey == "" {
		log.Fatal("❌ STRIPE_SECRET_KEY n'est pas défini")
	}
	stripe.Key = stripeKey

	// Setup DB connection
	DB, err := config.SetupDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	app := fiber.New()
	// Run migrations
	// db.Migrate(DB)
	routes.RegisterProductRoutes(app, DB)
	routes.RegisterUserRoutes(app, DB)
	// Connect to DB (assuming db.Connect() is for keeping the DB connection alive)
	db.Connect()


	// Enable CORS middleware
	app.Use(cors.New())

	// Define a route for testing
	app.Get("/", func(c *fiber.Ctx) error {
		log.Println("ReTech API status: launched 🚀")
		return c.SendString("ReTech API status: launched 🚀")
	})

	// Start the server on the specified port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to port 8080 if not provided
	}
	log.Fatal(app.Listen(":" + port))
}
