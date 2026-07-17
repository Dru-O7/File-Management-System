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
		{Name: "Alice Smith", Email: "alice@school.edu", PasswordHash: string(hash), Role: "vocational", SchoolID: &school.ID, ClassSection: "Department A"},
		{Name: "Bob Johnson", Email: "bob@school.edu", PasswordHash: string(hash), Role: "Teaching staff", SchoolID: &school.ID, ClassSection: "Department A", Subject: "Science"},
		{Name: "Charlie Brown", Email: "charlie@school.edu", PasswordHash: string(hash), Role: "School Admin", SchoolID: &school.ID},
		{Name: "David Smith", Email: "david@school.edu", PasswordHash: string(hash), Role: "non-teaching", SchoolID: &school.ID},
		{Name: "Diana Prince", Email: "diana@school.edu", PasswordHash: string(hash), Role: "Teaching staff", SchoolID: &school.ID, ClassSection: "Department B", Subject: "History"},
		{Name: "Evan Wright", Email: "evan@school.edu", PasswordHash: string(hash), Role: "Teaching staff", SchoolID: &school.ID, ClassSection: "Department C", Subject: "Mathematics"},
		{Name: "Fiona Gallagher", Email: "fiona@school.edu", PasswordHash: string(hash), Role: "Teaching staff", SchoolID: &school.ID, ClassSection: "Department D", Subject: "English"},
		{Name: "George Vance", Email: "george@school.edu", PasswordHash: string(hash), Role: "School Admin", SchoolID: &school.ID},
		{Name: "System Administrator", Email: "admin@school.edu", PasswordHash: string(hash), Role: "DHE", SchoolID: &school.ID},
	}

	for i := range users {
		users[i].ID = uuid.New()
		gormDB.Create(&users[i])
	}
	log.Println("Seeded school-scoped users.")

	// Establish Parent-Child link (David is Alice's parent)
	var alice, david models.User
	gormDB.First(&alice, "email = ?", "alice@school.edu")
	gormDB.First(&david, "email = ?", "david@school.edu")
	if alice.ID != uuid.Nil && david.ID != uuid.Nil {
		pc := models.ParentChild{
			ParentID: david.ID,
			ChildID:  alice.ID,
		}
		gormDB.Create(&pc)
		log.Println("Established Parent-Child relationship: David -> Alice")
	}

	// 4. Ensure document types are seeded
	var docTypeCount int64
	gormDB.Model(&models.DocumentType{}).Count(&docTypeCount)
	if docTypeCount == 0 {
		docTypes := []models.DocumentType{
			{
				SchoolID:          school.ID,
				Name:              "Staff Grievance",
				Slug:              "staff-grievance",
				WorkflowStages:    `[{"stage": 1, "role": "Teaching staff", "label": "Department Head", "optional": false}]`,
				RequiredFields:    `[]`,
				SlaHours:          72,
				NeedsParentCosign: false,
			},
			{
				SchoolID:          school.ID,
				Name:              "Infrastructure Issue",
				Slug:              "infrastructure-issue",
				WorkflowStages:    `[{"stage": 1, "role": "School Admin", "label": "School Admin Final approval", "optional": false}]`,
				RequiredFields:    `["reason", "urgency"]`,
				SlaHours:          120,
				NeedsParentCosign: false,
			},
			{
				SchoolID:          school.ID,
				Name:              "Disciplinary Issue",
				Slug:              "disciplinary-issue",
				WorkflowStages:    `[{"stage": 1, "role": "Teaching staff", "label": "Department Head", "optional": false}]`,
				RequiredFields:    `["event_name", "event_date"]`,
				SlaHours:          24,
				NeedsParentCosign: false,
			},
			{
				SchoolID:          school.ID,
				Name:              "Audit Report",
				Slug:              "audit-report",
				WorkflowStages:    `[{"stage": 1, "role": "School Admin", "label": "School Admin Approval", "optional": false}]`,
				RequiredFields:    `["audit_reason", "percentage"]`,
				SlaHours:          96,
				NeedsParentCosign: false,
			},
			{
				SchoolID:          school.ID,
				Name:              "Official Circular",
				Slug:              "official-circular",
				WorkflowStages:    `[]`,
				RequiredFields:    `[]`,
				SlaHours:          0,
				NeedsParentCosign: false,
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
