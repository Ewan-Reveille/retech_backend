package models

import (
	"time"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type ProductImage struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ProductID uuid.UUID
	ImageURL  string
	Alt       string
	CreatedAt time.Time
}
