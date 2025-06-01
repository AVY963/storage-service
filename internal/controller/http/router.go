package http

import (
	"tages/internal/auth"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter настраивает роутер для HTTP сервера
func SetupRouter(fileHandler *FileHandler, authHandler *AuthHandler, tokenManager *auth.TokenManager) *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
            return true // Разрешить любые Origin
        },
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
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

	// Статические файлы фронтенда (собранное React-приложение)
	router.Static("/assets", "./web/dist/assets")
	router.StaticFile("/vite.svg", "./web/dist/vite.svg")

	// Обработка всех остальных маршрутов - отдаем index.html для SPA
	router.NoRoute(func(c *gin.Context) {
		c.File("./web/dist/index.html")
	})

	return router
}
