package services

import (
	"context"

	"github.com/sunr3d/warehouse-control/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.2 --name=AuthService --output=../../../mocks --filename=mock_auth_service.go --with-expecter
type AuthService interface {
	Login(ctx context.Context, username, pass string) (string, error)
	ValidateToken(tokenStr string) (*models.JWTClaims, error)
}
