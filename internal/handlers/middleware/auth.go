package middleware

import (
	"net/http"
	"strings"

	"github.com/wb-go/wbf/ginext"

	"github.com/sunr3d/warehouse-control/internal/interfaces/services"
)

const (
	UserCtxKey = "user"
)

// AuthMiddleware - middleware для авторизации пользователя.
// Валидирует токен из заголовка Authorization и устанавливает claims в контекст.
// Если токен не валиден, возвращает ошибку 401 Unauthorized.
// Если токен валиден, устанавливает claims в контекст и пропускает запрос дальше.
func AuthMiddleware(authSvc services.AuthService) ginext.HandlerFunc {
	return func(c *ginext.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ginext.H{"error": "заголовок авторизации не найден"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ginext.H{"error": "неверный формат заголовка авторизации"})
			return
		}

		tokenStr := parts[1]
		claims, err := authSvc.ValidateToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ginext.H{"error": "невалидный токен"})
			return
		}

		c.Set(UserCtxKey, claims)
		c.Next()
	}
}
