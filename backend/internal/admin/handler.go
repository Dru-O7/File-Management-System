package admin

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetStats returns system-wide statistics
func (h *Handler) GetStats(c echo.Context) error {
	schoolIDStr, _ := c.Get("actor_school_id").(string)
	var schoolID *string
	if schoolIDStr != "" {
		schoolID = &schoolIDStr
	}
	stats, err := h.service.GetStats(schoolID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch stats"})
	}
	return c.JSON(http.StatusOK, stats)
}

// GetUsers returns all users across all schools
func (h *Handler) GetUsers(c echo.Context) error {
	schoolIDStr, _ := c.Get("actor_school_id").(string)
	var schoolID *string
	if schoolIDStr != "" {
		schoolID = &schoolIDStr
	}
	users, err := h.service.GetAllUsers(schoolID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch users"})
	}
	return c.JSON(http.StatusOK, users)
}

// CreateUser creates a new user
func (h *Handler) CreateUser(c echo.Context) error {
	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	user, err := h.service.CreateUser(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, user)
}

// UpdateUser updates an existing user
func (h *Handler) UpdateUser(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}
	var req UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	user, err := h.service.UpdateUser(id, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}

// DeleteUser deletes a user by ID
func (h *Handler) DeleteUser(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}
	if err := h.service.DeleteUser(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete user"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

// GetDocumentTypes returns all document types across all schools
func (h *Handler) GetDocumentTypes(c echo.Context) error {
	schoolIDStr, _ := c.Get("actor_school_id").(string)
	var schoolID *string
	if schoolIDStr != "" {
		schoolID = &schoolIDStr
	}
	dts, err := h.service.GetAllDocumentTypes(schoolID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch document types"})
	}
	return c.JSON(http.StatusOK, dts)
}

// CreateDocumentType creates a new document type
func (h *Handler) CreateDocumentType(c echo.Context) error {
	var req CreateDocTypeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	dt, err := h.service.CreateDocumentType(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, dt)
}

// UpdateDocumentType updates an existing document type
func (h *Handler) UpdateDocumentType(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid document type ID"})
	}
	var req UpdateDocTypeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	dt, err := h.service.UpdateDocumentType(id, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, dt)
}

// DeleteDocumentType deletes a document type
func (h *Handler) DeleteDocumentType(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid document type ID"})
	}
	if err := h.service.DeleteDocumentType(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete document type"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Document type deleted successfully"})
}

// GetSchools returns all schools
func (h *Handler) GetSchools(c echo.Context) error {
	schoolIDStr, _ := c.Get("actor_school_id").(string)
	var schoolID *string
	if schoolIDStr != "" {
		schoolID = &schoolIDStr
	}
	schools, err := h.service.GetAllSchools(schoolID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch schools"})
	}
	return c.JSON(http.StatusOK, schools)
}

// UpdateSchool updates a school's information
func (h *Handler) UpdateSchool(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid school ID"})
	}
	var req UpdateSchoolRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	school, err := h.service.UpdateSchool(id, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, school)
}
