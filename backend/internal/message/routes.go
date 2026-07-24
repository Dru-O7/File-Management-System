package message

import (
	"office-file-sharing/backend/internal/shared/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, handler *Handler, jwtSecret []byte) {
	messages := g.Group("/messages", middleware.AuthMiddleware(jwtSecret))

	messages.POST("", handler.SendMessage)
	messages.GET("/inbox", handler.GetInbox)
	messages.GET("/sent", handler.GetSent)
	messages.GET("/search-users", handler.SearchUsers)
	messages.GET("/by-email", handler.GetUserByEmail)
	messages.GET("/:id", handler.GetMessageDetails)
}
