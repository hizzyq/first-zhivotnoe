package handler

import (
	"chooki/internal/domain"
	"context"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MediaUploader interface {
	UploadMedia(ctx context.Context, objectName, contentType string, size int64, reader io.Reader) (string, error)
}

type EventPublisher interface {
	PublishMediaUploaded(ctx context.Context, event domain.MediaUploadedEvent) error
}

type Handler struct {
	s3Storage MediaUploader
	broker    EventPublisher
}

func NewHandler(broker EventPublisher, s3Storage MediaUploader) *Handler {
	return &Handler{
		s3Storage: s3Storage,
		broker:    broker,
	}
}

func (h *Handler) UploadMedia(c *gin.Context) {
	// 1. Достаем файл из запроса
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get file from request"})
		return
	}

	// 2. Достаем текстовые поля
	title := c.PostForm("title")
	description := c.PostForm("description")

	// 3. Получаем ID пользователя
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		userID = uuid.New().String()
	}

	// 4. Открываем поток файла
	fileStream, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer fileStream.Close()

	// 5. Загружаем в MinIO, используя ОРИГИНАЛЬНОЕ имя файла вместо ключа
	contentType := file.Header.Get("Content-Type")
	_, err = h.s3Storage.UploadMedia(c.Request.Context(), file.Filename, contentType, file.Size, fileStream)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload file to storage"})
		return
	}

	// 6. Отправляем событие в Kafka
	err = h.broker.PublishMediaUploaded(c.Request.Context(), domain.MediaUploadedEvent{
		UserID:      userID,
		Title:       title,
		Description: description,
		ImagePath:   file.Filename,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to publish event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "file processed successfully",
		"file_name": file.Filename,
	})
}
