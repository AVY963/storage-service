package http

import (
	"errors"
	"log"
	"net/http"
	"time"

	"tages/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	userUsecase *usecase.UserUsecase
}

func NewAuthHandler(userUsecase *usecase.UserUsecase) *AuthHandler {
	return &AuthHandler{
		userUsecase: userUsecase,
	}
}

// RegisterRequest структура для регистрации пользователя
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest структура для авторизации пользователя
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// TokenResponse структура ответа с токеном
type TokenResponse struct {
	AccessToken string   `json:"access_token"`
	User        UserInfo `json:"user"`
}

// UserInfo структура с информацией о пользователе
type UserInfo struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

// RefreshResponse структура ответа при обновлении токена
type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}

// RegisterHandler обрабатывает регистрацию нового пользователя
func (h *AuthHandler) RegisterHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Получаем детальные ошибки валидации
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			errorMap := make(map[string]string)
			for _, e := range validationErrors {
				field := e.Field()
				switch field {
				case "Email":
					errorMap[field] = "некорректный email адрес"
				case "Password":
					if e.Tag() == "min" {
						errorMap[field] = "пароль должен содержать минимум 6 символов"
					} else {
						errorMap[field] = "пароль обязателен"
					}
				default:
					errorMap[field] = "ошибка валидации"
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "ошибка валидации данных",
				"details": errorMap,
			})
		} else {
			// Общая ошибка
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "неверный формат данных: " + err.Error(),
			})
		}
		return
	}

	user, tokens, err := h.userUsecase.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		log.Printf("ERROR: Failed to register user: %v", err)
		if err == usecase.ErrUserAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{
				"error": "пользователь с таким email уже существует",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "ошибка регистрации: " + err.Error(),
			})
		}
		return
	}

	// Устанавливаем refresh token в httpOnly cookie
	setRefreshTokenCookie(c, tokens.RefreshToken)

	// Возвращаем access token и данные пользователя
	c.JSON(http.StatusOK, TokenResponse{
		AccessToken: tokens.AccessToken,
		User: UserInfo{
			ID:    user.ID,
			Email: user.Email,
		},
	})
}

// LoginHandler обрабатывает авторизацию пользователя
func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Получаем детальные ошибки валидации
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			errorMap := make(map[string]string)
			for _, e := range validationErrors {
				field := e.Field()
				switch field {
				case "Email":
					errorMap[field] = "некорректный email адрес"
				case "Password":
					errorMap[field] = "пароль обязателен"
				default:
					errorMap[field] = "ошибка валидации"
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "ошибка валидации данных",
				"details": errorMap,
			})
		} else {
			// Общая ошибка
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "неверный формат данных: " + err.Error(),
			})
		}
		return
	}

	user, tokens, err := h.userUsecase.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		log.Printf("ERROR: Failed to login: %v", err)
		if err == usecase.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "неверный email или пароль",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "ошибка авторизации: " + err.Error(),
			})
		}
		return
	}

	// Устанавливаем refresh token в httpOnly cookie
	setRefreshTokenCookie(c, tokens.RefreshToken)

	// Возвращаем access token и данные пользователя
	c.JSON(http.StatusOK, TokenResponse{
		AccessToken: tokens.AccessToken,
		User: UserInfo{
			ID:    user.ID,
			Email: user.Email,
		},
	})
}

// RefreshTokenHandler обрабатывает обновление токена
func (h *AuthHandler) RefreshTokenHandler(c *gin.Context) {
	// Получаем refresh token из cookie
	refreshToken, err := getRefreshTokenFromCookie(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "refresh token не найден",
		})
		return
	}

	// Обновляем токены
	tokens, err := h.userUsecase.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		log.Printf("ERROR: Failed to refresh token: %v", err)
		if err == usecase.ErrInvalidToken {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "неверный refresh token",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "ошибка обновления токена: " + err.Error(),
			})
		}
		return
	}

	// Устанавливаем новый refresh token в cookie
	setRefreshTokenCookie(c, tokens.RefreshToken)

	// Возвращаем новый access token
	c.JSON(http.StatusOK, RefreshResponse{
		AccessToken: tokens.AccessToken,
	})
}

// LogoutHandler обрабатывает выход из системы
func (h *AuthHandler) LogoutHandler(c *gin.Context) {
	// Получаем refresh token из cookie
	refreshToken, err := getRefreshTokenFromCookie(c)
	if err == nil {
		// Удаляем токен из базы данных
		if err := h.userUsecase.Logout(c.Request.Context(), refreshToken); err != nil {
			log.Printf("WARN: Failed to logout: %v", err)
		}
	}

	// Удаляем cookie с refresh token
	deleteRefreshTokenCookie(c)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "выход выполнен успешно",
	})
}

// Установка refresh token в cookie
func setRefreshTokenCookie(c *gin.Context, token string) {
	c.SetCookie(
		"refresh_token",
		token,
		int(7*24*time.Hour.Seconds()), // 7 дней
		"/",
		"",
		true, // Secure
		true, // HttpOnly
	)
}

// Удаление cookie с refresh token
func deleteRefreshTokenCookie(c *gin.Context) {
	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"",
		true,
		true,
	)
}

// Получение refresh token из cookie
func getRefreshTokenFromCookie(c *gin.Context) (string, error) {
	token, err := c.Cookie("refresh_token")
	if err != nil {
		return "", err
	}
	return token, nil
}
