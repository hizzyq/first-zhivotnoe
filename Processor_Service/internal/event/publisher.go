package event

import (
	"context"
	"encoding/json"
	"fmt"
	"processor/internal/domain"
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

func (p *Publisher) PublishMediaProcessed(ctx context.Context, event domain.MediaProcessedEvent) error {
	const op = "event.publisher.PublishMediaProcessed"

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("%s: failed to marshal event: %w", op, err)
	}

	message := kafka.Message{
		Key:   []byte(event.MediaID),
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
