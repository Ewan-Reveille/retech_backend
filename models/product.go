package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type productionCondition string

const (
	NEW       productionCondition = "new"
	VERY_GOOD productionCondition = "very good"
	GOOD      productionCondition = "good"
	USED      productionCondition = "used"
	FAIR      productionCondition = "fair"
	UNKNOWN   productionCondition = "unknown"
)

type Product struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Status      string    `json:"status,omitempty"`

	SellerID uuid.UUID `json:"seller_id"`
	Seller   User      `json:"-"`

	Condition string `json:"condition"`

	CategoryID uuid.UUID `json:"category_id"`
	Category   Category  `json:"-"`

	StripeProductID string `json:"-"`
	StripePriceID   string `json:"-"`

	Images     []ProductImage `json:"-"`
	Promotions []Promotion    `json:"-"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type ProductRepository interface {
	Create(product *Product) error
	GetByID(id uuid.UUID) (*Product, error)
	GetAll() ([]Product, error)
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
	if err := pm.DB.Preload("Seller").Preload("Category").Preload("Images").Preload("Promotions").First(&product, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (pm *ProductModel) GetAll() ([]Product, error) {
	var products []Product
	if err := pm.DB.Preload("Seller").Preload("Category").Preload("Images").Preload("Promotions").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (pm *ProductModel) Update(product *Product) error {
	return pm.DB.Save(product).Error
}

func (pm *ProductModel) Delete(id uuid.UUID) error {
	return pm.DB.Delete(&Product{}, "id = ?", id).Error
}
