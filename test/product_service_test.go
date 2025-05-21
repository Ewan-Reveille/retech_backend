// test/services/product_service_test.go
package test

import (
	"testing"

	"github.com/Ewan-Reveille/retech/models"
	"github.com/Ewan-Reveille/retech/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stripe/stripe-go/v81"
)

// Mock Stripe client
type MockStripeClient struct{}

func (m *MockStripeClient) NewProduct(params *stripe.ProductParams) (*stripe.Product, error) {
	return &stripe.Product{ID: "prod_test"}, nil
}

func (m *MockStripeClient) NewPrice(params *stripe.PriceParams) (*stripe.Price, error) {
	return &stripe.Price{ID: "price_test"}, nil
}

func TestProductService_CreateWithStripe(t *testing.T) {
	db := SetupTestDB(t)
	defer TeardownTestDB(db)

	// Manually create tables with SQLite-compatible syntax
	db.Exec(`CREATE TABLE IF NOT EXISTS products (
    id BLOB PRIMARY KEY,
    title TEXT,
    description TEXT,
    price REAL,
    seller_id BLOB,
    stripe_product_id TEXT,
    stripe_price_id TEXT,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME
)`)

	productModel := models.ProductModel{DB: db}
	mockStripe := &MockStripeClient{}

	productService := services.ProductService{
		Repo:         &productModel,
		DB:           db,
		StripeClient: mockStripe, // Inject the mock
	}

	product := models.Product{
		Title:       "Test",
		Description: "Test",
		Price:       100.0,
		SellerID:    uuid.New(),
		CategoryID:  uuid.New(),
	}

	err := productService.Create(&product)
	assert.NoError(t, err)
	assert.Equal(t, "prod_test", product.StripeProductID)
	assert.Equal(t, "price_test", product.StripePriceID)
}
