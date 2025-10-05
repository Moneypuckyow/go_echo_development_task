package auth

import "github.com/golang-jwt/jwt/v5"

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents successful login response
type LoginResponse struct {
	Message string `json:"message" example:"Login successful"`
	Token   string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User    string `json:"user" example:"admin"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid username or password"`
}

type JwtCustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
