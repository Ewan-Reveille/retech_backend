// test/testutils.go
package test

import (
	"testing"

	"github.com/Ewan-Reveille/retech/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// test/testutils.go
func SetupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Disable automatic UUID generation at DB level
	err = db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Category{},
		&models.Order{},
		&models.Payment{},
	)

	db.Exec(`CREATE TABLE IF NOT EXISTS users (
        id BLOB PRIMARY KEY,
        created_at DATETIME,
        updated_at DATETIME,
        deleted_at DATETIME,
        username TEXT UNIQUE,
        email TEXT UNIQUE,
        password TEXT
    )`)

	// Create categories table
	db.Exec(`CREATE TABLE IF NOT EXISTS categories (
        id BLOB PRIMARY KEY,
        name TEXT UNIQUE
    )`)
	// ..

	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func TeardownTestDB(db *gorm.DB) {
	sqlDB, _ := db.DB()
	sqlDB.Close()
}
