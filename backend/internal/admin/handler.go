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

// GetUsers returns all users across all schools scoped by hierarchy
func (h *Handler) GetUsers(c echo.Context) error {
	actorRole, _ := c.Get("actor_role").(string)
	actorSchoolIDStr, _ := c.Get("actor_school_id").(string)
	var actorSchoolID *uuid.UUID
	if actorSchoolIDStr != "" {
		if u, err := uuid.Parse(actorSchoolIDStr); err == nil {
			actorSchoolID = &u
		}
	}

	users, err := h.service.GetAllUsers(actorRole, actorSchoolID)
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

	actorRole, _ := c.Get("actor_role").(string)
	actorSchoolIDStr, _ := c.Get("actor_school_id").(string)
	var actorSchoolID *uuid.UUID
	if actorSchoolIDStr != "" {
		parsed, err := uuid.Parse(actorSchoolIDStr)
		if err == nil {
			actorSchoolID = &parsed
		}
	}

	user, err := h.service.CreateUser(req, actorRole, actorSchoolID)
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

	actorRole, _ := c.Get("actor_role").(string)
	actorSchoolIDStr, _ := c.Get("actor_school_id").(string)
	var actorSchoolID *uuid.UUID
	if actorSchoolIDStr != "" {
		parsed, err := uuid.Parse(actorSchoolIDStr)
		if err == nil {
			actorSchoolID = &parsed
		}
	}

	user, err := h.service.UpdateUser(id, req, actorRole, actorSchoolID)
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

	actorUserIDStr, _ := c.Get("user_id").(string)
	actorUserID, err := uuid.Parse(actorUserIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized: invalid user ID"})
	}

	if err := h.service.DeleteUser(id, actorUserID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

// GetDocumentTypes returns all document types across all schools
func (h *Handler) GetDocumentTypes(c echo.Context) error {
	actorRole, _ := c.Get("actor_role").(string)
	actorSchoolIDStr, _ := c.Get("actor_school_id").(string)
	var actorSchoolID *uuid.UUID
	if actorSchoolIDStr != "" {
		if u, err := uuid.Parse(actorSchoolIDStr); err == nil {
			actorSchoolID = &u
		}
	}
	dts, err := h.service.GetAllDocumentTypes(actorRole, actorSchoolID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, dts)
}

// CreateDocumentType creates a new document type
func (h *Handler) CreateDocumentType(c echo.Context) error {
	var req CreateDocTypeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	actorRole, _ := c.Get("actor_role").(string)
	actorSchoolIDStr, _ := c.Get("actor_school_id").(string)
	var actorSchoolID *uuid.UUID
	if actorSchoolIDStr != "" {
		if u, err := uuid.Parse(actorSchoolIDStr); err == nil {
			actorSchoolID = &u
		}
	}

	dt, err := h.service.CreateDocumentType(req, actorRole, actorSchoolID)
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

// GetRoles lists all roles
func (h *Handler) GetRoles(c echo.Context) error {
	actorRole, _ := c.Get("actor_role").(string)
	actorSchoolIDStr, _ := c.Get("actor_school_id").(string)
	var actorSchoolID *uuid.UUID
	if actorSchoolIDStr != "" {
		if u, err := uuid.Parse(actorSchoolIDStr); err == nil {
			actorSchoolID = &u
		}
	}

	roles, err := h.service.GetAllRoles(actorRole, actorSchoolID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, roles)
}

// CreateRole creates a new role
func (h *Handler) CreateRole(c echo.Context) error {
	var req CreateRoleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	actorRole, _ := c.Get("actor_role").(string)
	actorSchoolIDStr, _ := c.Get("actor_school_id").(string)
	var actorSchoolID *uuid.UUID
	if actorSchoolIDStr != "" {
		if u, err := uuid.Parse(actorSchoolIDStr); err == nil {
			actorSchoolID = &u
		}
	}

	role, err := h.service.CreateRole(req, actorRole, actorSchoolID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, role)
}

// UpdateRole updates an existing role
func (h *Handler) UpdateRole(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid role ID"})
	}
	var req UpdateRoleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	actorRole, _ := c.Get("actor_role").(string)
	actorSchoolIDStr, _ := c.Get("actor_school_id").(string)
	var actorSchoolID *uuid.UUID
	if actorSchoolIDStr != "" {
		if u, err := uuid.Parse(actorSchoolIDStr); err == nil {
			actorSchoolID = &u
		}
	}

	role, err := h.service.UpdateRole(id, req, actorRole, actorSchoolID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, role)
}

// DeleteRole deletes a role
func (h *Handler) DeleteRole(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid role ID"})
	}

	actorRole, _ := c.Get("actor_role").(string)
	actorSchoolIDStr, _ := c.Get("actor_school_id").(string)
	var actorSchoolID *uuid.UUID
	if actorSchoolIDStr != "" {
		if u, err := uuid.Parse(actorSchoolIDStr); err == nil {
			actorSchoolID = &u
		}
	}

	if err := h.service.DeleteRole(id, actorRole, actorSchoolID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Role deleted successfully"})
}

// ── Organization CRUD (SuperAdmin only) ──────────────────────────────────────

func (h *Handler) GetOrganizations(c echo.Context) error {
	actorRole, _ := c.Get("actor_role").(string)
	actorSchoolIDStr, _ := c.Get("actor_school_id").(string)
	var actorSchoolID *uuid.UUID
	if actorSchoolIDStr != "" {
		if u, err := uuid.Parse(actorSchoolIDStr); err == nil {
			actorSchoolID = &u
		}
	}

	orgs, err := h.service.GetAllOrganizations(actorRole, actorSchoolID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, orgs)
}

func (h *Handler) CreateOrganization(c echo.Context) error {
	var req CreateOrganizationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	actorRole, _ := c.Get("actor_role").(string)
	actorSchoolIDStr, _ := c.Get("actor_school_id").(string)
	var actorSchoolID *uuid.UUID
	if actorSchoolIDStr != "" {
		if u, err := uuid.Parse(actorSchoolIDStr); err == nil {
			actorSchoolID = &u
		}
	}

	org, err := h.service.CreateOrganization(req, actorRole, actorSchoolID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, org)
}

func (h *Handler) UpdateOrganization(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid organization ID"})
	}

	var req UpdateOrganizationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	actorRole, _ := c.Get("actor_role").(string)
	actorSchoolIDStr, _ := c.Get("actor_school_id").(string)
	var actorSchoolID *uuid.UUID
	if actorSchoolIDStr != "" {
		if u, err := uuid.Parse(actorSchoolIDStr); err == nil {
			actorSchoolID = &u
		}
	}

	org, err := h.service.UpdateOrganization(id, req, actorRole, actorSchoolID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, org)
}

func (h *Handler) DeleteOrganization(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid organization ID"})
	}

	actorRole, _ := c.Get("actor_role").(string)
	actorSchoolIDStr, _ := c.Get("actor_school_id").(string)
	var actorSchoolID *uuid.UUID
	if actorSchoolIDStr != "" {
		if u, err := uuid.Parse(actorSchoolIDStr); err == nil {
			actorSchoolID = &u
		}
	}

	if err := h.service.DeleteOrganization(id, actorRole, actorSchoolID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Organization deleted successfully"})
}

// ── Peer Connection Handlers ──────────────────────────────────────────────────

func (h *Handler) GetPeerConnections(c echo.Context) error {
	actorRole, _ := c.Get("actor_role").(string)
	actorSchoolIDStr, _ := c.Get("actor_school_id").(string)
	var actorSchoolID *uuid.UUID
	if actorSchoolIDStr != "" {
		if u, err := uuid.Parse(actorSchoolIDStr); err == nil {
			actorSchoolID = &u
		}
	}

	conns, err := h.service.GetPeerConnections(actorRole, actorSchoolID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, conns)
}

func (h *Handler) RequestPeerConnection(c echo.Context) error {
	var req CreatePeerConnectionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	actorRole, _ := c.Get("actor_role").(string)
	actorSchoolIDStr, _ := c.Get("actor_school_id").(string)
	var actorSchoolID *uuid.UUID
	if actorSchoolIDStr != "" {
		if u, err := uuid.Parse(actorSchoolIDStr); err == nil {
			actorSchoolID = &u
		}
	}

	conn, err := h.service.RequestPeerConnection(req, actorRole, actorSchoolID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, conn)
}

func (h *Handler) AcceptPeerConnection(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid connection ID"})
	}

	actorRole, _ := c.Get("actor_role").(string)
	actorSchoolIDStr, _ := c.Get("actor_school_id").(string)
	var actorSchoolID *uuid.UUID
	if actorSchoolIDStr != "" {
		if u, err := uuid.Parse(actorSchoolIDStr); err == nil {
			actorSchoolID = &u
		}
	}

	conn, err := h.service.AcceptPeerConnection(id, actorRole, actorSchoolID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, conn)
}

func (h *Handler) RejectPeerConnection(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid connection ID"})
	}

	actorRole, _ := c.Get("actor_role").(string)
	actorSchoolIDStr, _ := c.Get("actor_school_id").(string)
	var actorSchoolID *uuid.UUID
	if actorSchoolIDStr != "" {
		if u, err := uuid.Parse(actorSchoolIDStr); err == nil {
			actorSchoolID = &u
		}
	}

	conn, err := h.service.RejectPeerConnection(id, actorRole, actorSchoolID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, conn)
}

func (h *Handler) RevokePeerConnection(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid connection ID"})
	}

	actorRole, _ := c.Get("actor_role").(string)
	actorSchoolIDStr, _ := c.Get("actor_school_id").(string)
	var actorSchoolID *uuid.UUID
	if actorSchoolIDStr != "" {
		if u, err := uuid.Parse(actorSchoolIDStr); err == nil {
			actorSchoolID = &u
		}
	}

	conn, err := h.service.RevokePeerConnection(id, actorRole, actorSchoolID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, conn)
}
