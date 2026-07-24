package message

import (
	"office-file-sharing/backend/internal/shared/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateMessage(msg *models.Message) error
	GetMessageByID(id uuid.UUID) (*models.Message, error)
	GetInbox(recipientID uuid.UUID) ([]models.Message, error)
	GetSent(senderID uuid.UUID) ([]models.Message, error)
	MarkAsRead(id uuid.UUID) error
	SearchUsersByQuery(currentUserID uuid.UUID, query string) ([]models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateMessage(msg *models.Message) error {
	return r.db.Create(msg).Error
}

func (r *repository) GetMessageByID(id uuid.UUID) (*models.Message, error) {
	var msg models.Message
	err := r.db.Preload("Sender").Preload("Recipient").First(&msg, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (r *repository) GetInbox(recipientID uuid.UUID) ([]models.Message, error) {
	var msgs []models.Message
	err := r.db.Preload("Sender").Preload("Recipient").
		Where("recipient_id = ?", recipientID).
		Order("created_at desc").
		Find(&msgs).Error
	return msgs, err
}

func (r *repository) GetSent(senderID uuid.UUID) ([]models.Message, error) {
	var msgs []models.Message
	err := r.db.Preload("Sender").Preload("Recipient").
		Where("sender_id = ?", senderID).
		Order("created_at desc").
		Find(&msgs).Error
	return msgs, err
}

func (r *repository) MarkAsRead(id uuid.UUID) error {
	return r.db.Model(&models.Message{}).Where("id = ?", id).Update("is_read", true).Error
}

func (r *repository) SearchUsersByQuery(currentUserID uuid.UUID, query string) ([]models.User, error) {
	var users []models.User
	pattern := "%" + query + "%"
	err := r.db.Where("id != ? AND (LOWER(name) LIKE LOWER(?) OR LOWER(email) LIKE LOWER(?))", currentUserID, pattern, pattern).
		Limit(10).
		Find(&users).Error
	return users, err
}

func (r *repository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("LOWER(email) = LOWER(?)", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
