package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/zlog"

	"github.com/sunr3d/warehouse-control/internal/config"
	"github.com/sunr3d/warehouse-control/internal/interfaces/infra"
)

const (
	pingTimeout = 3 * time.Second
)

var _ infra.Database = (*postgresRepo)(nil)

type postgresRepo struct {
	*userRepo
	*itemRepo
	*itemHistoryRepo
}

// New - конструктор нового postgresRepo.
// Создает новое соединение с БД, пингует его, собирает репозитории и возвращает их.
func New(ctx context.Context, cfg config.DBConfig) (infra.Database, error) {
	options := &dbpg.Options{
		MaxOpenConns: cfg.MaxOpenConns,
		MaxIdleConns: cfg.MaxIdleConns,
	}

	db, err := dbpg.New(cfg.DSN, nil, options)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("не удалось создать БД")
		return nil, fmt.Errorf("не удалось создать БД: %w", err)
	}

	pCtx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	if err := db.Master.PingContext(pCtx); err != nil {
		zlog.Logger.Error().Err(err).Msg("ошибка при подключении к БД")
		if err := db.Master.Close(); err != nil {
			zlog.Logger.Error().Err(err).Msg("ошибка при закрытии соединения с БД")
		}
		return nil, fmt.Errorf("не удалось подключиться к БД: %w", err)
	}

	userRepo := &userRepo{db: db}
	itemRepo := &itemRepo{db: db}
	itemHistoryRepo := &itemHistoryRepo{db: db}

	return &postgresRepo{
		userRepo:        userRepo,
		itemRepo:        itemRepo,
		itemHistoryRepo: itemHistoryRepo,
	}, nil
}

// Close - метод для закрытия соединения с БД.
// Использовать через тайпкаст к io.Closer через defer после создания нового postgresRepo.
func (r *postgresRepo) Close() error {
	if err := r.userRepo.db.Master.Close(); err != nil {
		zlog.Logger.Error().Err(err).Msg("не удалось закрыть соединение с БД")
		return err
	}

	return nil
}
