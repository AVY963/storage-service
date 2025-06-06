package http

import (
	"io"
	"log"
	"net/http"

	"tages/internal/models"
	"tages/internal/usecase"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	fileUsecase *usecase.Usecase
}

func NewFileHandler(fileUsecase *usecase.Usecase) *FileHandler {
	return &FileHandler{
		fileUsecase: fileUsecase,
	}
}

// UploadHandler обрабатывает загрузку зашифрованных файлов
func (h *FileHandler) UploadHandler(c *gin.Context) {
	// Получаем ID пользователя из контекста
	userID, exists := GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	// Получаем файл и метаданные из формы
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "не удалось получить файл"})
		return
	}

	encryptedKey, err := c.FormFile("key")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "не удалось получить ключ шифрования"})
		return
	}

	nonce, err := c.FormFile("nonce")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "не удалось получить вектор инициализации"})
		return
	}

	filename := c.PostForm("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "имя файла не указано"})
		return
	}

	// Читаем данные из файлов
	encryptedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка чтения файла"})
		return
	}
	defer encryptedFile.Close()

	keyFile, err := encryptedKey.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка чтения ключа шифрования"})
		return
	}
	defer keyFile.Close()

	nonceFile, err := nonce.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка чтения вектора инициализации"})
		return
	}
	defer nonceFile.Close()

	// Читаем все данные в память
	fileData, err := io.ReadAll(encryptedFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка чтения файла"})
		return
	}

	keyData, err := io.ReadAll(keyFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка чтения файла"})
		return
	}
	nonceData, err := io.ReadAll(nonceFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка чтения вектора инициализации"})
		return
	}

	// Создаем структуру с зашифрованными данными
	encryptedData := &models.EncryptedFileData{
		File:     fileData,
		Key:      keyData,
		Nonce:    nonceData,
		Filename: filename,
	}

	// Сохраняем файл с указанием пользователя
	err = h.fileUsecase.Upload(c.Request.Context(), encryptedData, userID)
	if err != nil {
		log.Printf("ERROR: Failed to upload file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка загрузки файла: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  "файл успешно загружен",
		"filename": filename,
	})
}

// DownloadHandler обрабатывает скачивание зашифрованных файлов
func (h *FileHandler) DownloadHandler(c *gin.Context) {
	// Получаем ID пользователя из контекста
	userID, exists := GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "имя файла не указано"})
		return
	}

	// Получаем зашифрованные данные для конкретного пользователя
	encryptedData, err := h.fileUsecase.Download(c.Request.Context(), filename, userID)
	if err != nil {
		log.Printf("ERROR: Failed to download file: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "файл не найден"})
		return
	}

	// Отправляем зашифрованные данные клиенту
	c.JSON(http.StatusOK, gin.H{
		"file":  encryptedData.File,
		"key":   encryptedData.Key,
		"nonce": encryptedData.Nonce,
	})
}

// ListHandler отображает список файлов
func (h *FileHandler) ListHandler(c *gin.Context) {
	// Получаем ID пользователя из контекста
	userID, exists := GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	files, err := h.fileUsecase.ListFiles(c.Request.Context(), userID)
	if err != nil {
		log.Printf("ERROR: Failed to list files: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка получения списка файлов",
		})
		return
	}

	// Формируем ответ
	type fileInfo struct {
		Name      string `json:"name"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	response := make([]fileInfo, 0, len(files))
	for _, file := range files {
		response = append(response, fileInfo{
			Name:      file.Name,
			CreatedAt: file.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: file.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"files": response,
	})
}

// DeleteHandler удаляет файл
func (h *FileHandler) DeleteHandler(c *gin.Context) {
	// Получаем ID пользователя из контекста
	userID, exists := GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "имя файла не указано",
		})
		return
	}

	err := h.fileUsecase.DeleteFile(c.Request.Context(), filename, userID)
	if err != nil {
		log.Printf("ERROR: Failed to delete file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка удаления файла",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "файл успешно удален",
	})
}
