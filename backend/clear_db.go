package main

import (
	"fmt"
	"log"

	"office-file-sharing/backend/internal/shared/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=110085 dbname=office_files port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Clearing database...")

	// 1. Truncate document-related and organization tables
	tables := []string{
		"workflow_histories",
		"notifications",
		"attachments",
		"files",
		"documents",
		"peer_connections",
		"organizations",
	}

	for _, table := range tables {
		err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE;", table)).Error
		if err != nil {
			fmt.Printf("Error truncating %s: %v\n", table, err)
		} else {
			fmt.Printf("Truncated %s\n", table)
		}
	}

	// 2. Delete non-SuperAdmin roles (tenant-specific roles)
	_ = db.Exec("UPDATE roles SET parent_role_id = NULL;").Error
	err = db.Exec("DELETE FROM roles WHERE tenant_id IS NOT NULL;").Error
	if err != nil {
		log.Printf("Error clearing custom roles: %v\n", err)
	} else {
		log.Println("Cleared custom tenant roles.")
	}

	// 3. Delete all users except SuperAdmin
	err = db.Exec("DELETE FROM users WHERE role != 'SuperAdmin';").Error
	if err != nil {
		log.Printf("Error clearing users: %v\n", err)
	} else {
		log.Println("Cleared all non-SuperAdmin users.")
	}

	// 4. Delete all schools (tenants)
	err = db.Exec("TRUNCATE TABLE schools CASCADE;").Error
	if err != nil {
		log.Printf("Error clearing schools: %v\n", err)
	} else {
		log.Println("Cleared all schools/tenants.")
	}

	// 5. Ensure SuperAdmin user exists
	var sa models.User
	err = db.Where("role = ?", "SuperAdmin").First(&sa).Error
	if err != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal(err)
		}
		newSA := models.User{
			ID:           uuid.New(),
			Name:         "Super Admin",
			Email:        "superadmin@school.edu",
			PasswordHash: string(hash),
			Role:         "SuperAdmin",
			SchoolID:     nil,
		}
		if err := db.Create(&newSA).Error; err != nil {
			log.Printf("Error seeding SuperAdmin: %v\n", err)
		} else {
			log.Println("Seeded default SuperAdmin user (superadmin@school.edu).")
		}
	} else {
		log.Printf("SuperAdmin user already exists: %s\n", sa.Email)
	}

	log.Println("Database cleanup completed successfully.")
}
