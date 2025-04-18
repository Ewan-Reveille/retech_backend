package models

import (
	"time"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type Payment struct {
	gorm.Model
	ID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	OrderID uuid.UUID
	Method  string // enum: card, crypto, etc.
	Status  string // enum: checked, failed

	CreatedAt time.Time
}
