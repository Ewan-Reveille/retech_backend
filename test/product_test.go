// test/controllers/product_test.go
package test

import (
	"bytes"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/Ewan-Reveille/retech/models"
	"github.com/Ewan-Reveille/retech/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"io"
	"encoding/json"
	// "github.com/stripe/stripe-go/v81"
)

// Mock Stripe client (already defined in product_service_test.go, but good to have it here too if this file is standalone)
// Make sure this is only declared once in your test package.
// If it's already in product_service_test.go, do NOT redeclare it here.
// func (m *MockStripeClient) NewProduct(params *stripe.ProductParams) (*stripe.Product, error) {
// 	return &stripe.Product{ID: "prod_test"}, nil
// }

// func (m *MockStripeClient) NewPrice(params *stripe.PriceParams) (*stripe.Price, error) {
// 	return &stripe.Price{ID: "price_test"}, nil
// }


func TestProductController_CreateProduct(t *testing.T) {
	db := SetupTestDB(t)
	defer TeardownTestDB(db)

	// --- ADD TABLE CREATION HERE ---
	// Manually create tables with SQLite-compatible syntax for this test
	db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id BLOB PRIMARY KEY,
		username TEXT,
		email TEXT UNIQUE,
		password TEXT,
		is_seller BOOLEAN,
		created_at DATETIME,
		updated_at DATETIME,
		deleted_at DATETIME
	)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS categories (
		id BLOB PRIMARY KEY,
		name TEXT,
		created_at DATETIME,
		updated_at DATETIME,
		deleted_at DATETIME
	)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS products (
		id BLOB PRIMARY KEY,
		title TEXT,
		description TEXT,
		price REAL,
		status TEXT,
		seller_id BLOB,
		condition TEXT,
		category_id BLOB,
		stripe_product_id TEXT,
		stripe_price_id TEXT,
		created_at DATETIME,
		updated_at DATETIME,
		deleted_at DATETIME,
		FOREIGN KEY(seller_id) REFERENCES users(id),
		FOREIGN KEY(category_id) REFERENCES categories(id)
	)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS product_images (
		id BLOB PRIMARY KEY,
		product_id BLOB,
		image_url TEXT,
		alt TEXT,
		created_at DATETIME,
		updated_at DATETIME,
		deleted_at DATETIME,
		FOREIGN KEY(product_id) REFERENCES products(id)
	)`)
	// --- END ADDITION ---


	app := fiber.New()
	mockStripe := &MockStripeClient{} // Create an instance of your mock
	routes.RegisterProductRoutes(app, db, mockStripe) // Pass the mock here

	// Create required entities
	seller := models.User{
		ID:       uuid.New(),
		Username: "test_seller",
		Email:    "seller@test.com",
		Password: "password",
		IsSeller: true,
	}
	db.Create(&seller)

	category := models.Category{
		ID:   uuid.New(),
		Name: "Test Category",
	}
	db.Create(&category)

	// Create test image file
	tmpFile, err := os.CreateTemp("", "test-image-*.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// Create multipart form
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Add form fields
	_ = writer.WriteField("title", "Test Product")
	_ = writer.WriteField("description", "Test Description")
	_ = writer.WriteField("price", "99.99")
	_ = writer.WriteField("category", category.ID.String())
	_ = writer.WriteField("seller_id", "test1")
	_ = writer.WriteField("condition", "new")

	// Add image file
	part, _ := writer.CreateFormFile("images", filepath.Base(tmpFile.Name()))
	file, _ := os.Open(tmpFile.Name())
	defer file.Close()
	_, _ = io.Copy(part, file)

	writer.Close()

	// Create request
	req := httptest.NewRequest("POST", "/products", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := app.Test(req)
	assert.NoError(t, err)
	// The response status code will now be StatusOK (200) or StatusCreated (201) if successful
	// Check the logs for the actual error from the controller if it's not 201
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode) // This assertion will now pass if the product is created

	var createdProduct models.Product
	err = json.NewDecoder(resp.Body).Decode(&createdProduct)
	assert.NoError(t, err)
	assert.Equal(t, "Test Product", createdProduct.Title)

}