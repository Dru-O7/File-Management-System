package main

import (
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

	log.Println("Resetting and seeding database accounts...")

	// Clear existing users and cascade child references
	tx := gormDB.Exec("TRUNCATE TABLE users CASCADE;")
	if tx.Error != nil {
		log.Fatalf("Failed to truncate users table: %v", tx.Error)
	}

	// 1. Resolve or seed Greenwood High School
	var school models.School
	err := gormDB.First(&school, "slug = ?", "greenwood-high").Error
	if err != nil {
		school = models.School{
			ID:   uuid.New(),
			Name: "Greenwood High School",
			Slug: "greenwood-high",
		}
		gormDB.Create(&school)
		log.Println("Seeded school: Greenwood High School")
	}

	// 2. Hash default password
	hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash default password:", err)
	}

	// 3. Define users
	users := []models.User{
		{Name: "Aarav Sharma", Email: "aarav@school.edu", PasswordHash: string(hash), Role: "vocational", SchoolID: &school.ID, ClassSection: "Department A"},
		{Name: "Priya Patel", Email: "priya@school.edu", PasswordHash: string(hash), Role: "Teaching staff", SchoolID: &school.ID, ClassSection: "Department A", Subject: "Science"},
		{Name: "Rahul Gupta", Email: "rahul@school.edu", PasswordHash: string(hash), Role: "School Admin", SchoolID: &school.ID},
		{Name: "Deepak Singh", Email: "deepak@school.edu", PasswordHash: string(hash), Role: "non-teaching", SchoolID: &school.ID},
		{Name: "Neha Reddy", Email: "neha@school.edu", PasswordHash: string(hash), Role: "Teaching staff", SchoolID: &school.ID, ClassSection: "Department B", Subject: "History"},
		{Name: "Vikram Iyer", Email: "vikram@school.edu", PasswordHash: string(hash), Role: "Teaching staff", SchoolID: &school.ID, ClassSection: "Department C", Subject: "Mathematics"},
		{Name: "Meera Menon", Email: "meera@school.edu", PasswordHash: string(hash), Role: "Teaching staff", SchoolID: &school.ID, ClassSection: "Department D", Subject: "English"},
		{Name: "Gaurav Verma", Email: "gaurav@school.edu", PasswordHash: string(hash), Role: "School Admin", SchoolID: &school.ID},
		{Name: "System Administrator", Email: "admin@school.edu", PasswordHash: string(hash), Role: "DHE", SchoolID: &school.ID},
	}

	for i := range users {
		users[i].ID = uuid.New()
		gormDB.Create(&users[i])
	}
	log.Println("Seeded school-scoped users.")

	// 4. Ensure document types are seeded
	var docTypeCount int64
	gormDB.Model(&models.DocumentType{}).Count(&docTypeCount)
	if docTypeCount == 0 {
		docTypes := []models.DocumentType{
			{
				SchoolID:       school.ID,
				Name:           "Staff Grievance",
				Slug:           "staff-grievance",
				WorkflowStages: `[{"stage": 1, "role": "Teaching staff", "label": "Department Head", "optional": false}]`,
				RequiredFields: `[]`,
				SlaHours:       72,
			},
			{
				SchoolID:       school.ID,
				Name:           "Infrastructure Issue",
				Slug:           "infrastructure-issue",
				WorkflowStages: `[{"stage": 1, "role": "School Admin", "label": "School Admin Final approval", "optional": false}]`,
				RequiredFields: `["reason", "urgency"]`,
				SlaHours:       120,
			},
			{
				SchoolID:       school.ID,
				Name:           "Disciplinary Issue",
				Slug:           "disciplinary-issue",
				WorkflowStages: `[{"stage": 1, "role": "Teaching staff", "label": "Department Head", "optional": false}]`,
				RequiredFields: `["event_name", "event_date"]`,
				SlaHours:       24,
			},
			{
				SchoolID:       school.ID,
				Name:           "Audit Report",
				Slug:           "audit-report",
				WorkflowStages: `[{"stage": 1, "role": "School Admin", "label": "School Admin Approval", "optional": false}]`,
				RequiredFields: `["audit_reason", "percentage"]`,
				SlaHours:       96,
			},
			{
				SchoolID:       school.ID,
				Name:           "Official Circular",
				Slug:           "official-circular",
				WorkflowStages: `[]`,
				RequiredFields: `[]`,
				SlaHours:       0,
			},
		}
		for i := range docTypes {
			docTypes[i].ID = uuid.New()
			gormDB.Create(&docTypes[i])
		}
		log.Println("Database seeded with document types.")
	}

	log.Println("Database seeding completed successfully.")
}
