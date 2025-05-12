package http

import (
	"io"
	"log"
	"net/http"

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

// UploadHandler обрабатывает загрузку файлов
func (h *FileHandler) UploadHandler(c *gin.Context) {
	// Получаем файл из формы
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "не удалось получить файл",
		})
		return
	}
	defer file.Close()

	// Получаем имя файла
	filename := header.Filename

	// Получаем содержимое файла
	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка чтения файла",
		})
		return
	}

	// Сохраняем файл
	err = h.fileUsecase.Upload(c.Request.Context(), filename, data)
	if err != nil {
		log.Printf("ERROR: Failed to upload file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка загрузки файла: " + err.Error(),
		})
		return
	}

	// Формируем успешный ответ
	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  "файл успешно загружен",
		"filename": filename,
	})
}

// DownloadHandler обрабатывает скачивание файлов
func (h *FileHandler) DownloadHandler(c *gin.Context) {
	filename := c.Param("filename")

	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "имя файла не указано",
		})
		return
	}

	// Открываем файл для чтения
	fileReader, err := h.fileUsecase.Download(filename)
	if err != nil {
		log.Printf("ERROR: Failed to download file: %v", err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "файл не найден",
		})
		return
	}
	defer fileReader.Close()

	// Устанавливаем заголовки для скачивания
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")

	// Копируем содержимое файла в ответ
	_, err = io.Copy(c.Writer, fileReader)
	if err != nil {
		log.Printf("ERROR: Failed to send file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка отправки файла",
		})
		return
	}
}

// ListHandler отображает список файлов
func (h *FileHandler) ListHandler(c *gin.Context) {
	files, err := h.fileUsecase.ListFiles(c.Request.Context())
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
	filename := c.Param("filename")

	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "имя файла не указано",
		})
		return
	}

	err := h.fileUsecase.DeleteFile(c.Request.Context(), filename)
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
