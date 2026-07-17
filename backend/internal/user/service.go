package user

import (
	"github.com/google/uuid"
	"office-file-sharing/backend/internal/shared/models"
)

type Service interface {
	GetUsers(actorID uuid.UUID) ([]UserResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetUsers(actorID uuid.UUID) ([]UserResponse, error) {
	// Find actor profile
	var actor models.User
	err := s.repo.(*repository).db.First(&actor, "id = ?", actorID).Error
	if err != nil {
		return nil, err
	}

	allUsers, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var filtered []models.User
	for _, u := range allUsers {
		if u.ID == actorID {
			continue // skip self
		}
		// DHE or SuperAdmin sees everyone
		if actor.Role == "DHE" || actor.Role == "SuperAdmin" || actor.Role == "Admin" {
			filtered = append(filtered, u)
		} else {
			// Everyone else sees people in their own school/office
			if u.SchoolID != nil && actor.SchoolID != nil && *u.SchoolID == *actor.SchoolID {
				filtered = append(filtered, u)
			}
		}
	}

	responses := make([]UserResponse, len(filtered))
	for i, u := range filtered {
		responses[i] = UserResponse{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			Role:      u.Role,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		}
	}
	return responses, nil
}
