package services

import (
	"errors"
	"github.com/stripe/stripe-go/v81"
    "github.com/stripe/stripe-go/v81/product"
    "github.com/stripe/stripe-go/v81/price"
	"github.com/Ewan-Reveille/retech/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductService struct {
	Repo models.ProductRepository
	DB   *gorm.DB
}

func (ps *ProductService) Create(p *models.Product) error {
    // Basic validation
    if p.Title == "" || p.Price <= 0 {
        return errors.New("invalid title or price")
    }

    // Start transaction
    return ps.DB.Transaction(func(tx *gorm.DB) error {
        // 1) Create Stripe product
        spParams := &stripe.ProductParams{
            Name:        stripe.String(p.Title),
            Description: stripe.String(p.Description),
        }
        stripeProd, err := product.New(spParams)
        if err != nil {
            return err
        }

        // 2) Create Stripe price (in cents)
        priceParams := &stripe.PriceParams{
            UnitAmount: stripe.Int64(int64(p.Price * 100)),
            Currency:   stripe.String("eur"),
            Product:    stripe.String(stripeProd.ID),
        }
        stripePriceObj, err := price.New(priceParams)
        if err != nil {
            return err
        }

        // 3) Fill in Stripe IDs
        p.StripeProductID = stripeProd.ID
        p.StripePriceID   = stripePriceObj.ID

        // 4) Persist your own Product
        if err := ps.Repo.Create(p); err != nil {
            return err
        }

        return nil
    })
}


func (ps *ProductService) GetByID(id uuid.UUID) (*models.Product, error) {
	return ps.Repo.GetByID(id)
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
