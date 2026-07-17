package minios3

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Storage struct {
	client     *minio.Client
	bucketName string
}

func New(endpoint, accessKey, secretKey, bucketName string, useSSL bool) (*Storage, error) {
	const op = "storage.minio.new"
	// Инициализация клиента
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: failed to init minio client: %w", op, err)
	}

	ctx := context.Background()

	// Проверка на существование бакета
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to check bucket existance: %w", op, err)
	}

	// Создание бакета
	if !exists {
		err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("%s: failed to create a bucket: %w", op, err)
		}
	}

	return &Storage{
		client:     minioClient,
		bucketName: bucketName,
	}, nil
}

func (s *Storage) UploadMedia(ctx context.Context, objectName, contentType string, size int64, reader io.Reader) (string, error) {
	const op = "storage.minio.uploadfile"

	info, err := s.client.PutObject(ctx, s.bucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("%s: failed to put object: %w", op, err)
	}

	// Возвращается ключ(путь) файла
	return info.Key, nil
}
