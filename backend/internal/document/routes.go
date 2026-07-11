package document

import (
	"office-file-sharing/backend/internal/shared/middleware"

	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers all document workflows under AuthMiddleware
func RegisterRoutes(g *echo.Group, handler *Handler, jwtSecret []byte) {
	r := g.Group("")
	r.Use(middleware.AuthMiddleware(jwtSecret))

	r.POST("/documents", handler.Upload)
	r.GET("/documents", handler.List)
	r.GET("/document-types", handler.GetDocumentTypes)
	r.GET("/documents/:id", handler.GetDetails)
	r.GET("/documents/:id/download", handler.Download)
	r.GET("/attachments/:id/download", handler.DownloadAttachment)
	r.PUT("/documents/:id/replace", handler.Replace)
	r.POST("/documents/:id/action", handler.TakeAction)
	r.POST("/documents/:id/recall", handler.Recall)
	r.POST("/documents/:id/notes", handler.AppendNote)
	r.PUT("/documents/:id/draft", handler.SaveDraft)
	r.POST("/documents/:id/attachments", handler.AddAttachment)
	r.GET("/notifications", handler.GetNotifications)
	r.GET("/reports", handler.GetReports)
}
