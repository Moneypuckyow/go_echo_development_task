package main

import (
	"go-echo/config"
	_ "go-echo/docs"
	"go-echo/internal/auth"
	"go-echo/internal/user"
	"go-echo/pkg/middleware"
	"go-echo/pkg/validator"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title           Go Echo API
// @version         1.0
// @description     RESTful API for user management with JWT authentication

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:9090
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your bearer token in the format **Bearer &lt;token&gt;**
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

	e.Logger.Fatal(e.Start(":9090"))
}
