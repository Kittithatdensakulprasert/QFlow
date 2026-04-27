package db

import (
	"log"
	"qflow/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	DB = db
	return db
}

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&domain.User{},
		&domain.OTP{},
		&domain.Category{},
		&domain.Provider{},
		&domain.Zone{},
		&domain.Queue{},
		&domain.Notification{},
	)
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}
	log.Println("database migration completed")
}
