// Package db handles database connections and shared pool configuration
package db

import (
	"log"

	"office-file-sharing/backend/internal/shared/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Init initializes connection to PostgreSQL and runs AutoMigrate on shared schemas
func Init(databaseURL string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")

	err = db.AutoMigrate(
		&models.School{},
		&models.User{},
		&models.ParentChild{},
		&models.DocumentType{},
		&models.Document{},
		&models.DocumentPendingApprover{},
		&models.WorkflowHistory{},
		&models.Notification{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migration completed")

	return db
}
