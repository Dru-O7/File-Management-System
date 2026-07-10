package user

import "github.com/labstack/echo/v4"

// RegisterRoutes registers user-related endpoints under the Echo group
func RegisterRoutes(g *echo.Group, handler *Handler) {
	g.GET("/users", handler.GetUsers)
}
