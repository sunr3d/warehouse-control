package inventorysvc

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/sunr3d/warehouse-control/mocks"
	"github.com/sunr3d/warehouse-control/models"
)

// TestInventorySvc_AddItem - тесты для метода AddItem
func TestInventorySvc_AddItem_OK(t *testing.T) {
	mockDB := mocks.NewDatabase(t)
	svc := New(mockDB)

	item := &models.Item{
		Name:        "Товар 1",
		Description: "Описание товара",
		Quantity:    10,
	}

	mockDB.EXPECT().
		Create(mock.Anything, 1, item).
		Return(1, nil)

	id, err := svc.AddItem(context.Background(), 1, item)

	assert.NoError(t, err)
	assert.Equal(t, 1, id)
}

func TestInventorySvc_AddItem_ErrZeroQuantity(t *testing.T) {
	mockDB := mocks.NewDatabase(t)
	svc := New(mockDB)

	item := &models.Item{
		Name:        "Товар 1",
		Description: "Описание товара",
		Quantity:    0,
	}

	id, err := svc.AddItem(context.Background(), 1, item)

	assert.Error(t, err)
	assert.Equal(t, 0, id)
	assert.Contains(t, err.Error(), "quantity должно быть больше 0")
}

func TestInventorySvc_AddItem_ErrNegativeQuantity(t *testing.T) {
	mockDB := mocks.NewDatabase(t)
	svc := New(mockDB)

	item := &models.Item{
		Name:        "Товар 1",
		Description: "Описание товара",
		Quantity:    -5,
	}

	id, err := svc.AddItem(context.Background(), 1, item)

	assert.Error(t, err)
	assert.Equal(t, 0, id)
	assert.Contains(t, err.Error(), "quantity должно быть больше 0")
}

func TestInventorySvc_AddItem_ErrDBFailed(t *testing.T) {
	mockDB := mocks.NewDatabase(t)
	svc := New(mockDB)

	item := &models.Item{
		Name:        "Товар 1",
		Description: "Описание товара",
		Quantity:    10,
	}

	mockDB.EXPECT().
		Create(mock.Anything, 1, item).
		Return(0, fmt.Errorf("database error"))

	id, err := svc.AddItem(context.Background(), 1, item)

	assert.Error(t, err)
	assert.Equal(t, 0, id)
	assert.Contains(t, err.Error(), "db.Create")
}

// TestInventorySvc_GetInventory - тесты для метода GetInventory
func TestInventorySvc_GetInventory_OK(t *testing.T) {
	mockDB := mocks.NewDatabase(t)
	svc := New(mockDB)

	expectedItems := []models.Item{
		{ID: 1, Name: "Товар 1", Description: "Описание 1", Quantity: 10},
		{ID: 2, Name: "Товар 2", Description: "Описание 2", Quantity: 20},
	}

	mockDB.EXPECT().
		List(mock.Anything).
		Return(expectedItems, nil)

	items, err := svc.GetInventory(context.Background())

	assert.NoError(t, err)
	assert.Len(t, items, 2)
	assert.Equal(t, expectedItems, items)
}

func TestInventorySvc_GetInventory_ErrDBFailed(t *testing.T) {
	mockDB := mocks.NewDatabase(t)
	svc := New(mockDB)

	mockDB.EXPECT().
		List(mock.Anything).
		Return(nil, fmt.Errorf("database error"))

	items, err := svc.GetInventory(context.Background())

	assert.Error(t, err)
	assert.Nil(t, items)
	assert.Contains(t, err.Error(), "db.List")
}

func TestInventorySvc_GetInventory_EmptyList(t *testing.T) {
	mockDB := mocks.NewDatabase(t)
	svc := New(mockDB)

	mockDB.EXPECT().
		List(mock.Anything).
		Return([]models.Item{}, nil)

	items, err := svc.GetInventory(context.Background())

	assert.NoError(t, err)
	assert.Len(t, items, 0)
}

// TestInventorySvc_UpdateItem - тесты для метода UpdateItem
func TestInventorySvc_UpdateItem_OK(t *testing.T) {
	mockDB := mocks.NewDatabase(t)
	svc := New(mockDB)

	item := &models.Item{
		Name:        "Обновленный товар",
		Description: "Обновленное описание",
		Quantity:    15,
	}

	mockDB.EXPECT().
		Update(mock.Anything, 1, 1, item).
		Return(nil)

	err := svc.UpdateItem(context.Background(), 1, 1, item)

	assert.NoError(t, err)
}

func TestInventorySvc_UpdateItem_OKZeroQuantity(t *testing.T) {
	mockDB := mocks.NewDatabase(t)
	svc := New(mockDB)

	item := &models.Item{
		Name:        "Обновленный товар",
		Description: "Обновленное описание",
		Quantity:    0,
	}

	mockDB.EXPECT().
		Update(mock.Anything, 1, 1, item).
		Return(nil)

	err := svc.UpdateItem(context.Background(), 1, 1, item)

	assert.NoError(t, err)
}

