package event

import (
	"context"
	"encoding/json"
	"log"
	"processor/internal/domain"

	"github.com/segmentio/kafka-go"
)

type MediaProcessor interface {
	ProcessMedia(ctx context.Context, event domain.MediaUploadedEvent) error
}

type Consumer struct {
	reader    *kafka.Reader
	processor MediaProcessor
}

func NewConsumer(brokers []string, topic, groupID string, processor MediaProcessor) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	return &Consumer{
		reader:    reader,
		processor: processor,
	}
}

func (c *Consumer) Listen(ctx context.Context) error {
	const op = "event.consumer.Listen"

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping listening kafka")
			return ctx.Err()
		default:
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return nil
				}
				log.Printf("%s: error while reading message: %w", op, err)
				continue
			}
			var event domain.MediaUploadedEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Printf("%s: error while unmarhaslling message: %v", op, err)
				c.reader.CommitMessages(ctx, msg)
				continue
			}
			if err := c.processor.ProcessMedia(ctx, event); err != nil {
				log.Printf("%s: error to process media: %v", op, err)
				continue
			}
			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				log.Printf("%s: error to commit messages: %v", op, err)
			}
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
