package admin

import (
	"time"

	"github.com/google/uuid"
)

// ── System Stats ──────────────────────────────────────────────────────────────

type SystemStats struct {
	TotalUsers        int64 `json:"total_users"`
	TotalDocuments    int64 `json:"total_documents"`
	TotalSchools      int64 `json:"total_schools"`
	TotalDocumentTypes int64 `json:"total_document_types"`
	ActiveDocuments   int64 `json:"active_documents"`
	ApprovedDocuments int64 `json:"approved_documents"`
	SLABreaches       int64 `json:"sla_breaches"`
	PendingDocuments  int64 `json:"pending_documents"`
}

// ── User Management ───────────────────────────────────────────────────────────

type UserResponse struct {
	ID           uuid.UUID  `json:"ID"`
	Name         string     `json:"Name"`
	Email        string     `json:"Email"`
	Role         string     `json:"Role"`
	SchoolID     *uuid.UUID `json:"SchoolID"`
	SchoolName   string     `json:"SchoolName"`
	ClassSection string     `json:"ClassSection"`
	Subject      string     `json:"Subject"`
	Phone        string     `json:"Phone"`
	CreatedAt    time.Time  `json:"CreatedAt"`
}

type CreateUserRequest struct {
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	Password     string     `json:"password"`
	Role         string     `json:"role"`
	SchoolID     *uuid.UUID `json:"school_id"`
	ClassSection string     `json:"class_section"`
	Subject      string     `json:"subject"`
	Phone        string     `json:"phone"`
}

type UpdateUserRequest struct {
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	Role         string     `json:"role"`
	SchoolID     *uuid.UUID `json:"school_id"`
	ClassSection string     `json:"class_section"`
	Subject      string     `json:"subject"`
	Phone        string     `json:"phone"`
	Password     string     `json:"password"` // optional — only set if non-empty
}

// ── Document Type Management ──────────────────────────────────────────────────

type DocumentTypeResponse struct {
	ID                uuid.UUID `json:"ID"`
	SchoolID          uuid.UUID `json:"SchoolID"`
	SchoolName        string    `json:"SchoolName"`
	Name              string    `json:"Name"`
	Slug              string    `json:"Slug"`
	WorkflowStages    string    `json:"WorkflowStages"`
	RequiredFields    string    `json:"RequiredFields"`
	SlaHours          int       `json:"SlaHours"`
	Active            bool      `json:"Active"`
}

type CreateDocTypeRequest struct {
	SchoolID          uuid.UUID `json:"school_id"`
	Name              string    `json:"name"`
	Slug              string    `json:"slug"`
	WorkflowStages    string    `json:"workflow_stages"`
	RequiredFields    string    `json:"required_fields"`
	SlaHours          int       `json:"sla_hours"`
}

type UpdateDocTypeRequest struct {
	Name              string `json:"name"`
	Slug              string `json:"slug"`
	WorkflowStages    string `json:"workflow_stages"`
	RequiredFields    string `json:"required_fields"`
	SlaHours          int    `json:"sla_hours"`
	Active            bool   `json:"active"`
}

// ── School Management ─────────────────────────────────────────────────────────

type SchoolResponse struct {
	ID        uuid.UUID `json:"ID"`
	Name      string    `json:"Name"`
	Slug      string    `json:"Slug"`
	Settings  string    `json:"Settings"`
	CreatedAt time.Time `json:"CreatedAt"`
}

type UpdateSchoolRequest struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Settings string `json:"settings"`
}
