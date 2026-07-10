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

		switch actor.Role {
		case "Principal":
			// Principal sees everyone in school
			if u.SchoolID != nil && actor.SchoolID != nil && *u.SchoolID == *actor.SchoolID {
				filtered = append(filtered, u)
			}
		case "Teacher":
			// Teacher sees principal, parents, and students of their class section
			if u.Role == "Principal" || u.Role == "Teacher" || u.Role == "Parent" || (u.Role == "Student" && u.ClassSection == actor.ClassSection) {
				filtered = append(filtered, u)
			}
		case "Student":
			// Student only sees teachers and Principal (approvers)
			if u.Role == "Teacher" || u.Role == "Principal" {
				filtered = append(filtered, u)
			}
		case "Parent":
			// Parent sees teachers, Principal, and their own children
			isChild := false
			if u.Role == "Student" {
				var count int64
				s.repo.(*repository).db.Model(&models.ParentChild{}).
					Where("parent_id = ? AND child_id = ?", actorID, u.ID).
					Count(&count)
				isChild = count > 0
			}
			if u.Role == "Teacher" || u.Role == "Principal" || isChild {
				filtered = append(filtered, u)
			}
		default:
			filtered = append(filtered, u)
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
