package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Producer struct {
	writerBugCreated  *kafka.Writer
	logger            *zap.Logger
}

func NewProducer(brokers []string, logger *zap.Logger) *Producer {

	writerBugCreated := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   EventBugCreated,
	})


	return &Producer{
		writerBugCreated:  writerBugCreated,
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


func (p *Producer) Close() error {

	if err := p.writerBugCreated.Close(); err != nil {
		return err
	}
	return nil
}