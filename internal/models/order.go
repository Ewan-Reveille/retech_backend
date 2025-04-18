package models

import (
	"time"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type Order struct {
	gorm.Model
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	BuyerID    uuid.UUID
	ProductID  uuid.UUID
	PaymentID  uuid.UUID
	ShippingID uuid.UUID
	Status     string // pending, paid, shipped, received, canceled, dispute

	CreatedAt time.Time
	UpdatedAt time.Time

	Buyer    User
	Product  Product
	Payment  Payment
	Shipping Shipping
}
