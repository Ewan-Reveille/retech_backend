package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/product"
	"github.com/stripe/stripe-go/v81/price"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Load environment variables from a .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize Stripe with your secret key
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	log.Println("Stripe Secret Key:", stripe.Key)

	// Create a new Fiber application
	app := fiber.New()
	app.Use(cors.New())

	// Endpoint to add a product
	app.Post("/api/products", func(c *fiber.Ctx) error {
		// Struct to parse JSON data from the request
		var reqData struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Price       int64  `json:"price"`
			ImageUrl    string `json:"imageUrl"`
		}

		// Parse the JSON data from the request body
		if err := c.BodyParser(&reqData); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		// Create a product on Stripe
		productParams := &stripe.ProductParams{
			Name:        stripe.String(reqData.Name),
			Description: stripe.String(reqData.Description),
			Images:      stripe.StringSlice([]string{reqData.ImageUrl}),
		}

		stripeProduct, err := product.New(productParams)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		// Create a price for this product
		priceParams := &stripe.PriceParams{
			UnitAmount: stripe.Int64(reqData.Price * 100), // Convert to cents
			Currency:   stripe.String("eur"),
			Product:    stripe.String(stripeProduct.ID),
		}

		stripePrice, err := price.New(priceParams)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		// Return the response with product and price information
		return c.JSON(fiber.Map{
			"message":      "Product successfully added to Stripe",
			"stripeProduct": stripeProduct,
			"stripePrice": stripePrice,
		})
	})

	// Start the server on port 8080
	log.Println("🚀 Server running on port 8080")
	log.Fatal(app.Listen(":8080"))
}

