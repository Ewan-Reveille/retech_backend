package services

import (
	"errors"
	"github.com/Ewan-Reveille/retech/models"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/price"
	"github.com/stripe/stripe-go/v81/product"
	"gorm.io/gorm"

	// "log"
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
func (ps *ProductService) Create(p *models.Product) error {
    if p.Title == "" || p.Price <= 0 {
        return errors.New("invalid title or price")
    }

    return ps.DB.Transaction(func(tx *gorm.DB) error {
        // Use the injected Stripe client
        stripeProd, err := ps.StripeClient.NewProduct(&stripe.ProductParams{
            Name:        stripe.String(p.Title),
            Description: stripe.String(p.Description),
        })
        if err != nil {
            return err
        }

        priceParams := &stripe.PriceParams{
            UnitAmount: stripe.Int64(int64(p.Price * 100)),
            Currency:   stripe.String("eur"),
            Product:    stripe.String(stripeProd.ID),
        }
        
        stripePriceObj, err := ps.StripeClient.NewPrice(priceParams)
        if err != nil {
            return err
        }

        p.StripeProductID = stripeProd.ID
        p.StripePriceID = stripePriceObj.ID

        return ps.Repo.Create(p)
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
