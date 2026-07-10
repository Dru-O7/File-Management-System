package main

import (
	"log"

	"office-file-sharing/backend/internal/db"
	"office-file-sharing/backend/internal/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db.InitDB()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	api := e.Group("/api")

	// Protected routes
	r := api.Group("")
	r.Use(handlers.AuthMiddleware)

	r.POST("/documents", handlers.UploadDocument)
	r.GET("/documents", handlers.GetDocuments)
	r.GET("/documents/:id", handlers.GetDocumentDetails)
	r.GET("/documents/:id/download", handlers.DownloadDocument)
	r.PUT("/documents/:id/replace", handlers.ReplaceDocument)
	r.POST("/documents/:id/action", handlers.DocumentAction)

	log.Println("Document & Workflow Service starting on port :8082...")
	log.Fatal(e.Start(":8082"))
}
