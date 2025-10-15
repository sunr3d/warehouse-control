package services

import (
	"context"

	"github.com/sunr3d/warehouse-control/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.2 --name=InventoryService --output=../../../mocks --filename=mock_inventory_service.go --with-expecter
type InventoryService interface {
	AddItem(ctx context.Context, userID int, item *models.Item) (int, error)
	GetInventory(ctx context.Context) ([]models.Item, error)
	UpdateItem(ctx context.Context, userID, id int, item *models.Item) error
	DeleteItem(ctx context.Context, userID, id int) error

	GetItemHistory(ctx context.Context, id int) ([]models.ItemHistory, error)
}
