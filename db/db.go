package db

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// Charge .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Failed to load .env file")
	}

	// dsn := os.Getenv("DATABASE_URL")
	// if dsn == "" {
	// 	log.Fatal("❌ DATABASE_URL is not set")
	// }

	// log.Println("📦 Connecting to:", dsn)

	dsn := os.Getenv("DATABASE_URL")
		db, err := gorm.Open(postgres.New(postgres.Config{
			DSN: dsn,
			PreferSimpleProtocol: true,
		}), &gorm.Config{})

	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	DB = db
	log.Println("✅ Connected to Supabase PostgreSQL!")
}
