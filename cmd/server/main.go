// @title           Retech API
// @version         1.0
// @description     This is the backend API for Retech
// @termsOfService  http://swagger.io/terms/

// @contact.name   Ewan Réveillé--Adam
// @contact.url    https://aia-handicap.com
// @contact.email  ewan@aia-handicap.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
package main

import (
    "log"
    "os"

    "github.com/Ewan-Reveille/retech/config"
    "github.com/Ewan-Reveille/retech/db"
    "github.com/Ewan-Reveille/retech/routes"
    "github.com/Ewan-Reveille/retech/services"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    swagger "github.com/arsmn/fiber-swagger/v2"
    "github.com/joho/godotenv"
    "github.com/stripe/stripe-go/v81"
    _ "github.com/Ewan-Reveille/retech/docs"
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

	realStripeClient := &services.StripeClientImpl{}
	// Run migrations
	// db.Migrate(DB)
	routes.RegisterProductRoutes(app, DB, realStripeClient)
	routes.RegisterUserRoutes(app, DB)
	routes.RegisterCategoryRoutes(app, DB)
	routes.RegisterOrderRoutes(app, DB)
	app.Get("/swagger/*", swagger.HandlerDefault)
	// Connect to DB (assuming db.Connect() is for keeping the DB connection alive)
	db.Connect()

	// Enable CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: false,
	}))

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
	for _, r := range app.GetRoutes() {
		log.Printf("Route registered: %s %s", r.Method, r.Path)
	}
	log.Fatal(app.Listen(":" + port))
}
