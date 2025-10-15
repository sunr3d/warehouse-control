package postgres

import (
	"context"
	"fmt"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"

	"github.com/sunr3d/warehouse-control/internal/interfaces/infra"
	"github.com/sunr3d/warehouse-control/models"
)

const (
	qCreateItem = `
	INSERT INTO items (item_name, item_description, quantity) 
	VALUES ($1, $2, $3) 
	RETURNING id`

	qListItems = `
	SELECT id, item_name, item_description, quantity, created_at, updated_at
	FROM items`

	qUpdateItem = `
	UPDATE items SET item_name = $2, item_description = $3, quantity = $4 
	WHERE id = $1`

	qDeleteItem = `
	DELETE FROM items 
	WHERE id = $1`

	qSetUserID = `
	SET LOCAL warehouse.user_id = $1`
)

var _ infra.ItemRepo = (*itemRepo)(nil)

type itemRepo struct {
	db *dbpg.DB
}

// Create - метод для создания нового item в БД.
func (r *itemRepo) Create(ctx context.Context, userID int, item *models.Item) (int, error) {
	tx, err := r.db.Master.BeginTx(ctx, nil)
	if err != nil {
		zlog.Logger.Error().Err(err).Msgf("Create: не удалось начать транзакцию для item: %v", item)
		return 0, fmt.Errorf("не удалось начать транзакцию для item: %v: %w", item, err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(
		ctx,
		qSetUserID,
		userID,
	)
	if err != nil {
		zlog.Logger.Error().Err(err).Msgf("Create: не удалось установить userID: %d", userID)

		return 0, fmt.Errorf("не удалось установить userID: %d: %w", userID, err)
	}

	row := tx.QueryRowContext(
		ctx,
		qCreateItem,
		item.Name,
		item.Description,
		item.Quantity,
	)
	var id int
	if err := row.Scan(&id); err != nil {
		zlog.Logger.Error().Err(err).Msgf("Create: не удалось перевести данные из строки в структуру для item: %v", item)

		return 0, fmt.Errorf("не удалось перевести данные из строки в структуру для item: %v: %w", item, err)
	}

	if err := tx.Commit(); err != nil {
		zlog.Logger.Error().Err(err).Msgf("Create: не удалось завершить транзакцию для item: %v", item)

		return 0, fmt.Errorf("не удалось завершить транзакцию для item: %v: %w", item, err)
	}

	return id, nil
}

// List - метод для получения всех items из БД.
func (r *itemRepo) List(ctx context.Context) ([]models.Item, error) {
	strategy := retry.Strategy{
		Attempts: 3,
	}

	rows, err := r.db.QueryWithRetry(
		ctx,
		strategy,
		qListItems,
	)
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Msgf("List: не удалось выполнить запрос List")

		return nil, fmt.Errorf("не удалось выполнить запрос List: %w", err)
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.Quantity,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			zlog.Logger.Error().
				Err(err).
				Msgf("List: не удалось перевести данные из строки в структуру для item")

			return nil, fmt.Errorf("не удалось перевести данные из строки в структуру для item: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		zlog.Logger.Error().
			Err(err).
			Msgf("List: не удалось получить все строки")

		return nil, fmt.Errorf("не удалось получить все строки: %w", err)
	}

	return items, nil
}

// Update - метод для обновления item в БД.
func (r *itemRepo) Update(ctx context.Context, userID, id int, item *models.Item) error {
	tx, err := r.db.Master.BeginTx(ctx, nil)
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Msgf("Update: не удалось начать транзакцию для item: %v", item)

		return fmt.Errorf("не удалось начать транзакцию для item: %v: %w", item, err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(
		ctx,
		qSetUserID,
		userID,
	)
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Msgf("Update: не удалось установить userID: %d", userID)

		return fmt.Errorf("не удалось установить userID: %d: %w", userID, err)
	}

	result, err := tx.ExecContext(
		ctx,
		qUpdateItem,
		id,
		item.Name,
		item.Description,
		item.Quantity,
	)
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Msgf("Update: не удалось выполнить запрос Update для item: %v", item)

		return fmt.Errorf("не удалось выполнить запрос Update для item: %v: %w", item, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Msgf("Update: не удалось получить количество строк, обновленных запросом Update для item: %v", item)

		return fmt.Errorf("не удалось получить количество строк, обновленных запросом Update для item: %v: %w", item, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("item с id %d не найден", id)
	}

	if err := tx.Commit(); err != nil {
		zlog.Logger.Error().
			Err(err).
			Msgf("Update: не удалось завершить транзакцию для item: %v", item)

		return fmt.Errorf("не удалось завершить транзакцию для item: %v: %w", item, err)
	}

	return nil
}

// Delete - метод для удаления item из БД.
func (r *itemRepo) Delete(ctx context.Context, userID, id int) error {
	tx, err := r.db.Master.BeginTx(ctx, nil)
	if err != nil {
		zlog.Logger.Error().Err(err).Msgf("Delete: не удалось начать транзакцию для item: %v", id)

		return fmt.Errorf("не удалось начать транзакцию для item: %v: %w", id, err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(
		ctx,
		qSetUserID,
		userID,
	)
	if err != nil {
		zlog.Logger.Error().Err(err).Msgf("Delete: не удалось установить userID: %d", userID)

		return fmt.Errorf("не удалось установить userID: %d: %w", userID, err)
	}

	result, err := tx.ExecContext(
		ctx,
		qDeleteItem,
		id,
	)
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Msgf("Delete: не удалось выполнить запрос Delete для item: %v", id)

		return fmt.Errorf("не удалось выполнить запрос Delete для item: %v: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Msgf("Delete: не удалось получить количество строк, удаленных запросом Delete для item: %v", id)

		return fmt.Errorf("не удалось получить количество строк, удаленных запросом Delete для item: %d: %w", id, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("item с id %d не найден", id)
	}

	if err := tx.Commit(); err != nil {
		zlog.Logger.Error().Err(err).Msgf("Delete: не удалось завершить транзакцию для item: %v", id)

		return fmt.Errorf("не удалось завершить транзакцию для item: %v: %w", id, err)
	}

	return nil
}