func TestInventorySvc_UpdateItem_ErrNegativeQuantity(t *testing.T) {
	mockDB := mocks.NewDatabase(t)
	svc := New(mockDB)

	item := &models.Item{
		Name:        "Обновленный товар",
		Description: "Обновленное описание",
		Quantity:    -5,
	}

	err := svc.UpdateItem(context.Background(), 1, 1, item)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "quantity должно быть больше или равно 0")
}

func TestInventorySvc_UpdateItem_ErrItemNotFound(t *testing.T) {
	mockDB := mocks.NewDatabase(t)
	svc := New(mockDB)

	item := &models.Item{
		Name:        "Обновленный товар",
		Description: "Обновленное описание",
		Quantity:    15,
	}

	mockDB.EXPECT().
		Update(mock.Anything, 1, 999, item).
		Return(fmt.Errorf("item с id 999 не найден"))

	err := svc.UpdateItem(context.Background(), 1, 999, item)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "не найден")
}

func TestInventorySvc_UpdateItem_ErrDBFailed(t *testing.T) {
	mockDB := mocks.NewDatabase(t)
	svc := New(mockDB)

	item := &models.Item{
		Name:        "Обновленный товар",
		Description: "Обновленное описание",
		Quantity:    15,
	}

	mockDB.EXPECT().
		Update(mock.Anything, 1, 1, item).
		Return(fmt.Errorf("database connection error"))

	err := svc.UpdateItem(context.Background(), 1, 1, item)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db.Update")
}

// TestInventorySvc_DeleteItem - тесты для метода DeleteItem
func TestInventorySvc_DeleteItem_OK(t *testing.T) {
	mockDB := mocks.NewDatabase(t)
	svc := New(mockDB)

	mockDB.EXPECT().
		Delete(mock.Anything, 1, 1).
		Return(nil)

	err := svc.DeleteItem(context.Background(), 1, 1)

	assert.NoError(t, err)
}

func TestInventorySvc_DeleteItem_ErrItemNotFound(t *testing.T) {
	mockDB := mocks.NewDatabase(t)
	svc := New(mockDB)

	mockDB.EXPECT().
		Delete(mock.Anything, 1, 999).
		Return(fmt.Errorf("item с id 999 не найден"))

	err := svc.DeleteItem(context.Background(), 1, 999)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "не найден")
}

func TestInventorySvc_DeleteItem_ErrDBFailed(t *testing.T) {
	mockDB := mocks.NewDatabase(t)
	svc := New(mockDB)

	mockDB.EXPECT().
		Delete(mock.Anything, 1, 1).
		Return(fmt.Errorf("database connection error"))

	err := svc.DeleteItem(context.Background(), 1, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db.Delete")
}

// TestInventorySvc_GetItemHistory - тесты для метода GetItemHistory
func TestInventorySvc_GetItemHistory_OK(t *testing.T) {
	mockDB := mocks.NewDatabase(t)
	svc := New(mockDB)

	oldValue := `{"id":1,"name":"Старое название","quantity":10}`
	newValue := `{"id":1,"name":"Новое название","quantity":15}`

	expectedHistory := []models.ItemHistory{
		{ID: 1, ItemID: 1, UserID: 1, Operation: "INSERT", OldValue: nil, NewValue: &newValue},
		{ID: 2, ItemID: 1, UserID: 1, Operation: "UPDATE", OldValue: &oldValue, NewValue: &newValue},
	}

	mockDB.EXPECT().
		GetByItemID(mock.Anything, 1).
		Return(expectedHistory, nil)

	history, err := svc.GetItemHistory(context.Background(), 1)

	assert.NoError(t, err)
	assert.Len(t, history, 2)
	assert.Equal(t, expectedHistory, history)
}

func TestInventorySvc_GetItemHistory_ErrItemNotFound(t *testing.T) {
	mockDB := mocks.NewDatabase(t)
	svc := New(mockDB)

	mockDB.EXPECT().
		GetByItemID(mock.Anything, 999).
		Return([]models.ItemHistory{}, nil)

	history, err := svc.GetItemHistory(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, history)
	assert.Contains(t, err.Error(), "не найдена")
}

func TestInventorySvc_GetItemHistory_ErrDBFailed(t *testing.T) {
	mockDB := mocks.NewDatabase(t)
	svc := New(mockDB)

	mockDB.EXPECT().
		GetByItemID(mock.Anything, 1).
		Return(nil, fmt.Errorf("database error"))

	history, err := svc.GetItemHistory(context.Background(), 1)

	assert.Error(t, err)
	assert.Nil(t, history)
	assert.Contains(t, err.Error(), "db.GetByItemID")
}
