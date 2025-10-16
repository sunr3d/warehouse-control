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
	qGetByItemID = `
	SELECT id, item_id, user_id, operation, old_value, new_value, changed_at
	FROM items_history
	WHERE item_id = $1`
)

var _ infra.ItemHistoryRepo = (*itemHistoryRepo)(nil)

type itemHistoryRepo struct {
	db *dbpg.DB
}

// GetByItemID - метод для получения истории изменений для конкретного itemID.
func (r *itemHistoryRepo) GetByItemID(ctx context.Context, itemID int) ([]models.ItemHistory, error) {
	strategy := retry.Strategy{
		Attempts: 3,
	}

	rows, err := r.db.QueryWithRetry(
		ctx,
		strategy,
		qGetByItemID,
		itemID,
	)
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Int("item_id", itemID).
			Msg("GetByItemID: не удалось выполнить запрос GetByItemID")

		return nil, fmt.Errorf("не удалось выполнить запрос GetByItemID: %w", err)
	}
	defer rows.Close()

	var ledger []models.ItemHistory
	for rows.Next() {
		var itemHistory models.ItemHistory
		if err := rows.Scan(
			&itemHistory.ID,
			&itemHistory.ItemID,
			&itemHistory.UserID,
			&itemHistory.Operation,
			&itemHistory.OldValue,
			&itemHistory.NewValue,
			&itemHistory.ChangedAt,
		); err != nil {
			zlog.Logger.Error().
				Err(err).
				Int("item_id", itemID).
				Msg("GetByItemID: не удалось перевести данные из строки в структуру")

			return nil, fmt.Errorf("не удалось перевести данные из строки в структуру: %w", err)
		}

		ledger = append(ledger, itemHistory)
	}

	if err := rows.Err(); err != nil {
		zlog.Logger.Error().
			Err(err).
			Int("item_id", itemID).
			Msg("GetByItemID: не удалось получить все строки")

		return nil, fmt.Errorf("не удалось получить все строки: %w", err)
	}

	return ledger, nil
}
