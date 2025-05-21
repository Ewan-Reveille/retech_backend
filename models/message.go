package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	// // ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	SenderID   uuid.UUID
	ReceiverID uuid.UUID
	Content    string
	ReadAt     *time.Time
	CreatedAt  time.Time

	Sender   User `gorm:"foreignKey:SenderID"`
	Receiver User `gorm:"foreignKey:ReceiverID"`
}

type MessageRepository interface {
	Create(message *Message) error
	GetByID(id uuid.UUID) (*Message, error)
	GetBySenderID(senderID uuid.UUID) ([]Message, error)
	GetByReceiverID(receiverID uuid.UUID) ([]Message, error)
	Update(message *Message) error
	Delete(id uuid.UUID) error
}

type MessageModel struct {
	DB *gorm.DB
}

func (mm *MessageModel) Create(message *Message) error {
	return mm.DB.Create(message).Error
}

func (mm *MessageModel) GetByID(id uuid.UUID) (*Message, error) {
	var message Message
	if err := mm.DB.Preload("Sender").Preload("Receiver").First(&message, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &message, nil
}

func (mm *MessageModel) GetBySenderID(senderID uuid.UUID) ([]Message, error) {
	var messages []Message
	if err := mm.DB.Preload("Sender").Preload("Receiver").Where("sender_id = ?", senderID).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (mm *MessageModel) GetByReceiverID(receiverID uuid.UUID) ([]Message, error) {
	var messages []Message
	if err := mm.DB.Preload("Sender").Preload("Receiver").Where("receiver_id = ?", receiverID).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (mm *MessageModel) Update(message *Message) error {
	return mm.DB.Save(message).Error
}

func (mm *MessageModel) Delete(id uuid.UUID) error {
	var message Message
	if err := mm.DB.First(&message, "id = ?", id).Error; err != nil {
		return err
	}
	return mm.DB.Delete(&message).Error
}
