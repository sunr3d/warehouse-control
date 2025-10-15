package authsvc

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/sunr3d/warehouse-control/internal/interfaces/infra"
	"github.com/sunr3d/warehouse-control/internal/interfaces/services"
	"github.com/sunr3d/warehouse-control/models"
)

const (
	tokenTTL = 12 * time.Hour
)

var _ services.AuthService = (*authSvc)(nil)

type authSvc struct {
	db        infra.Database
	jwtSecret string
}

// New - конструктор сервиса авторизации.
func New(db infra.Database, jwtSecret string) services.AuthService {
	return &authSvc{db: db, jwtSecret: jwtSecret}
}

// Login - метод для авторизации пользователя, возвращает токен.
func (s *authSvc) Login(ctx context.Context, username, pass string) (string, error) {
	user, err := s.db.GetByUsername(ctx, username)
	if err != nil {
		return "", fmt.Errorf("db.GetByUsername: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(pass))
	if err != nil {
		return "", fmt.Errorf("bcrypt.CompareHashAndPassword: %w", err)
	}

	claims := &models.JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("token.SignedString: %w", err)
	}

	return tokenStr, nil
}

// ValidateToken - метод для валидации токена.
func (s *authSvc) ValidateToken(tokenStr string) (*models.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &models.JWTClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("jwt.ParseWithClaims: %w", err)
	}

	if claims, ok := token.Claims.(*models.JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("невалидный токен")
}
