package auth

import (
	"office-file-sharing/backend/internal/shared/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetByEmail(email string) (*models.User, error)
	Create(user *models.User) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetByEmail(email string) (*models.User, error) {
	var u models.User
	if err := r.db.Preload("School").First(&u, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) Create(user *models.User) error {
	return r.db.Create(user).Error
}
