package message

import (
	"errors"
	"office-file-sharing/backend/internal/shared/models"
	"strings"

	"github.com/google/uuid"
)

type Service interface {
	SendMessage(senderID uuid.UUID, req SendMessageRequest) (*MessageResponse, error)
	GetInbox(recipientID uuid.UUID) ([]MessageResponse, error)
	GetSent(senderID uuid.UUID) ([]MessageResponse, error)
	GetMessageDetails(id, currentUserID uuid.UUID) (*MessageResponse, error)
	SearchUsers(currentUserID uuid.UUID, query string) ([]UserSearchResponse, error)
	GetUserByEmail(email string) (*UserSearchResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) SendMessage(senderID uuid.UUID, req SendMessageRequest) (*MessageResponse, error) {
	if req.RecipientID == uuid.Nil {
		return nil, errors.New("recipient is required")
	}
	if strings.TrimSpace(req.Subject) == "" {
		return nil, errors.New("subject is required")
	}
	if strings.TrimSpace(req.Body) == "" {
		return nil, errors.New("message body is required")
	}

	msg := &models.Message{
		ID:          uuid.New(),
		SenderID:    senderID,
		RecipientID: req.RecipientID,
		Subject:     strings.TrimSpace(req.Subject),
		Body:        strings.TrimSpace(req.Body),
		IsRead:      false,
	}

	if err := s.repo.CreateMessage(msg); err != nil {
		return nil, err
	}

	fetchedMsg, err := s.repo.GetMessageByID(msg.ID)
	if err != nil {
		return nil, err
	}

	return s.toMessageResponse(fetchedMsg), nil
}

func (s *service) GetInbox(recipientID uuid.UUID) ([]MessageResponse, error) {
	msgs, err := s.repo.GetInbox(recipientID)
	if err != nil {
		return nil, err
	}

	responses := make([]MessageResponse, len(msgs))
	for i, m := range msgs {
		responses[i] = *s.toMessageResponse(&m)
	}
	return responses, nil
}

func (s *service) GetSent(senderID uuid.UUID) ([]MessageResponse, error) {
	msgs, err := s.repo.GetSent(senderID)
	if err != nil {
		return nil, err
	}

	responses := make([]MessageResponse, len(msgs))
	for i, m := range msgs {
		responses[i] = *s.toMessageResponse(&m)
	}
	return responses, nil
}

func (s *service) GetMessageDetails(id, currentUserID uuid.UUID) (*MessageResponse, error) {
	msg, err := s.repo.GetMessageByID(id)
	if err != nil {
		return nil, errors.New("message not found")
	}

	if msg.SenderID != currentUserID && msg.RecipientID != currentUserID {
		return nil, errors.New("unauthorized to view this message")
	}

	// Mark as read if current user is the recipient
	if msg.RecipientID == currentUserID && !msg.IsRead {
		_ = s.repo.MarkAsRead(id)
		msg.IsRead = true
	}

	return s.toMessageResponse(msg), nil
}

func (s *service) SearchUsers(currentUserID uuid.UUID, query string) ([]UserSearchResponse, error) {
	query = strings.TrimSpace(query)
	if len(query) < 1 {
		return []UserSearchResponse{}, nil
	}

	users, err := s.repo.SearchUsersByQuery(currentUserID, query)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 && strings.Contains(query, "@") {
		exactUser, err := s.repo.GetUserByEmail(query)
		if err == nil && exactUser != nil && exactUser.ID != currentUserID {
			users = append(users, *exactUser)
		}
	}

	res := make([]UserSearchResponse, len(users))
	for i, u := range users {
		res[i] = UserSearchResponse{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
			Role:  u.Role,
		}
	}
	return res, nil
}

func (s *service) GetUserByEmail(email string) (*UserSearchResponse, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return nil, errors.New("email is required")
	}

	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("user not found with this email")
	}

	return &UserSearchResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (s *service) toMessageResponse(m *models.Message) *MessageResponse {
	return &MessageResponse{
		ID:             m.ID,
		SenderID:       m.SenderID,
		SenderName:     m.Sender.Name,
		SenderEmail:    m.Sender.Email,
		SenderRole:     m.Sender.Role,
		RecipientID:    m.RecipientID,
		RecipientName:  m.Recipient.Name,
		RecipientEmail: m.Recipient.Email,
		RecipientRole:  m.Recipient.Role,
		Subject:        m.Subject,
		Body:           m.Body,
		IsRead:         m.IsRead,
		CreatedAt:      m.CreatedAt,
	}
}
