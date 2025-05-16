package test

import (
	"github.com/Ewan-Reveille/retech/models"
	"github.com/Ewan-Reveille/retech/routes"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

func SetupApp() *fiber.App {
	// Base SQLite en mémoire pour les tests
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to the test database: %v", err)
	}

	// Dans SetupApp
	db.Exec("DELETE FROM users")

	// 🛠️ Exécute la migration
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("failed to migrate test database: %v", err)
	}

	app := fiber.New()

	// Ajoute les routes avec la DB
	routes.RegisterUserRoutes(app, db)

	return app
}
