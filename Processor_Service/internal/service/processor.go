package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"processor/internal/domain"
)

type MediaStorage interface {
	UploadMedia(ctx context.Context, objectName, contentType string, size int64, reader io.Reader) (string, error)
	DownloadMedia(ctx context.Context, objectName string) (io.ReadCloser, error)
}

type EventPublisher interface {
	PublishMediaProcessed(ctx context.Context, event domain.MediaProcessedEvent) error
}

type ProcessorService struct {
	storage MediaStorage
	broker  EventPublisher
}

func NewProcessorService(storage MediaStorage, broker EventPublisher) *ProcessorService {
	return &ProcessorService{
		storage: storage,
		broker:  broker,
	}
}

func (s *ProcessorService) ProcessMedia(ctx context.Context, event domain.MediaUploadedEvent) error {
	const op = "service.processor.ProcessMedia"

	// Создаем временную папку
	tempDir, err := os.MkdirTemp("", "media_process_*")
	if err != nil {
		return fmt.Errorf("%s: failed to create temp dir: %w", op, err)
	}
	defer os.RemoveAll(tempDir)

	// Формируем пути к исходному и обработанному файлу на диске
	inputPath := filepath.Join(tempDir, "input_"+event.MediaPath)
	outputPath := filepath.Join(tempDir, "processed_"+event.MediaPath)

	// Скачиваем файл из MinIO
	fileReader, err := s.storage.DownloadMedia(ctx, event.MediaPath)
	if err != nil {
		return fmt.Errorf("%s: error while downloading: %w", op, err)
	}
	defer fileReader.Close()

	// Сохраняем скачанный поток во временный входной файл
	inputFile, err := os.Create(inputPath)
	if err != nil {
		return fmt.Errorf("%s: failed to create input file: %w", op, err)
	}
	if _, err := io.Copy(inputFile, fileReader); err != nil {
		inputFile.Close()
		return fmt.Errorf("%s: failed to save downloaded stream to disk: %w", op, err)
	}
	inputFile.Close()

	// Запускаем FFmpeg обработку
	if err := ProcessFile(ctx, inputPath, outputPath, event.ContentType); err != nil {
		return fmt.Errorf("%s: ffmpeg processing failed: %w", op, err)
	}

	// Открываем обработанный файл для загрузки обратно в MinIO
	processedFile, err := os.Open(outputPath)
	if err != nil {
		return fmt.Errorf("%s: failed to open processed file: %w", op, err)
	}
	defer processedFile.Close()

	fileInfo, err := processedFile.Stat()
	if err != nil {
		return fmt.Errorf("%s: failed to stat processed file: %w", op, err)
	}

	// Загружаем результат в MinIO
	processedName := "processed_" + event.MediaPath
	mediaPath, err := s.storage.UploadMedia(ctx, processedName, event.ContentType, fileInfo.Size(), processedFile)
	if err != nil {
		return fmt.Errorf("%s: failed to upload file: %w", op, err)
	}

	// Публикуем событие
	err = s.broker.PublishMediaProcessed(ctx, domain.MediaProcessedEvent{
		UserID:      event.UserID,
		MediaID:     event.MediaID,
		Title:       event.Title,
		Description: event.Description,
		MediaPath:   mediaPath,
		ContentType: event.ContentType,
	})
	if err != nil {
		return fmt.Errorf("%s: failed to publish event: %w", op, err)
	}

	return nil
}

func ProcessFile(ctx context.Context, inputPath, outputPath, contentType string) error {
	const op = "service.processor.ProcessFile"

	isImage := contentType == "pic" || strings.HasPrefix(contentType, "image/")
	isVideo := contentType == "vid" || strings.HasPrefix(contentType, "video/")

	var args []string

	if isImage {
		args = []string{
			"-y",
			"-i", inputPath,
			"-vf", "scale=600:-2",
			"-q:v", "80",
			outputPath,
		}
	} else if isVideo {
		args = []string{
			"-y",
			"-i", inputPath,
			"-vf", "scale='min(720\\,iw)':-2",
			"-c:v", "libx264",
			"-crf", "26",
			"-preset", "fast",
			"-c:a", "aac",
			"-b:a", "128k",
			"-movflags", "+faststart",
			outputPath,
		}
	} else {
		return fmt.Errorf("%s: unsupported content type: %s", op, contentType)
	}

	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%s: ffmpeg error: %w, log: %s", op, err, string(output))
	}

	return nil
}
