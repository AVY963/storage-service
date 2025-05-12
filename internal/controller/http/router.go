package http

import (
	"tages/internal/auth"

	"github.com/gin-gonic/gin"
)

// SetupRouter настраивает роутер для HTTP сервера
func SetupRouter(fileHandler *FileHandler, authHandler *AuthHandler, tokenManager *auth.TokenManager) *gin.Engine {
	router := gin.Default()

	// Создаем middleware для авторизации
	authMiddleware := NewAuthMiddleware(tokenManager)

	// API группа
	api := router.Group("/api")

	// Маршруты для авторизации (публичные)
	authRoutes := api.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.RegisterHandler)
		authRoutes.POST("/login", authHandler.LoginHandler)
		authRoutes.POST("/refresh", authHandler.RefreshTokenHandler)
		authRoutes.POST("/logout", authHandler.LogoutHandler)
	}

	// Маршруты для работы с файлами (защищенные)
	filesRoutes := api.Group("/files")
	filesRoutes.Use(authMiddleware.Middleware())
	{
		filesRoutes.POST("/upload", fileHandler.UploadHandler)
		filesRoutes.GET("/list", fileHandler.ListHandler)
		filesRoutes.GET("/download/:filename", fileHandler.DownloadHandler)
		filesRoutes.DELETE("/delete/:filename", fileHandler.DeleteHandler)
	}

	return router
}
