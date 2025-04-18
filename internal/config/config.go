package config

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "os"
)

func SetupDB() (*gorm.DB, error) {
    // Get the database URL from the environment
    dsn := os.Getenv("DATABASE_URL")
    
    // Open the connection to the database using GORM
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    
    return db, nil
}