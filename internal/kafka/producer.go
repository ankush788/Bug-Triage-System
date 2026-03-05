package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// Producer publishes events to Kafka topics
type Producer struct {
	writers map[string]*kafka.Writer
	logger  *zap.Logger
}

func NewProducer(brokers []string, logger *zap.Logger) *Producer {
	return &Producer{
		writers: make(map[string]*kafka.Writer),
		logger:  logger,
	}
}

// getWriter returns (or creates) a writer for a topic
func (p *Producer) getWriter(topic string) *kafka.Writer {
	if w, exists := p.writers[topic]; exists {
		return w
	}

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"}, // Would be parameterized
		Topic:   topic,
	})

	p.writers[topic] = w
	return w
}

// PublishBugCreatedEvent publishes a bug created event
func (p *Producer) PublishBugCreatedEvent(ctx context.Context, event *BugCreatedEvent) error {
	data, err := event.ToJSON()
	if err != nil {
		p.logger.Error("failed to marshal bug_created event", zap.Error(err))
		return err
	}

	w := p.getWriter(EventBugCreated)
	err = w.WriteMessages(ctx, kafka.Message{
		Value: data,
	})

	if err != nil {
		p.logger.Error("failed to publish bug_created event", zap.Error(err))
		return err
	}

	p.logger.Info("published bug_created event", zap.Int64("bug_id", event.BugID))
	return nil
}

// PublishBugAnalyzedEvent publishes a bug analyzed event
func (p *Producer) PublishBugAnalyzedEvent(ctx context.Context, event *BugAnalyzedEvent) error {
	data, err := event.ToJSON()
	if err != nil {
		p.logger.Error("failed to marshal bug_analyzed event", zap.Error(err))
		return err
	}

	w := p.getWriter(EventBugAnalyzed)
	err = w.WriteMessages(ctx, kafka.Message{
		Value: data,
	})

	if err != nil {
		p.logger.Error("failed to publish bug_analyzed event", zap.Error(err))
		return err
	}

	p.logger.Info("published bug_analyzed event", zap.Int64("bug_id", event.BugID))
	return nil
}

// Close closes all writers
func (p *Producer) Close() error {
	for _, w := range p.writers {
		if err := w.Close(); err != nil {
			p.logger.Error("failed to close writer", zap.Error(err))
			return err
		}
	}
	return nil
}
