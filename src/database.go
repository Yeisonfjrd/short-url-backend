package main

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func initDatabase() {
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        log.Fatal("DATABASE_URL must be set")
    }

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect database: %v", err)
    }

    err = db.AutoMigrate(&URLMapping{})
    if err != nil {
        log.Fatalf("Failed to migrate database: %v", err)
    }

    DB = db
    log.Println("NeonDB connected and migrated successfully")
}
