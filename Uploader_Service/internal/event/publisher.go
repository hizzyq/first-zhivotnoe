package publisher

import (
	"chooki/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type Publisher struct {
	writer *kafka.Writer
}

func NewPublisher(brokers []string, topic string) (*Publisher, error) {
	const op = "event.publisher.NewPublisher"

	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		WriteTimeout: 10 * time.Second,
		Async:        false,
	}

	return &Publisher{writer: writer}, nil
}

func (p *Publisher) PublishMediaUploaded(ctx context.Context, event domain.MediaUploadedEvent) error {
	const op = "event.publisher.PublishMediaUploaded"

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("%s: failed to marshal event: %w", op, err)
	}

	message := kafka.Message{
		Key:   []byte(event.FileName),
		Value: payload,
	}

	err = p.writer.WriteMessages(ctx, message)
	if err != nil {
		return fmt.Errorf("%s: failed to write message: %w", op, err)
	}
	return nil
}

func (p *Publisher) Close() error {
	return p.writer.Close()
}
