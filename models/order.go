package models

import (
	"time"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type Order struct {
	gorm.Model
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	BuyerID    uuid.UUID `json:"buyer_id"`
	ProductID  uuid.UUID `json:"product_id"`
	PaymentID  uuid.UUID `json:"payment_id"`
	ShippingID uuid.UUID `json:"shipping_id"`
	Status     string    `json:"status"`
	Commission  float64   `json:"commission"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Buyer    User
	Product  Product
	Payment  Payment
	Shipping Shipping
}

type OrderRepository interface {
	Create(order *Order) error
	GetByID(id uuid.UUID) (*Order, error)
	GetByBuyerID(buyerID uuid.UUID) ([]Order, error)
	// GetBySellerId(buyerID uuid.UUID) ([]Order, error)
	GetByProductID(productID uuid.UUID) ([]Order, error)
	GetByPaymentID(paymentID uuid.UUID) ([]Order, error)
	GetByShippingID(shippingID uuid.UUID) ([]Order, error)
	Update(order *Order) error
	Delete(id uuid.UUID) error
}

type OrderModel struct {
	DB *gorm.DB
}

func (om *OrderModel) Create(order *Order) error {
	return om.DB.Create(order).Error
}

func (om *OrderModel) GetByID(id uuid.UUID) (*Order, error) {
	var order Order
	if err := om.DB.Preload("Buyer").Preload("Product").Preload("Payment").Preload("Shipping").First(&order, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (om *OrderModel) GetByBuyerID(buyerID uuid.UUID) ([]Order, error) {
	var orders []Order
	if err := om.DB.Preload("Buyer").Preload("Product").Preload("Payment").Preload("Shipping").Where("buyer_id = ?", buyerID).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (om *OrderModel) GetByProductID(productID uuid.UUID) ([]Order, error) {
	var orders []Order
	if err := om.DB.Preload("Buyer").Preload("Product").Preload("Payment").Preload("Shipping").Where("product_id = ?", productID).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (om *OrderModel) GetByPaymentID(paymentID uuid.UUID) ([]Order, error) {
	var orders []Order
	if err := om.DB.Preload("Buyer").Preload("Product").Preload("Payment").Preload("Shipping").Where("payment_id = ?", paymentID).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (om *OrderModel) GetByShippingID(shippingID uuid.UUID) ([]Order, error) {
	var orders []Order
	if err := om.DB.Preload("Buyer").Preload("Product").Preload("Payment").Preload("Shipping").Where("shipping_id = ?", shippingID).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (om *OrderModel) Update(order *Order) error {
	return om.DB.Save(order).Error
}

func (om *OrderModel) Delete(id uuid.UUID) error {
	var order Order
	if err := om.DB.First(&order, "id = ?", id).Error; err != nil {
		return err
	}
	return om.DB.Delete(&order).Error
}