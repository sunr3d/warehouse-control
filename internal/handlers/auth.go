package httphandlers

import (
	"net/http"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

// login - handler для авторизации пользователя.
func (h *handler) login(c *ginext.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Logger.Warn().
			Err(err).
			Msg("login: некорректный JSON запрос")
		c.JSON(http.StatusBadRequest, ginext.H{"error": "некорректный запрос"})
		return
	}

	zlog.Logger.Info().
		Str("username", req.Username).
		Msg("login: попытка авторизации пользователя")

	token, err := h.authSvc.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Msg("login: неверные учетные данные")
		c.JSON(http.StatusUnauthorized, ginext.H{"error": "неверные учетные данные"})
		return
	}

	zlog.Logger.Info().
		Str("username", req.Username).
		Msg("login: успешная авторизация пользователя")

	resp := loginResp{
		Username: req.Username,
		Token:    token,
	}

	c.JSON(http.StatusOK, resp)
}
