package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model    `swaggerignore:"true"`
	// ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name string    `gorm:"unique"`

	Products []Product
}

// @swagger:model
type SwaggerCategory struct {
    ID   uuid.UUID `json:"id"`
    Name string    `json:"name"`
}

type CategoryRepository interface {
	Create(category *Category) error
	GetByID(id uuid.UUID) (*Category, error)
	GetByName(name string) (*Category, error)
	GetAll() ([]Category, error)
	Update(category *Category) error
	Delete(id uuid.UUID) error
}

type CategoryModel struct {
	DB *gorm.DB
}

func (cm *CategoryModel) Create(category *Category) error {
	return cm.DB.Create(category).Error
}

// func (cm *CategoryModel) GetByID(id uuid.UUID) (*Category, error) {
// 	var category Category
// 	if err := cm.DB.Preload("Products").First(&category, "id = ?", id).Error; err != nil {
// 		return nil, err
// 	}
// 	return &category, nil
// }

func (cm *CategoryModel) GetByID(id uuid.UUID) (*Category, error) {
    var category Category
    if err := cm.DB.First(&category, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return &category, nil
}

func (cm *CategoryModel) GetByName(name string) (*Category, error) {
	var category Category
	if err := cm.DB.First(&category, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (cm *CategoryModel) GetAll() ([]Category, error) {
	var categories []Category
	if err := cm.DB.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (cm *CategoryModel) Update(category *Category) error {
	return cm.DB.Save(category).Error
}

func (cm *CategoryModel) Delete(id uuid.UUID) error {
	return cm.DB.Delete(&Category{}, "id = ?", id).Error
}
