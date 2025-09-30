package auth

import (
	"database/sql"
	"go-echo/config"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service   *Service
	db        *sql.DB
	loginUser *LoginUser
}

type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoadLoginCredentials membaca variabel environment untuk login
func LoadLoginCredentials() *LoginUser {
	if err := godotenv.Load("config/.env"); err != nil {
		log.Println("⚠️ .env file not found, using system environment variables")
	}

	return &LoginUser{
		Username: os.Getenv("LOGIN_USERNAME"),
		Password: os.Getenv("LOGIN_PASSWORD"),
	}
}

func NewHandler(service *Service, db *sql.DB) *Handler {
	return &Handler{
		service:   service,
		db:        db,
		loginUser: LoadLoginCredentials(),
	}
}

func (h *Handler) Login(c echo.Context) error {
	var req LoginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// Validasi username dan password dari environment variables
	if req.Username != h.loginUser.Username || req.Password != h.loginUser.Password {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "Invalid username or password",
		})
	}

	// Generate token jika login berhasil
	token, err := h.service.GenerateToken(req.Username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to generate token",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Login successful",
		"token":   token,
		"user":    req.Username,
	})
}

func RegisterRoutes(e *echo.Echo, cfg *config.Config, db *sql.DB) {
	svc := NewService(cfg.JWTSecret)
	h := NewHandler(svc, db)

	e.POST("/login", h.Login)
}
