package auth

import (
	"office-file-sharing/backend/internal/shared/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	CheckAdminAccess(roleName string) bool
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

func (r *repository) CheckAdminAccess(roleName string) bool {
	if roleName == "SuperAdmin" || roleName == "Admin" || roleName == "DHE" || roleName == "School Admin" {
		return true
	}

	var role models.Role
	if err := r.db.First(&role, "role_name = ?", roleName).Error; err != nil {
		return false
	}

	curr := &role
	for curr != nil {
		if curr.IsAdminAccess {
			return true
		}
		if curr.ParentRoleID == nil {
			break
		}
		var parent models.Role
		if err := r.db.First(&parent, "id = ?", *curr.ParentRoleID).Error; err != nil {
			break
		}
		curr = &parent
	}
	return false
}
