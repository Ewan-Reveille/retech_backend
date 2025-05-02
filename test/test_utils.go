package test

import (
	"github.com/Ewan-Reveille/retech/models"
	"github.com/Ewan-Reveille/retech/routes"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupTestApp() (*fiber.App, *gorm.DB) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	app := fiber.New()

	// Migration des tables nécessaires pour le test
	db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Category{},
		&models.ProductCondition{},
		&models.ProductImage{},
		&models.Promotion{},
	)

	routes.RegisterProductRoutes(app.Group("/api"), db)

	return app, db
}
