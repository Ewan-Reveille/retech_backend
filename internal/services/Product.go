package services

import (
    "errors"
    "github.com/Ewan-Reveille/retech/internal/models"
	"github.com/google/uuid"
)

type ProductService struct {
    Repo models.ProductRepository
}

func (ps *ProductService) Create(product *models.Product) error {
    if product.Price <= 0 {
        return errors.New("price must be positive")
    }

    return ps.Repo.Create(product)
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
