package httphandlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/wb-go/wbf/ginext"

	"github.com/sunr3d/warehouse-control/models"
)

// createItem - handler для создания нового item.
func (h *handler) createItem(c *ginext.Context) {
	userClaims, _ := c.Get("user")
	claims := userClaims.(*models.JWTClaims)
	userID := claims.UserID

	var req itemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "некорректный запрос"})
		return
	}

	item := &models.Item{
		Name:        req.Name,
		Description: req.Description,
		Quantity:    req.Quantity,
	}

	id, err := h.invSvc.AddItem(c.Request.Context(), userID, item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "не удалось создать item"})
		return
	}

	c.JSON(http.StatusCreated, ginext.H{"id": id})
}

// getItems - handler для получения всех items.
func (h *handler) getItems(c *ginext.Context) {
	items, err := h.invSvc.GetInventory(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "не удалось получить items"})
		return
	}

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
		c.JSON(http.StatusBadRequest, ginext.H{"error": "некорректный запрос"})
		return
	}

	userClaims, _ := c.Get("user")
	claims := userClaims.(*models.JWTClaims)
	userID := claims.UserID

	var req itemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "некорректный запрос"})
		return
	}

	item := &models.Item{
		Name:        req.Name,
		Description: req.Description,
		Quantity:    req.Quantity,
	}

	err = h.invSvc.UpdateItem(c.Request.Context(), userID, id, item)
	if err != nil {
		if strings.Contains(err.Error(), "не найден") {
			c.JSON(http.StatusNotFound, ginext.H{"error": "item с id " + strconv.Itoa(id) + " не найден"})
			return
		}
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "не удалось обновить item"})
		return
	}

	c.JSON(http.StatusOK, ginext.H{"id": id, "message": "item успешно обновлен"})
}

// deleteItem - handler для удаления item.
func (h *handler) deleteItem(c *ginext.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "некорректный запрос"})
		return
	}

	userClaims, _ := c.Get("user")
	claims := userClaims.(*models.JWTClaims)
	userID := claims.UserID

	err = h.invSvc.DeleteItem(c.Request.Context(), userID, id)
	if err != nil {
		if strings.Contains(err.Error(), "не найден") {
			c.JSON(http.StatusNotFound, ginext.H{"error": "item с id " + strconv.Itoa(id) + " не найден"})
			return
		}
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "не удалось удалить item"})
		return
	}

	c.JSON(http.StatusOK, ginext.H{"id": id, "message": "item успешно удален"})
}
