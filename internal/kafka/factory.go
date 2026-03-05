package kafka

import (
	"go.uber.org/zap"
)

// NewProducerWithBrokers initializes a Kafka producer with specified brokers
func NewProducerWithBrokers(brokers []string, logger *zap.Logger) *Producer {
	return NewProducer(brokers, logger)
}

// NewConsumerWithBrokers initializes a Kafka consumer with specified brokers
func NewConsumerWithBrokers(brokers []string, topic string, groupID string, logger *zap.Logger) *Consumer {
	return NewConsumer(brokers, topic, groupID, logger)
}
