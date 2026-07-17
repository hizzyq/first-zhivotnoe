package main

import (
	"chooki/internal/config" // импортируем наш новый конфиг
	handler "chooki/internal/delivery/http"
	publisher "chooki/internal/event"
	minios3 "chooki/internal/storage/minio"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Инициализируем конфигурацию
	cfg := config.MustLoad()

	// 2. Инициализируем MinIO, беря параметры из конфига
	s3Storage, err := minios3.New(
		cfg.S3.Endpoint,
		cfg.S3.AccessKey,
		cfg.S3.SecretKey,
		cfg.S3.BucketName,
		cfg.S3.UseSSL == "true" || cfg.S3.UseSSL == "1", // само приведется в bool
	)
	if err != nil {
		log.Fatalf("Unable to init MinIO storage: %v", err)
	}

	// 3. Инициализируем Кафку из конфига
	// Продюсер принимает слайс строк, поэтому оборачиваем адрес в []string
	kafkaBroker, err := publisher.NewPublisher([]string{cfg.Broker.Address}, cfg.Broker.QueueName)
	if err != nil {
		log.Fatalf("Failed to init new publisher: %v", err)
	}
	defer kafkaBroker.Close()

	// 4. Настраиваем Gin
	router := gin.Default()

	// Ограничение памяти на основе конфига
	router.MaxMultipartMemory = int64(cfg.Limits.MaxSizeMB) << 20

	h := handler.NewHandler(kafkaBroker, s3Storage)
	router.POST("/upload", h.UploadMedia)

	log.Printf("server starting on port %s in %s mode...", cfg.Server.Port, cfg.Env)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("мама... %v", err)
	}
}
