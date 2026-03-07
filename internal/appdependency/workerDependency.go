package appdependency

import (
	"bug_triage/internal/config"
	"bug_triage/internal/database"
	"bug_triage/internal/kafka"
	"bug_triage/internal/repository"
	"bug_triage/internal/worker"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// WorkerDependencies holds dependencies required for the worker service
type WorkerDependencies struct {
	DB           *sqlx.DB
	BugRepo      repository.BugRepository
	KafkaProducer *kafka.Producer
	KafkaConsumer *kafka.Consumer
	AIAnalyzer   *worker.SimpleAIAnalyzer
	Logger       *zap.Logger
}

// NewWorkerDependencies initializes dependencies required for the worker
func NewWorkerDependencies(cfg *config.Config, log *zap.Logger) (*WorkerDependencies, error) {
	// Initialize database
	db, err := database.NewPostgresConnection(cfg.DBUrl, log)
	if err != nil {
		return nil, err
	}

	// Initialize repositories
	bugRepo := repository.NewPostgresBugRepo(db)

	// Initialize Kafka producer and consumer
	kafkaProducer := kafka.NewProducerWithBrokers([]string{cfg.KafkaBroker}, log)
	kafkaConsumer := kafka.NewConsumerWithBrokers(
		[]string{cfg.KafkaBroker},
		kafka.EventBugCreated,
		"bug-analyzer-worker-group",
		log,
	)

	// Initialize AI analyzer
	aiAnalyzer := worker.NewSimpleAIAnalyzer(log)

	return &WorkerDependencies{
		DB:            db,
		BugRepo:       bugRepo,
		KafkaProducer: kafkaProducer,
		KafkaConsumer: kafkaConsumer,
		AIAnalyzer:    aiAnalyzer,
		Logger:        log,
	}, nil
}

// Close closes all closeable worker dependencies
func (d *WorkerDependencies) Close() error {
	if d.DB != nil {
		d.DB.Close()
	}
	if d.KafkaProducer != nil {
		d.KafkaProducer.Close()
	}
	if d.KafkaConsumer != nil {
		d.KafkaConsumer.Close()
	}
	return nil
}
