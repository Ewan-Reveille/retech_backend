package db

import (
	"github.com/Ewan-Reveille/retech/models"
	"gorm.io/gorm"
)

func Migrate(DB *gorm.DB) {
	// Migrate all models at once
	err := DB.AutoMigrate(
		&models.User{},
		// &models.Category{},
		// &models.Message{},
		// &models.Order{},
		// &models.Payment{},
		// &models.ProductCondition{},
		&models.Product{},
		// &models.ProductImage{},
		// &models.Promotion{},
		// &models.Report{},
		// &models.Review{},
		// &models.Shipping{},
	)

	if err != nil {
		panic("migration failed: " + err.Error())
	}
}
