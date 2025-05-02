package services

import (
	"errors"

	"github.com/Ewan-Reveille/retech/models"
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
	return ps.Repo.Create(payment)
}

func (ps *PaymentService) GetByID(id uuid.UUID) (*models.Payment, error) {
	payment, err := ps.Repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (ps *PaymentService) GetByOrderID(orderID uuid.UUID) (*models.Payment, error) {
	payment, err := ps.Repo.GetByOrderID(orderID)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (ps *PaymentService) Update(payment *models.Payment) error {
	existing, err := ps.Repo.GetByID(payment.PaymentID)
	if err != nil {
		return err
	}

	if existing.Status == "completed" {
		return errors.New("cannot update a completed payment")
	}

	return ps.Repo.Update(payment)
}

func (ps *PaymentService) Delete(id uuid.UUID) error {
	payment, err := ps.Repo.GetByID(id)
	if err != nil {
		return err
	}

	if payment.Status == "completed" {
		return errors.New("cannot delete a completed payment")
	}

	return ps.Repo.Delete(id)
}

func (ps *PaymentService) GetByUserID() ([]models.Payment, error) {
	// Implement logic to get payments by user ID
	// This is a placeholder for the actual implementation
	return nil, errors.New("not implemented")
}