package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// MessageHandler is a function that processes a Kafka message
type MessageHandler func(ctx context.Context, message []byte) error

// Consumer consumes messages from a Kafka topic
type Consumer struct {
	reader *kafka.Reader
	logger *zap.Logger
	topic  string
}

func NewConsumer(brokers []string, topic string, groupID string, logger *zap.Logger) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		StartOffset:    kafka.LastOffset,
		CommitInterval: time.Second,
		MaxBytes:       10e6, // 10MB
	})

	return &Consumer{
		reader: reader,
		logger: logger,
		topic:  topic,
	}
}

// StartConsuming starts consuming messages from the topic
// Blocks until context is cancelled or an error occurs
func (c *Consumer) StartConsuming(ctx context.Context, handler MessageHandler) error {
	c.logger.Info("starting consumer", zap.String("topic", c.topic))

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("consumer context cancelled")
			return ctx.Err()
		default:
		}

		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			c.logger.Error("failed to fetch message", zap.Error(err))
			return err
		}

		c.logger.Debug("received message", zap.String("topic", c.topic), zap.Int64("offset", msg.Offset))

		if err := handler(ctx, msg.Value); err != nil {
			c.logger.Error("handler error", zap.Error(err), zap.String("topic", c.topic))
			// Continue consuming even on handler error
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			c.logger.Error("failed to commit message", zap.Error(err))
		}
	}
}

// Close closes the consumer reader
func (c *Consumer) Close() error {
	if err := c.reader.Close(); err != nil {
		c.logger.Error("failed to close reader", zap.Error(err))
		return err
	}
	return nil
}
