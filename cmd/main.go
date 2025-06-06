package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tages/internal/auth"
	"tages/internal/config"
	handler "tages/internal/controller/http"
	"tages/internal/repository/pg"
	"tages/internal/usecase"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Загружаем конфигурацию
	cfg, err := config.InitConfig("")
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// Подключаемся к базе данных
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.PG.ConnectionString())
	if err != nil {
		log.Fatalf("failed to create pool: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	// Создаем репозитории и usecase
	pgRepo := pg.New(pool)

	fileUsecase := usecase.New(pgRepo)

	// Создаем менеджер JWT
	tokenManager := auth.NewTokenManager(cfg.JWT)
	userUsecase := usecase.NewUserUsecase(pgRepo, tokenManager)

	// Создаем HTTP обработчики
	fileHandler := handler.NewFileHandler(fileUsecase)
	authHandler := handler.NewAuthHandler(userUsecase)

	// Настраиваем роутер
	router := handler.SetupRouter(fileHandler, authHandler, tokenManager)

	// Создаем HTTP сервер
	server := &http.Server{
		Addr:         ":" + cfg.HTTP.Port,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.HTTP.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.HTTP.WriteTimeout) * time.Second,
	}

	// Запускаем сервер в горутине
	go func() {
		log.Printf("Starting HTTP server on :%s", cfg.HTTP.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Ожидаем сигнала для корректного завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Создаем контекст с таймаутом для корректного завершения
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(cfg.HTTP.ShutdownTimeout)*time.Second,
	)
	defer cancel()

	// Корректно останавливаем сервер
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server gracefully stopped")
}
