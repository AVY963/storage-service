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
		// Извлекаем токен из заголовка
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "отсутствует заголовок авторизации",
			})
			return
		}

		// Проверяем формат токена
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "неверный формат токена",
			})
			return
		}

		token := headerParts[1]

		// Проверяем валидность токена
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

		// Добавляем ID пользователя в контекст запроса
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
