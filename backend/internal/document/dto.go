package document

import (
	"time"

	"office-file-sharing/backend/internal/shared/models"

	"github.com/google/uuid"
)

type ActionRequest struct {
	ActorID   uuid.UUID  `json:"actor_id"`
	TargetID  *uuid.UUID `json:"target_id"` // Used if sending back or routing elsewhere
	Action    string     `json:"action"`    // Approve, Reject, Sent Back
	Remarks   string     `json:"remarks"`
	Signature string     `json:"signature"`
}

type AttachmentResponse struct {
	ID         uuid.UUID `json:"ID"`
	DocumentID uuid.UUID `json:"DocumentID"`
	Filename   string    `json:"Filename"`
	UploadedBy uuid.UUID `json:"UploadedBy"`
	CreatedAt  time.Time `json:"CreatedAt"`
}

type DocumentResponse struct {
	ID             uuid.UUID             `json:"ID"`
	Filename       string                `json:"Filename"`
	FilePath       string                `json:"FilePath"`
	UploaderID     uuid.UUID             `json:"UploaderID"`
	CurrentOwnerID uuid.UUID             `json:"CurrentOwnerID"`
	Status         models.DocumentStatus `json:"Status"`
	Title          string                `json:"Title"`
	Description    string                `json:"Description"`
	UniqueNumber   string                `json:"UniqueNumber"`
	Tags           string                `json:"Tags"`
	Category       string                `json:"Category"`
	Priority       string                `json:"Priority"`
	Direction      string                `json:"Direction"`
	AssignedAt     time.Time             `json:"AssignedAt"`
	ReferralOwnerID *uuid.UUID            `json:"ReferralOwnerID"`
	NotingSheet    string                `json:"NotingSheet"`
	DraftSpace     string                `json:"DraftSpace"`
	CreatedAt      time.Time             `json:"CreatedAt"`
	UpdatedAt      time.Time             `json:"UpdatedAt"`

	Uploader     models.User          `json:"Uploader"`
	CurrentOwner models.User          `json:"CurrentOwner"`
	Attachments  []AttachmentResponse `json:"Attachments"`
}

type HistoryResponse struct {
	ID         uuid.UUID             `json:"ID"`
	DocumentID uuid.UUID             `json:"DocumentID"`
	ActorID    uuid.UUID             `json:"ActorID"`
	TargetID   *uuid.UUID            `json:"TargetID"`
	Action     models.WorkflowAction `json:"Action"`
	Remarks    string                `json:"Remarks"`
	Signature  string                `json:"Signature"`
	CreatedAt  time.Time             `json:"CreatedAt"`

	Actor  models.User  `json:"Actor"`
	Target *models.User `json:"Target"`
}

type DocumentDetailsResponse struct {
	Document DocumentResponse  `json:"document"`
	History  []HistoryResponse `json:"history"`
}
