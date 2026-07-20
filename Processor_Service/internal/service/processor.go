package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
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

	fileReader, err := s.storage.DownloadMedia(ctx, event.MediaPath)
	if err != nil {
		return fmt.Errorf("%s: error while downloading: %w", op, err)
	}
	defer fileReader.Close()

	tempDir, err := os.MkdirTemp("", "media_process_*")
	if err != nil {
		return fmt.Errorf("%s: failed to create temp: %w", op, err)
	}
	defer os.RemoveAll(tempDir)

	filename := "processed_" + event.MediaPath
	fullPath := filepath.Join(tempDir, filename)

	// TODO: ffmpeg обработчик

	file, err := os.Open(fullPath)
	if err != nil {
		return fmt.Errorf("%s: failed to open file: %w", op, err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("%s: failed to stat file: %w", op, err)
	}
	processedSize := fileInfo.Size()

	mediaPath, err := s.storage.UploadMedia(ctx, event.Title, event.ContentType, processedSize, file)
	if err != nil {
		return fmt.Errorf("%s: failed to upload file: %w", op, err)
	}

	s.broker.PublishMediaProcessed(ctx, domain.MediaProcessedEvent{
		UserID:      event.UserID,
		MediaID:     event.MediaID,
		Title:       event.Title,
		Description: event.Description,
		MediaPath:   mediaPath,
		ContentType: event.ContentType,
	})

	return nil
}

func ProcessFile(ctx context.Context, filename, dir, contentType string) error {
	const op = "service.processor.ProcessFile"

	if contentType == "pic" {
		args := []string{
			"-y",
			"-i", filename,
			"-vf", "scale=600:-2",
			"-q:v", "80",
			dir,
		}

		cmd := exec.CommandContext(ctx, "ffmpeg", args...)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("%s: ffmpeg image error: %w, %s", op, err, string(output))
		}
		return nil
	}

	if contentType == "vid" {
		args := []string{
			"-y",
			"-i", filename,
			"-vf", "scale='min(720\\,iw)':-2",
			"-c:v", "libx264",
			"-crf", "26",
			"-preset", "fast",
			"-c:a", "aac",
			"-b:a", "128k",
			"-movflags", "+faststart",
			dir,
		}

		cmd := exec.CommandContext(ctx, "ffmpeg", args...)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("%s: ffmpeg video error: %w, log: %s", err, string(output))
		}
		return nil
	}
	return nil
}
