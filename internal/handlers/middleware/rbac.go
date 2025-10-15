package middleware

import (
	"net/http"
	"slices"

	"github.com/wb-go/wbf/ginext"

	"github.com/sunr3d/warehouse-control/models"
)

// RBACMiddleware - middleware для проверки прав доступа.
// Проверяет наличие роли пользователя в списке разрешенных ролей.
// Если роль не найдена, возвращает ошибку 403 Forbidden.
// Если роль найдена, пропускает запрос дальше.
func RBACMiddleware(roles ...string) ginext.HandlerFunc {
	return func(c *ginext.Context) {
		userClaims, exists := c.Get(UserCtxKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ginext.H{"error": "пользователь не авторизован"})
			return
		}

		claims, ok := userClaims.(*models.JWTClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ginext.H{"error": "неверный тип claims"})
			return
		}

		if !slices.Contains(roles, claims.Role) {
			c.AbortWithStatusJSON(http.StatusForbidden, ginext.H{"error": "недостаточно прав"})
			return
		}

		c.Next()
	}
}
