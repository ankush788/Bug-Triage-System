package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Producer struct {
	writerBugCreated  *kafka.Writer
	writerBugAnalyzed *kafka.Writer
	logger            *zap.Logger
}

func NewProducer(brokers []string, logger *zap.Logger) *Producer {

	writerBugCreated := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   EventBugCreated,
	})

	writerBugAnalyzed := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   EventBugAnalyzed,
	})

	return &Producer{
		writerBugCreated:  writerBugCreated,
		writerBugAnalyzed: writerBugAnalyzed,
		logger:            logger,
	}
}

func (p *Producer) PublishBugCreatedEvent(ctx context.Context, event *BugCreatedEvent) error {
	data, err := event.ToJSON()
	if err != nil {
		p.logger.Error("failed to marshal bug_created event", zap.Error(err))
		return err
	}

	err = p.writerBugCreated.WriteMessages(ctx, kafka.Message{
		Value: data,
	})

	if err != nil {
		p.logger.Error("failed to publish bug_created event", zap.Error(err))
		return err
	}

	p.logger.Info("published bug_created event", zap.Int64("bug_id", event.BugID))
	return nil
}

func (p *Producer) PublishBugAnalyzedEvent(ctx context.Context, event *BugAnalyzedEvent) error {
	data, err := event.ToJSON()
	if err != nil {
		p.logger.Error("failed to marshal bug_analyzed event", zap.Error(err))
		return err
	}

	err = p.writerBugAnalyzed.WriteMessages(ctx, kafka.Message{
		Value: data,
	})

	if err != nil {
		p.logger.Error("failed to publish bug_analyzed event", zap.Error(err))
		return err
	}

	p.logger.Info("published bug_analyzed event", zap.Int64("bug_id", event.BugID))
	return nil
}

func (p *Producer) Close() error {

	if err := p.writerBugCreated.Close(); err != nil {
		return err
	}

	if err := p.writerBugAnalyzed.Close(); err != nil {
		return err
	}

	return nil
}