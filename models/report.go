package models

import (
	"time"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type Report struct {
	gorm.Model
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ReporterID uuid.UUID
	ReportedID uuid.UUID
	Reason     string
	Status     string // enum: reviewed, ignored, canceled

	CreatedAt time.Time
	UpdatedAt time.Time

	Reporter User `gorm:"foreignKey:ReporterID"`
	Reported User `gorm:"foreignKey:ReportedID"`
}
