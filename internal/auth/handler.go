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

// Login godoc
// @Summary      User login
// @Description  Authenticate user and return token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials body LoginRequest true "Login credentials"
// @Success      200 {object} LoginResponse
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Router       /login [post]
func (h *Handler) Login(c echo.Context) error {
	var req LoginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid input"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	// Validasi username dan password
	if req.Username != h.loginUser.Username || req.Password != h.loginUser.Password {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Invalid username or password",
		})
	}

	// Generate token
	token, err := h.service.GenerateToken(req.Username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to generate token",
		})
	}

	return c.JSON(http.StatusOK, LoginResponse{
		Message: "Login successful",
		Token:   token,
		User:    req.Username,
	})
}

func RegisterRoutes(e *echo.Echo, cfg *config.Config, db *sql.DB) {
	svc := NewService(cfg.JWTSecret)
	h := NewHandler(svc, db)

	e.POST("/login", h.Login)
}
