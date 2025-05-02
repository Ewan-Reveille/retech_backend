package services

import (
	"errors"

	"github.com/Ewan-Reveille/retech/models"
	"github.com/google/uuid"
)

type MessageService struct {
	Repo models.MessageRepository
}

func (ms *MessageService) Create(message *models.Message) error {
	if message.Content == "" {
		return errors.New("message content cannot be empty")
	}

	if message.SenderID == uuid.Nil || message.ReceiverID == uuid.Nil {
		return errors.New("sender and receiver IDs cannot be empty")
	}

	return ms.Repo.Create(message)
}

func (ms *MessageService) GetByID(id uuid.UUID) (*models.Message, error) {
	message, err := ms.Repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (ms *MessageService) GetBySenderID(senderID uuid.UUID) ([]models.Message, error) {
	messages, err := ms.Repo.GetBySenderID(senderID)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (ms *MessageService) GetByReceiverID(receiverID uuid.UUID) ([]models.Message, error) {
	messages, err := ms.Repo.GetByReceiverID(receiverID)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (ms *MessageService) Update(message *models.Message) error {

	existing, err := ms.Repo.GetByID(message.ID)
	if err != nil {
		return err
	}

	if existing.Content == "" {
		return errors.New("message content cannot be empty")
	}
	if existing.ID != message.ID {
		return errors.New("message ID mismatch")
	}
	return ms.Repo.Update(message)
}

func (ms *MessageService) Delete(id uuid.UUID) error {
	message, err := ms.Repo.GetByID(id)
	if err != nil {
		return err
	}

	if message.ID != id {
		return errors.New("message ID mismatch")
	}
	return ms.Repo.Delete(id)
}
