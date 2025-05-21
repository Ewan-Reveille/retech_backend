package models

import (
	"time"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type PaymentMethod string
type PaymentStatus string

const (
	MethodPaypal PaymentMethod = "paypal"
	MethodCard   PaymentMethod = "card"
	MethodCrypto PaymentMethod = "crypto"
)

const (
	StatusPending   PaymentStatus = "pending"
	StatusChecked PaymentStatus = "checked"
	StatusFailed    PaymentStatus = "failed"
)

type Payment struct {
	gorm.Model
	// PaymentID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	PaymentID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID   uuid.UUID
	Amount    float64
	OrderID   uuid.UUID
	Method    PaymentMethod
	Status    PaymentStatus
	CreatedAt time.Time
}

type PaymentRepository interface {
	Create(payment *Payment) error
	GetByID(id uuid.UUID) (*Payment, error)
	GetByOrderID(orderID uuid.UUID) (*Payment, error)
	Update(payment *Payment) error
	Delete(id uuid.UUID) error
}

type PaymentModel struct {
	DB *gorm.DB
}

func (pm *PaymentModel) Create(payment *Payment) error {
	return pm.DB.Create(payment).Error
}

func (pm *PaymentModel) GetByID(id uuid.UUID) (*Payment, error) {
	var payment Payment
	if err := pm.DB.Preload("Order").First(&payment, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (pm *PaymentModel) GetByOrderID(orderID uuid.UUID) (*Payment, error) {
	var payment Payment
	if err := pm.DB.Preload("Order").First(&payment, "order_id = ?", orderID).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (pm *PaymentModel) GetByUserID(userID uuid.UUID) ([]Payment, error) {
	var payments []Payment
	if err := pm.DB.Where("user_id = ?", userID).Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

func (pm *PaymentModel) Update(payment *Payment) error {
	return pm.DB.Save(payment).Error
}

func (pm *PaymentModel) Delete(id uuid.UUID) error {
	return pm.DB.Delete(&Payment{}, id).Error
}

