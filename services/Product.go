package services

import (
	"errors"
	"github.com/Ewan-Reveille/retech/models"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/price"
	"github.com/stripe/stripe-go/v81/product"
	"gorm.io/gorm"

	"log"
)
type StripeClientImpl struct{}

func (s *StripeClientImpl) NewProduct(params *stripe.ProductParams) (*stripe.Product, error) {
	return product.New(params)
}

// NewPrice calls the Stripe API to create a new price.
func (s *StripeClientImpl) NewPrice(params *stripe.PriceParams) (*stripe.Price, error) {
	return price.New(params)
}
type StripeClient interface {
	NewProduct(params *stripe.ProductParams) (*stripe.Product, error)
	NewPrice(params *stripe.PriceParams) (*stripe.Price, error)
}

type ProductService struct {
	Repo         models.ProductRepository
	DB           *gorm.DB
	StripeClient StripeClient
}

// services/Product.go
// services/Product.go
func (s *ProductService) Create(product *models.Product) error {
	log.Println("[ProductService.Create] Starting transaction")
	
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Create Stripe Product
		log.Printf("[ProductService.Create] Creating Stripe product, service: %+v", s)
		if s.StripeClient == nil {
			log.Fatal("[ProductService.Create] StripeClient is nil!")
		}
		
		stripeProduct, err := s.StripeClient.NewProduct(&stripe.ProductParams{
			Name: stripe.String(product.Title),
		})
		if err != nil {
			log.Printf("[ProductService.Create] Stripe product error: %v", err)
			return err
		}
		log.Printf("[ProductService.Create] Created Stripe product: %s", stripeProduct.ID)

		// 2. Create Stripe Price
		log.Printf("[ProductService.Create] Creating Stripe price for product %s", stripeProduct.ID)
		stripePrice, err := s.StripeClient.NewPrice(&stripe.PriceParams{
			Product:    stripe.String(stripeProduct.ID),
			UnitAmount: stripe.Int64(int64(product.Price * 100)),
			Currency:   stripe.String("eur"),
		})
		if err != nil {
			log.Printf("[ProductService.Create] Stripe price error: %v", err)
			return err
		}
		log.Printf("[ProductService.Create] Created Stripe price: %s", stripePrice.ID)

		// 3. Save to DB
		log.Println("[ProductService.Create] Saving product to database")
		product.StripeProductID = stripeProduct.ID
		product.StripePriceID = stripePrice.ID
		
		if err := tx.Create(product).Error; err != nil {
			log.Printf("[ProductService.Create] Database save error: %v", err)
			return err
		}
		log.Printf("[ProductService.Create] Product saved successfully. ID: %s", product.ID)
		
		return nil
	})
}

func (ps *ProductService) GetByID(id uuid.UUID) (*models.Product, error) {
	return ps.Repo.GetByID(id)
}

func (ps *ProductService) GetAll() ([]models.Product, error) {
	products, err := ps.Repo.GetAll()
	if err != nil {
		return nil, err
	}

	for i := range products {
		// Fetch Stripe product details
		stripeProduct, err := product.Get(products[i].StripeProductID, nil)
		if err != nil {
			return nil, err
		}
		products[i].Title = stripeProduct.Name
		products[i].Description = stripeProduct.Description
	}

	return products, nil
}

func (ps *ProductService) Update(product *models.Product) error {
	existing, err := ps.Repo.GetByID(product.ID)
	if err != nil {
		return err
	}

	if existing.Status == "sold" {
		return errors.New("cannot update a sold product")
	}

	return ps.Repo.Update(product)
}

func (ps *ProductService) Delete(id uuid.UUID, userID uuid.UUID) error {
	product, err := ps.Repo.GetByID(id)
	if err != nil {
		return err
	}

	if product.SellerID != userID {
		return errors.New("unauthorized")
	}

	return ps.Repo.Delete(id)
}
