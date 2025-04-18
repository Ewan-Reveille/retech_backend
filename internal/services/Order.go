package services

import (
	"errors"
	"github.com/Ewan-Reveille/retech/internal/models"
	"github.com/google/uuid"
)

type OrderService struct {
	Repo models.OrderRepository
}

func (os *OrderService) Create(order *models.Order) error {
	if order.BuyerID == uuid.Nil {
		return errors.New("buyer ID cannot be empty")
	}
	if order.ProductID == uuid.Nil {
		return errors.New("product ID cannot be empty")
	}
	
	order.ID = uuid.New() // ⬅️ Remplace utils.GenerateUUID()
	return os.Repo.Create(order)
}

func (os *OrderService) GetByID(id uuid.UUID) (*models.Order, error) {
	order, err := os.Repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (os *OrderService) GetByBuyerID(buyerId uuid.UUID) ([]models.Order, error) {
	orders, err := os.Repo.GetByBuyerID(buyerId)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (os *OrderService) GetBySellerId(sellerId uuid.UUID) ([]models.Order, error) {
	orders, err := os.Repo.GetBySellerId(sellerId)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (os *OrderService) GetByProductID(productID uuid.UUID) ([]models.Order, error) {
	orders, err := os.Repo.GetByProductID(productID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (os *OrderService) Update(order *models.Order) error {
	existing, err := os.Repo.GetByID(order.ID)
	if err != nil {
		return err
	}

	if existing.BuyerID != order.BuyerID {
		return errors.New("buyer ID mismatch")
	}
	if existing.ProductID != order.ProductID {
		return errors.New("product ID mismatch")
	}

	return os.Repo.Update(order)
}


func (os *OrderService) Delete(id uuid.UUID) error {
	order, err := os.Repo.GetByID(id)
	if err != nil {
		return err
	}

	if order.BuyerID == uuid.Nil {
		return errors.New("buyer ID cannot be empty")
	}
	if order.ProductID == uuid.Nil {
		return errors.New("product ID cannot be empty")
	}

	if order.ProductID == uuid.Nil {
		return errors.New("product ID cannot be empty")
	}
	return os.Repo.Delete(id)
}