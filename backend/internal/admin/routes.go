package admin

import (
	"net/http"
	"office-file-sharing/backend/internal/shared/middleware"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"office-file-sharing/backend/internal/shared/models"
)

// adminAccessMiddleware ensures the user has the "Admin" role.
func adminAccessMiddleware(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userIDStr, ok := c.Get("user_id").(string)
			if !ok || userIDStr == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			}

			var user models.User
			if err := db.First(&user, "id = ?", userIDStr).Error; err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not found"})
			}

			if user.Role != "Admin" && user.Role != "SuperAdmin" && user.Role != "DHE" && user.Role != "School Admin" {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied: Admin role required"})
			}

			// Store role and school for downstream scoping
			c.Set("actor_role", user.Role)
			if user.SchoolID != nil {
				c.Set("actor_school_id", user.SchoolID.String())
			}

			return next(c)
		}
	}
}

// RegisterRoutes registers all admin API routes under /api/admin
func RegisterRoutes(g *echo.Group, handler *Handler, jwtSecret []byte, db *gorm.DB) {
	admin := g.Group("/admin")

	// JWT auth first, then Admin/SuperAdmin role check
	admin.Use(middleware.AuthMiddleware(jwtSecret))
	admin.Use(adminAccessMiddleware(db))

	// Stats
	admin.GET("/stats", handler.GetStats)

	// User management
	admin.GET("/users", handler.GetUsers)
	admin.POST("/users", handler.CreateUser)
	admin.PUT("/users/:id", handler.UpdateUser)
	admin.DELETE("/users/:id", handler.DeleteUser)

	// Document type management
	admin.GET("/document-types", handler.GetDocumentTypes)
	admin.POST("/document-types", handler.CreateDocumentType)
	admin.PUT("/document-types/:id", handler.UpdateDocumentType)
	admin.DELETE("/document-types/:id", handler.DeleteDocumentType)

	// School management (Admin only)
	admin.GET("/schools", handler.GetSchools)
	admin.PUT("/schools/:id", handler.UpdateSchool)
}

