package models

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	UserID   int
	Username string
	Role     string
	jwt.RegisteredClaims
}
