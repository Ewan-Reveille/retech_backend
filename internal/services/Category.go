package services

import (
	"errors"
	"github.com/Ewan-Reveille/retech/internal/models"
	"github.com/google/uuid"
)

type CategoryService struct {
	Repo models.CategoryRepository
}

func (cs *CategoryService) Create(category *models.Category) error {
	if category.Name == "" {
		return errors.New("category name cannot be empty")
	}

	existing, err := cs.Repo.GetByName(category.Name)
	if err == nil && existing != nil {
		return errors.New("category already exists")
	}

	return cs.Repo.Create(category)
}

func (cs *CategoryService) GetByID(id uuid.UUID) (*models.Category, error) {
	category, err := cs.Repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (cs *CategoryService) GetByName(name string) (*models.Category, error) {
	category, err := cs.Repo.GetByName(name)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (cs *CategoryService) GetAll() ([]models.Category, error) {
	categories, err := cs.Repo.GetAll()
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (cs *CategoryService) Update(category *models.Category) error {
	existing, err := cs.Repo.GetByID(category.ID)
	if err != nil {
		return err
	}

	if existing.Name == "" {
		return errors.New("category name cannot be empty")
	}
	if existing.ID != category.ID {
		return errors.New("category ID mismatch")
	}
	return cs.Repo.Update(category)
}

func (cs *CategoryService) Delete(id uuid.UUID) error {
	category, err := cs.Repo.GetByID(id)
	if err != nil {
		return err
	}

	if len(category.Products) > 0 {
		return errors.New("cannot delete category with products")
	}

	return cs.Repo.Delete(id)
}