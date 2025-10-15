package infra

import (
	"context"

	"github.com/sunr3d/warehouse-control/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.2 --name=Database --output=../../../mocks --filename=mock_database.go --with-expecter
type Database interface {
	UserRepo
	ItemRepo
	ItemHistoryRepo
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.2 --name=UserRepo --output=../../../mocks --filename=mock_user_repo.go --with-expecter
type UserRepo interface {
	GetByUsername(ctx context.Context, username string) (*models.User, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.2 --name=ItemRepo --output=../../../mocks --filename=mock_item_repo.go --with-expecter
type ItemRepo interface {
	Create(ctx context.Context, userID int, item *models.Item) (int, error)
	List(ctx context.Context) ([]models.Item, error)
	Update(ctx context.Context, userID, id int, item *models.Item) error
	Delete(ctx context.Context, userID, id int) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.2 --name=ItemHistoryRepo --output=../../../mocks --filename=mock_item_history_repo.go --with-expecter
type ItemHistoryRepo interface {
	GetByItemID(ctx context.Context, itemID int) ([]models.ItemHistory, error)
}
