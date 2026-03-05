package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"bug_triage/internal/config"
	"bug_triage/internal/database"
	"bug_triage/internal/kafka"
	"bug_triage/internal/logger"
	"bug_triage/internal/repository"
	"bug_triage/internal/worker"

	"go.uber.org/zap"
)

// Bug Analyzer Worker
// This is a standalone service that consumes bug_created events from Kafka,
// performs AI analysis, and publishes bug_analyzed events.
//
// Can be run as a separate process/container:
//   go run ./cmd/worker
//
// For distributed processing, run multiple instances with the same consumer group.

func main() {
	// Initialize logger
	logger.Init()
	defer logger.Sync()

	log := logger.Log

	// Load configuration
	cfg := config.Load()
	log.Info("bug analyzer worker starting")

	// Initialize database
	db, err := database.NewPostgresConnection(cfg.DBUrl, log)
	if err != nil {
		log.Fatal("failed to initialize database", zap.Error(err))
	}
	defer db.Close()

	// Initialize repositories
	bugRepo := repository.NewPostgresBugRepo(db)

	// Initialize Kafka producer and consumer
	kafkaProducer := kafka.NewProducerWithBrokers([]string{cfg.KafkaBroker}, log)
	defer kafkaProducer.Close()

	kafkaConsumer := kafka.NewConsumerWithBrokers(
		[]string{cfg.KafkaBroker},
		kafka.EventBugCreated,
		"bug-analyzer-worker-group",
		log,
	)
	defer kafkaConsumer.Close()

	log.Info("kafka initialized")

	// Initialize AI analyzer
	aiAnalyzer := worker.NewSimpleAIAnalyzer(log)

	// Initialize bug analyzer
	bugAnalyzer := worker.NewBugAnalyzer(
		kafkaConsumer,
		bugRepo,
		kafkaProducer,
		aiAnalyzer,
		log,
	)

	// Start consumer in a goroutine
	consumerErrors := make(chan error, 1)
	go func() {
		if err := bugAnalyzer.Start(context.Background()); err != nil {
			consumerErrors <- err
		}
	}()

	log.Info("bug analyzer worker started")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigChan:
		log.Info("shutdown signal received")
	case err := <-consumerErrors:
		log.Fatal("worker error", zap.Error(err))
	}

	log.Info("worker stopped")
}
