package models

import (
	"time"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type User struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Username    string    `gorm:"uniqueIndex"`
	Email       string    `gorm:"uniqueIndex"`
	Password    string
	PhoneNumber string
	Address     string
	IsVerified  bool
	IsSeller    bool
	IsAdmin     bool
	Rating      float64

	CreatedAt time.Time
	UpdatedAt time.Time

	Products         []Product `gorm:"foreignKey:SellerID"`
	Orders           []Order   `gorm:"foreignKey:BuyerID"`
	Reviews          []Review  `gorm:"foreignKey:ReviewerID"`
	MessagesSent     []Message `gorm:"foreignKey:SenderID"`
	MessagesReceived []Message `gorm:"foreignKey:ReceiverID"`
	Reports          []Report  `gorm:"foreignKey:ReporterID"`
}
