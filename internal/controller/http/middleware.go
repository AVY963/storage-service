package http

import (
	"log"
	"net/http"
	"strings"

	"tages/internal/auth"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware middleware для проверки JWT токена
type AuthMiddleware struct {
	tokenManager *auth.TokenManager
}

// NewAuthMiddleware создает новое middleware для авторизации
func NewAuthMiddleware(tokenManager *auth.TokenManager) *AuthMiddleware {
	return &AuthMiddleware{
		tokenManager: tokenManager,
	}
}

// Middleware проверяет JWT токен в заголовке Authorization
func (m *AuthMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string

		// 1. Попытка взять из заголовка Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}

		// 2. Если не найден — пробуем взять из cookie
		if token == "" {
			cookie, err := c.Cookie("access_token")
			if err == nil {
				token = cookie
			}
		}

		// 3. Если токена всё ещё нет — 401
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "токен авторизации не найден",
			})
			return
		}

		// 4. Парсим токен
		claims, err := m.tokenManager.ParseAccessToken(token)
		if err != nil {
			if err == auth.ErrExpiredToken {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "срок действия токена истек",
				})
			} else {
				log.Printf("ERROR: Invalid token: %v", err)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "неверный токен",
				})
			}
			return
		}

		// 5. Добавляем ID в контекст
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
// GetUserID извлекает ID пользователя из контекста
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, false
	}

	id, ok := userID.(uint)
	return id, ok
}
