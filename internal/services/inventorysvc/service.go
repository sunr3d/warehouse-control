package inventorysvc

import (
	"context"
	"fmt"
	"strings"

	"github.com/sunr3d/warehouse-control/internal/interfaces/infra"
	"github.com/sunr3d/warehouse-control/internal/interfaces/services"
	"github.com/sunr3d/warehouse-control/models"
)

var _ services.InventoryService = (*inventorySvc)(nil)

type inventorySvc struct {
	db infra.Database
}

// New - конструктор нового inventorySvc.
func New(db infra.Database) services.InventoryService {
	return &inventorySvc{db: db}
}

// AddItem - метод для добавления нового item в БД.
func (s *inventorySvc) AddItem(ctx context.Context, userID int, item *models.Item) (int, error) {
	if item.Quantity <= 0 {
		return 0, fmt.Errorf("quantity должно быть больше 0")
	}

	id, err := s.db.Create(ctx, userID, item)
	if err != nil {
		return 0, fmt.Errorf("db.Create: %w", err)
	}

	return id, nil
}

// GetInventory - метод для получения всех items из БД.
func (s *inventorySvc) GetInventory(ctx context.Context) ([]models.Item, error) {
	items, err := s.db.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("db.List: %w", err)
	}

	return items, nil
}

// UpdateItem - метод для обновления item в БД.
func (s *inventorySvc) UpdateItem(ctx context.Context, userID, id int, item *models.Item) error {
	if item.Quantity < 0 {
		return fmt.Errorf("quantity должно быть больше или равно 0")
	}

	if err := s.db.Update(ctx, userID, id, item); err != nil {
		if strings.Contains(err.Error(), "не найден") {
			return fmt.Errorf("item с id %d не найден", id)
		}

		return fmt.Errorf("db.Update: %w", err)
	}

	return nil
}

// DeleteItem - метод для удаления item из БД.
func (s *inventorySvc) DeleteItem(ctx context.Context, userID, id int) error {
	if err := s.db.Delete(ctx, userID, id); err != nil {
		if strings.Contains(err.Error(), "не найден") {
			return fmt.Errorf("item с id %d не найден", id)
		}

		return fmt.Errorf("db.Delete: %w", err)
	}

	return nil
}

// GetItemHistory - метод для получения истории изменений для конкретного itemID.
func (s *inventorySvc) GetItemHistory(ctx context.Context, id int) ([]models.ItemHistory, error) {
	history, err := s.db.GetByItemID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("db.GetByItemID: %w", err)
	}

	if len(history) == 0 {
		return nil, fmt.Errorf("история изменений для item с id %d не найдена", id)
	}

	return history, nil
}
