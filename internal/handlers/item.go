package httphandlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"

	"github.com/sunr3d/warehouse-control/models"
)

// createItem - handler для создания нового item.
func (h *handler) createItem(c *ginext.Context) {
	userClaims, _ := c.Get("user")
	claims := userClaims.(*models.JWTClaims)
	userID := claims.UserID

	var req itemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Logger.Warn().
			Err(err).
			Msg("createItem: некорректный JSON запрос")
		c.JSON(http.StatusBadRequest, ginext.H{"error": "некорректный запрос"})
		return
	}

	zlog.Logger.Info().
		Int("user_id", userID).
		Str("item_name", req.Name).
		Int("quantity", req.Quantity).
		Msg("createItem: попытка создания нового item")

	item := &models.Item{
		Name:        req.Name,
		Description: req.Description,
		Quantity:    req.Quantity,
	}

	id, err := h.invSvc.AddItem(c.Request.Context(), userID, item)
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Int("user_id", userID).
			Str("item_name", req.Name).
			Int("quantity", req.Quantity).
			Msg("createItem: не удалось создать item")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "не удалось создать item"})
		return
	}

	zlog.Logger.Info().
		Int("user_id", userID).
		Int("item_id", id).
		Str("item_name", req.Name).
		Int("quantity", req.Quantity).
		Msg("createItem: item успешно создан")

	c.JSON(http.StatusCreated, ginext.H{"id": id})
}

// getItems - handler для получения всех items.
func (h *handler) getItems(c *ginext.Context) {
	userClaims, _ := c.Get("user")
	claims := userClaims.(*models.JWTClaims)
	userID := claims.UserID

	zlog.Logger.Info().
		Int("user_id", userID).
		Msg("getItems: попытка получить все items")

	items, err := h.invSvc.GetInventory(c.Request.Context())
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Int("user_id", userID).
			Msg("getItems: не удалось получить items")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "не удалось получить items"})
		return
	}

	zlog.Logger.Info().
		Int("user_id", userID).
		Int("items_count", len(items)).
		Msg("getItems: все items успешно получены")

	var resp []itemResp
	for _, item := range items {
		resp = append(resp, itemResp{
			ID:          item.ID,
			Quantity:    item.Quantity,
			Name:        item.Name,
			Description: item.Description,
			CreatedAt:   item.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   item.UpdatedAt.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, resp)
}

// updateItem - handler для обновления item.
func (h *handler) updateItem(c *ginext.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		zlog.Logger.Warn().
			Err(err).
			Msg("updateItem: некорректный запрос")
		c.JSON(http.StatusBadRequest, ginext.H{"error": "некорректный запрос"})
		return
	}

	userClaims, _ := c.Get("user")
	claims := userClaims.(*models.JWTClaims)
	userID := claims.UserID

	var req itemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Logger.Warn().
			Err(err).
			Msg("updateItem: некорректный JSON запрос")
		c.JSON(http.StatusBadRequest, ginext.H{"error": "некорректный запрос"})
		return
	}

	zlog.Logger.Info().
		Int("user_id", userID).
		Int("item_id", id).
		Str("item_name", req.Name).
		Int("quantity", req.Quantity).
		Msg("updateItem: попытка обновления item")

	item := &models.Item{
		Name:        req.Name,
		Description: req.Description,
		Quantity:    req.Quantity,
		UpdatedAt:   time.Now(),
	}

	err = h.invSvc.UpdateItem(c.Request.Context(), userID, id, item)
	if err != nil {
		if strings.Contains(err.Error(), "не найден") {
			zlog.Logger.Warn().
				Err(err).
				Int("user_id", userID).
				Int("item_id", id).
				Msg("updateItem: item не найден")
			c.JSON(http.StatusNotFound, ginext.H{"error": "item с id " + strconv.Itoa(id) + " не найден"})
			return
		}
		zlog.Logger.Error().
			Err(err).
			Int("user_id", userID).
			Int("item_id", id).
			Msg("updateItem: не удалось обновить item")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "не удалось обновить item"})
		return
	}

	zlog.Logger.Info().
		Int("user_id", userID).
		Int("item_id", id).
		Str("item_name", req.Name).
		Int("quantity", req.Quantity).
		Msg("updateItem: item успешно обновлен")

	c.JSON(http.StatusOK, ginext.H{"id": id, "message": "item успешно обновлен"})
}

// deleteItem - handler для удаления item.
func (h *handler) deleteItem(c *ginext.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		zlog.Logger.Warn().
			Err(err).
			Msg("deleteItem: некорректный запрос")
		c.JSON(http.StatusBadRequest, ginext.H{"error": "некорректный запрос"})
		return
	}

	userClaims, _ := c.Get("user")
	claims := userClaims.(*models.JWTClaims)
	userID := claims.UserID

	err = h.invSvc.DeleteItem(c.Request.Context(), userID, id)
	if err != nil {
		if strings.Contains(err.Error(), "не найден") {
			zlog.Logger.Warn().
				Err(err).
				Int("user_id", userID).
				Int("item_id", id).
				Msg("deleteItem: item не найден")
			c.JSON(http.StatusNotFound, ginext.H{"error": "item с id " + strconv.Itoa(id) + " не найден"})
			return
		}
		zlog.Logger.Error().
			Err(err).
			Int("user_id", userID).
			Int("item_id", id).
			Msg("deleteItem: не удалось удалить item")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "не удалось удалить item"})
		return
	}

	zlog.Logger.Info().
		Int("user_id", userID).
		Int("item_id", id).
		Msg("deleteItem: item успешно удален")

	c.JSON(http.StatusOK, ginext.H{"id": id, "message": "item успешно удален"})
}
