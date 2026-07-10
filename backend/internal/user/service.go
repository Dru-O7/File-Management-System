package user

type Service interface {
	GetUsers() ([]UserResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetUsers() ([]UserResponse, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	responses := make([]UserResponse, len(users))
	for i, u := range users {
		responses[i] = UserResponse{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		}
	}
	return responses, nil
}
