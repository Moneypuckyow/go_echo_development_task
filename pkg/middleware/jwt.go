package middleware

import (
	echojwt "github.com/labstack/echo-jwt/v4"
)

func JWTMiddleware(secret string) echojwt.Config {
	return echojwt.Config{
		SigningKey: []byte(secret),
	}
}
