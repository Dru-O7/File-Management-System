package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	authURL, err := url.Parse("http://localhost:8081")
	if err != nil {
		log.Fatal("Invalid Auth Service URL:", err)
	}

	docURL, err := url.Parse("http://localhost:8082")
	if err != nil {
		log.Fatal("Invalid Document Service URL:", err)
	}

	authProxy := httputil.NewSingleHostReverseProxy(authURL)
	docProxy := httputil.NewSingleHostReverseProxy(docURL)

	e.Any("/*", func(c echo.Context) error {
		reqPath := c.Request().URL.Path
		if strings.HasPrefix(reqPath, "/api/auth") || reqPath == "/api/users" {
			authProxy.ServeHTTP(c.Response().Writer, c.Request())
			return nil
		}
		if strings.HasPrefix(reqPath, "/api/documents") {
			docProxy.ServeHTTP(c.Response().Writer, c.Request())
			return nil
		}
		return c.JSON(http.StatusNotFound, map[string]string{"error": "API path not found in Gateway"})
	})

	log.Println("API Gateway starting on port :8080...")
	log.Fatal(e.Start(":8080"))
}
