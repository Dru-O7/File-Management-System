package auth

import "github.com/labstack/echo/v4"

// RegisterRoutes sets up auth endpoints under the Echo group
func RegisterRoutes(g *echo.Group, handler *Handler) {
	g.POST("/auth/login", handler.Login)
	g.POST("/auth/signup", handler.Signup)
}
