package message

import (
	"time"

	"github.com/google/uuid"
)

type SendMessageRequest struct {
	RecipientID uuid.UUID `json:"recipient_id"`
	Subject     string    `json:"subject"`
	Body        string    `json:"body"`
}

type MessageResponse struct {
	ID             uuid.UUID `json:"id"`
	SenderID       uuid.UUID `json:"sender_id"`
	SenderName     string    `json:"sender_name"`
	SenderEmail    string    `json:"sender_email"`
	SenderRole     string    `json:"sender_role"`
	RecipientID    uuid.UUID `json:"recipient_id"`
	RecipientName  string    `json:"recipient_name"`
	RecipientEmail string    `json:"recipient_email"`
	RecipientRole  string    `json:"recipient_role"`
	Subject        string    `json:"subject"`
	Body           string    `json:"body"`
	IsRead         bool      `json:"is_read"`
	CreatedAt      time.Time `json:"created_at"`
}

type UserSearchResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
}
