package user

import (
	"fmt"
	"net/http"

	"office-file-sharing/backend/internal/shared/config"
	"office-file-sharing/backend/internal/shared/email"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetUsers(c echo.Context) error {
	actorIDStr := c.Get("user_id").(string)
	actorID, err := uuid.Parse(actorIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID in token"})
	}

	users, err := h.service.GetUsers(actorID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch users"})
	}
	return c.JSON(http.StatusOK, users)
}

func (h *Handler) SendManualEmail(c echo.Context) error {
	type Request struct {
		To      string `json:"to"`
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}

	var req Request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if req.To == "" || req.Subject == "" || req.Body == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Fields 'to', 'subject', and 'body' are required"})
	}

	cfg := config.Load()
	if cfg.SMTPHost == "" || cfg.SMTPUser == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "SMTP server is not configured in backend .env file"})
	}

	err := email.SendMail(cfg, []string{req.To}, req.Subject, req.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to send email: %v", err)})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Email sent successfully"})
}
