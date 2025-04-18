package models

import (
	"time"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type Message struct {
	gorm.Model
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	SenderID   uuid.UUID
	ReceiverID uuid.UUID
	Content    string
	ReadAt     *time.Time
	CreatedAt  time.Time

	Sender   User `gorm:"foreignKey:SenderID"`
	Receiver User `gorm:"foreignKey:ReceiverID"`
}
