package httphandlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

// getItemHistory - ручка для получения истории изменений конкретного itemID.
func (h *handler) getItemHistory(c *ginext.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Int("id", id).
			Msg("getItemHistory: некорректный запрос")
		c.JSON(
			http.StatusBadRequest,
			ginext.H{"error": "некорректный запрос: " + err.Error()},
		)
		return
	}

	history, err := h.invSvc.GetItemHistory(c.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "не найден") {
			zlog.Logger.Error().
				Err(err).
				Int("id", id).
				Msg("getItemHistory: история изменений не найдена")
			c.JSON(http.StatusNotFound, ginext.H{"error": "история изменений для item с id " + strconv.Itoa(id) + " не найдена"})
			return
		}
		zlog.Logger.Error().
			Err(err).
			Int("id", id).
			Msg("getItemHistory: не удалось получить историю изменений")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "не удалось получить историю изменений"})
		return
	}

	var resp getItemHistoryResp
	resp.ItemID = id
	for _, item := range history {
		itemHist := itemHistoryResp{
			UserID:    item.UserID,
			Operation: item.Operation,
			ChangedAt: item.ChangedAt.Format(time.RFC3339),
		}
		if item.OldValue != nil {
			itemHist.OldValue = *item.OldValue
		}
		if item.NewValue != nil {
			itemHist.NewValue = *item.NewValue
		}

		resp.Items = append(resp.Items, itemHist)
	}

	c.JSON(http.StatusOK, resp)
}
