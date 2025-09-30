package main

import (
	"go-echo/config"
	"go-echo/internal/auth"
	"go-echo/internal/user"
	"go-echo/pkg/middleware"
	"go-echo/pkg/validator"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	// Load config
	cfg := config.Load()

	// Connect DB
	db := config.ConnectDB(cfg)

	// Init Echo
	e := echo.New()
	e.Validator = validator.NewValidator()

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Public routes
	auth.RegisterRoutes(e, cfg, db)

	// Protected routes
	r := e.Group("/")

	r.Use(echojwt.WithConfig(middleware.JWTMiddleware(cfg.JWTSecret)))
	user.RegisterRoutes(r, db)

	e.Logger.Fatal(e.Start(":8080"))
}
