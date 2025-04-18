package models

import (
	"time"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type Promotion struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ProductID uuid.UUID
	Type      string // highlighting, boost
	StartDate time.Time
	EndDate   time.Time
	CreatedAt time.Time
}
