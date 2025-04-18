package services

import (
    "errors"
    "github.com/Ewan-Reveille/retech/internal/models"
	"github.com/google/uuid"
)

type PaymentService struct {
	Repo models.PaymentRepository
}

func (ps *PaymentService) Create(payment *models.Payment) error {

	// Validate payment details
	if payment.UserID == uuid.Nil {
		return errors.New("invalid user ID")
	}

	if payment.Amount <= 0 {
		return errors.New("invalid payment amount")
	}
	
}