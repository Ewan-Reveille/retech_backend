package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductCondition struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string    // new, very good, used, etc.
	Description string
}
