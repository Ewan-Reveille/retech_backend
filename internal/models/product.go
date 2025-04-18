package models

import (
	"time"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type Product struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Title       string
	Description string
	Price       float64
	Status      string // enum: online, sold, deleted, pending

	SellerID uuid.UUID
	Seller   User

	CategoryID uuid.UUID
	Category   Category

	ConditionID uuid.UUID
	Condition   ProductCondition

	CreatedAt time.Time
	UpdatedAt time.Time

	Images     []ProductImage
	Promotions []Promotion
}


type ProductRepository interface {
	Create(product *Product) error
	GetByID(id uuid.UUID) (*Product, error)
	Update(product *Product) error
	Delete(id uuid.UUID) error
}

type ProductModel struct {
	DB *gorm.DB
}

func (pm *ProductModel) Create(product *Product) error {
	return pm.DB.Create(product).Error
}

func (pm *ProductModel) GetByID(id uuid.UUID) (*Product, error) {
	var product Product
	if err := pm.DB.Preload("Seller").Preload("Category").Preload("Condition").Preload("Images").Preload("Promotions").First(&product, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (pm *ProductModel) Update(product *Product) error {
	return pm.DB.Save(product).Error
}

func (pm *ProductModel) Delete(id uuid.UUID) error {
	return pm.DB.Delete(&Product{}, "id = ?", id).Error
}
