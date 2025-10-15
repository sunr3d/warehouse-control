package authsvc

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/sunr3d/warehouse-control/mocks"
	"github.com/sunr3d/warehouse-control/models"
)

// TestAuthSvc_Login - тесты для метода Login
func TestAuthSvc_Login_OK(t *testing.T) {
	mockDB := mocks.NewDatabase(t)
	svc := New(mockDB, "test-secret")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	user := &models.User{
		ID:           1,
		Username:     "admin123",
		PasswordHash: string(hashedPassword),
		Role:         "admin",
	}

	mockDB.EXPECT().
		GetByUsername(mock.Anything, "admin123").
		Return(user, nil)

	token, err := svc.Login(context.Background(), "admin123", "password")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := svc.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, 1, claims.UserID)
	assert.Equal(t, "admin123", claims.Username)
	assert.Equal(t, "admin", claims.Role)
}

func TestAuthSvc_Login_ErrInvalidPassword(t *testing.T) {
	mockDB := mocks.NewDatabase(t)
	svc := New(mockDB, "test-secret")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	user := &models.User{
		ID:           1,
		Username:     "admin123",
		PasswordHash: string(hashedPassword),
		Role:         "admin",
	}

	mockDB.EXPECT().
		GetByUsername(mock.Anything, "admin123").
		Return(user, nil)

	token, err := svc.Login(context.Background(), "admin123", "wrong-password")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "bcrypt.CompareHashAndPassword")
}

func TestAuthSvc_Login_ErrUserNotFound(t *testing.T) {
	mockDB := mocks.NewDatabase(t)
	svc := New(mockDB, "test-secret")

	mockDB.EXPECT().
		GetByUsername(mock.Anything, "nonexistent").
		Return(nil, fmt.Errorf("пользователь nonexistent не найден"))

	token, err := svc.Login(context.Background(), "nonexistent", "password")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "db.GetByUsername")
}

// TestAuthSvc_ValidateToken - тесты для метода ValidateToken
func TestAuthSvc_ValidateToken_OK(t *testing.T) {
	mockDB := mocks.NewDatabase(t)
	svc := New(mockDB, "test-secret")

	claims := &models.JWTClaims{
		UserID:   1,
		Username: "admin123",
		Role:     "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte("test-secret"))

	resultClaims, err := svc.ValidateToken(tokenStr)

	assert.NoError(t, err)
	assert.NotNil(t, resultClaims)
	assert.Equal(t, 1, resultClaims.UserID)
	assert.Equal(t, "admin123", resultClaims.Username)
	assert.Equal(t, "admin", resultClaims.Role)
}

func TestAuthSvc_ValidateToken_ErrInvalidToken(t *testing.T) {
	mockDB := mocks.NewDatabase(t)
	svc := New(mockDB, "test-secret")

	claims, err := svc.ValidateToken("invalid-token")

	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Contains(t, err.Error(), "jwt.ParseWithClaims")
}

func TestAuthSvc_ValidateToken_ErrExpired(t *testing.T) {
	mockDB := mocks.NewDatabase(t)
	svc := New(mockDB, "test-secret")

	claims := &models.JWTClaims{
		UserID:   1,
		Username: "admin123",
		Role:     "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte("test-secret"))

	resultClaims, err := svc.ValidateToken(tokenStr)

	assert.Error(t, err)
	assert.Nil(t, resultClaims)
}

func TestAuthSvc_ValidateToken_ErrWrongSecret(t *testing.T) {
	mockDB := mocks.NewDatabase(t)
	svc := New(mockDB, "test-secret")

	claims := &models.JWTClaims{
		UserID:   1,
		Username: "admin123",
		Role:     "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte("wrong-secret"))

	resultClaims, err := svc.ValidateToken(tokenStr)

	assert.Error(t, err)
	assert.Nil(t, resultClaims)
}