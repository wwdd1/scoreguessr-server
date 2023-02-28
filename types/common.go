package types

import (
	"github.com/golang-jwt/jwt/v4"
)

type JWTClaims struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	ProfilePicture string `json:"profilePicture"`
	jwt.RegisteredClaims
}

type JWTCookieInfo struct {
	Cookie string
	MaxAge int
}

const (
	JWT_EXPIRED = 1
	JWT_INVALID = 2
	JWT_VALID   = 3
)
