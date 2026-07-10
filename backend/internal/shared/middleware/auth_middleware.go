// Package middleware provides HTTP middleware for authentication, security, and logging
package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type JWTCustomClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

// AuthMiddleware creates a JWT validation middleware with the injected secret key
func AuthMiddleware(jwtSecret []byte) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			var tokenString string

			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) != 2 || parts[0] != "Bearer" {
					return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid authorization header format"})
				}
				tokenString = parts[1]
			} else {
				tokenString = c.QueryParam("token")
				if tokenString == "" {
					return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing authorization token"})
				}
			}
			claims := &JWTCustomClaims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return jwtSecret, nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid or expired authorization token"})
			}

			// Store values in context for other handlers
			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)
			c.Set("name", claims.Name)

			return next(c)
		}
	}
}
