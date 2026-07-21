package main

import (
	"fmt"
	"log"

	"office-file-sharing/backend/internal/shared/config"
	"office-file-sharing/backend/internal/shared/db"
	"office-file-sharing/backend/internal/shared/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg := config.Load()
	gormDB := db.Init(cfg.DatabaseURL)

	log.Println("Seeding database - resetting all tables except SuperAdmin...")

	// 1. Truncate document-related and organization tables
	tables := []string{
		"workflow_histories",
		"notifications",
		"attachments",
		"files",
		"documents",
		"peer_connections",
		"organizations",
		"schools",
	}

	for _, table := range tables {
		err := gormDB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE;", table)).Error
		if err != nil {
			log.Printf("Warning: error truncating %s: %v\n", table, err)
		} else {
			log.Printf("Truncated %s\n", table)
		}
	}

	// 2. Delete all roles except SuperAdmin
	err := gormDB.Exec("DELETE FROM roles WHERE role_name != 'SuperAdmin';").Error
	if err != nil {
		log.Printf("Warning: error clearing roles: %v\n", err)
	}

	// 3. Delete all users except SuperAdmin
	err = gormDB.Exec("DELETE FROM users WHERE role != 'SuperAdmin';").Error
	if err != nil {
		log.Printf("Warning: error clearing users: %v\n", err)
	}

	// 4. Ensure SuperAdmin role exists
	var superAdminRole models.Role
	err = gormDB.Where("role_name = ? AND tenant_id IS NULL", "SuperAdmin").First(&superAdminRole).Error
	if err != nil {
		newRoleID := uuid.New()
		superAdminRole = models.Role{
			ID:            newRoleID,
			RoleName:      "SuperAdmin",
			IsAdminAccess: true,
			ParentRoleID:  nil,
			TenantID:      nil,
			CreatedBy:     "System",
			Path:          "/" + newRoleID.String() + "/",
		}
		if err := gormDB.Create(&superAdminRole).Error; err != nil {
			log.Fatalf("Failed to seed SuperAdmin role: %v", err)
		}
		log.Println("Seeded SuperAdmin role.")
	} else {
		log.Println("SuperAdmin role already exists.")
	}

	// 5. Ensure default SuperAdmin user exists
	var superAdminUser models.User
	err = gormDB.Where("role = ?", "SuperAdmin").First(&superAdminUser).Error
	if err != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal(err)
		}
		superAdminUser = models.User{
			ID:           uuid.New(),
			Name:         "Super Admin",
			Email:        "superadmin@school.edu",
			PasswordHash: string(hash),
			Role:         "SuperAdmin",
			SchoolID:     nil,
		}
		if err := gormDB.Create(&superAdminUser).Error; err != nil {
			log.Fatalf("Failed to seed SuperAdmin user: %v", err)
		}
		log.Println("Seeded default SuperAdmin user (superadmin@school.edu).")
	} else {
		log.Println("SuperAdmin user already exists.")
	}

	log.Println("Database seeding completed successfully.")
}
