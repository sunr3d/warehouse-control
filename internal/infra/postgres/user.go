package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"

	"github.com/sunr3d/warehouse-control/internal/interfaces/infra"
	"github.com/sunr3d/warehouse-control/models"
)

const (
	qGetByUsername = `
	SELECT id, username, password_hash, user_role
	FROM users
	WHERE username = $1`
)

var _ infra.UserRepo = (*userRepo)(nil)

type userRepo struct {
	db *dbpg.DB
}

// GetByUsername - метод для получения пользователя по username.
func (r *userRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	strategy := retry.Strategy{
		Attempts: 3,
	}
	var user models.User

	row, err := r.db.QueryRowWithRetry(
		ctx,
		strategy,
		qGetByUsername,
		username,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("пользователь %s не найден", username)
		}
		zlog.Logger.Error().Err(err).Msgf("GetByUsername: пользователь %s не найден", username)

		return nil, fmt.Errorf("не удалось выполнить запрос GetByUsername для пользователя %s: %w", username, err)
	}

	if err := row.Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
	); err != nil {
		zlog.Logger.Error().Err(err).Msgf("GetByUsername: не удалось перевести данные из строки в структуру для пользователя %s", username)

		return nil, fmt.Errorf("не удалось перевести данные из строки в структуру для пользователя %s: %w", username, err)
	}

	return &user, nil
}
