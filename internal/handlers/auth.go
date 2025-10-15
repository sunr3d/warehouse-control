package httphandlers

import (
	"net/http"

	"github.com/wb-go/wbf/ginext"
)

// login - handler для авторизации пользователя.
func (h *handler) login(c *ginext.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "некорректный запрос"})
		return
	}

	token, err := h.authSvc.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ginext.H{"error": "неверные учетные данные"})
		return
	}

	resp := loginResp{
		Username: req.Username,
		Token:    token,
	}

	c.JSON(http.StatusOK, resp)
}
