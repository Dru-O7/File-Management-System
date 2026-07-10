package auth

import "office-file-sharing/backend/internal/shared/models"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}
