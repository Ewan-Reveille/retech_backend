package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Review struct {
	gorm.Model
	// ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	ReviewerID uuid.UUID
	ReviewedID uuid.UUID
	Rating     int
	Comment    string
	CreatedAt  time.Time

	Reviewer User `gorm:"foreignKey:ReviewerID"`
	Reviewed User `gorm:"foreignKey:ReviewedID"`
}
