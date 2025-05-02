package models

import (
	"time"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type Shipping struct {
	gorm.Model
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Carrier        string
	TrackingNumber string
	Status         string // prepared, shipped, transit, delivered

	CreatedAt time.Time
}
