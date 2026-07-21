package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"office-file-sharing/backend/internal/admin"
	"office-file-sharing/backend/internal/auth"
	"office-file-sharing/backend/internal/document"
	"office-file-sharing/backend/internal/shared/config"
	"office-file-sharing/backend/internal/shared/db"
	"office-file-sharing/backend/internal/shared/models"
	"office-file-sharing/backend/internal/user"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/time/rate"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()
	database := db.Init(cfg.DatabaseURL)
	seedData(database)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Rate limiter specifically for authentication endpoints to prevent brute force
	authRateLimiter := middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Skipper: func(c echo.Context) bool {
			return !strings.HasPrefix(c.Request().URL.Path, "/api/auth/")
		},
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{
				Rate:      rate.Limit(5.0 / 60.0), // 5 requests per minute
				Burst:     5,
				ExpiresIn: 3 * time.Minute,
			},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			return ctx.RealIP(), nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusTooManyRequests, map[string]string{"error": "Too many requests. Please try again later."})
		},
	})

	e.Use(authRateLimiter)

	api := e.Group("/api")

	// Repositories
	authRepo := auth.NewRepository(database)
	userRepo := user.NewRepository(database)
	docRepo := document.NewRepository(database)
	adminRepo := admin.NewRepository(database)

	// Services
	authService := auth.NewService(authRepo, []byte(cfg.JWTSecret))
	userService := user.NewService(userRepo)
	docService := document.NewService(docRepo, "./uploads")
	adminService := admin.NewService(adminRepo)

	// Handlers
	authHandler := auth.NewHandler(authService)
	userHandler := user.NewHandler(userService, database)
	docHandler := document.NewHandler(docService)
	adminHandler := admin.NewHandler(adminService)

	// Register Modular Routes
	auth.RegisterRoutes(api, authHandler)
	user.RegisterRoutes(api, userHandler, []byte(cfg.JWTSecret))
	document.RegisterRoutes(api, docHandler, []byte(cfg.JWTSecret))
	admin.RegisterRoutes(api, adminHandler, []byte(cfg.JWTSecret), database)



	log.Println("Modular Academic Monolith starting on port :8080...")
	log.Fatal(e.Start(":8080"))
}



func seedData(gormDB *gorm.DB) {
	// Ensure SuperAdmin role exists
	var superAdminRole models.Role
	err := gormDB.Where("role_name = ? AND tenant_id IS NULL", "SuperAdmin").First(&superAdminRole).Error
	if err != nil {
		newRoleID := uuid.New()
		superAdminRole = models.Role{
			ID:            newRoleID,
			RoleName:      "SuperAdmin",
			IsAdminAccess: true,
			ParentRoleID:  nil,
			TenantID:      nil,
			CreatedBy:     "System",
			Path:          "/" + newRoleID.String() + "/",
		}
		if err := gormDB.Create(&superAdminRole).Error; err != nil {
			log.Fatalf("Failed to seed SuperAdmin role: %v", err)
		}
		log.Println("Seeded SuperAdmin role.")
	}

	// Ensure default SuperAdmin user exists
	var superAdminUser models.User
	err = gormDB.Where("role = ?", "SuperAdmin").First(&superAdminUser).Error
	if err != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal(err)
		}
		superAdminUser = models.User{
			ID:           uuid.New(),
			Name:         "Super Admin",
			Email:        "superadmin@school.edu",
			PasswordHash: string(hash),
			Role:         "SuperAdmin",
			SchoolID:     nil,
		}
		if err := gormDB.Create(&superAdminUser).Error; err != nil {
			log.Fatalf("Failed to seed SuperAdmin user: %v", err)
		}
		log.Println("Seeded default SuperAdmin user (superadmin@school.edu).")
	}
}
