package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"bug_triage/internal/appdependency"
	"bug_triage/internal/config"
	"bug_triage/internal/logger"
	"bug_triage/internal/migration"
	"bug_triage/internal/worker"

	"go.uber.org/zap"
)

// Bug Analyzer Worker
// This is a standalone service that consumes bug_created events from Kafka,
// performs AI analysis, and publishes bug_analyzed events.

// Reason to run on seperate server
//Independent scaling: API server handles user requests while workers process Kafka jobs. so, both can scale independently and seperatly based on their requirment
//Fault isolation: If the worker crashes or heavy processing occurs, the API server continues serving requests without downtime.
//Industry microservice pattern: Background jobs (event consumers) are usually deployed as separate services from the main API is industry trend



func main() {
	// Initialize logger
	logger.Init()
	defer logger.Sync()

	log := logger.Log

	// Load configuration
	cfg := config.Load()
	log.Info("bug analyzer worker starting")

	// run migrations before starting worker (worker also interacts with database)
	if err := migration.Run(cfg.DBUrl, log); err != nil {
		log.Fatal("database migration failed", zap.Error(err))
	}

	// Initialize all worker dependencies
	deps, err := appdependency.NewWorkerDependencies(cfg, log)
	if err != nil {
		log.Fatal("failed to initialize worker dependencies", zap.Error(err))
	}
	defer deps.Close()

	// Initialize bug analyzer
	bugAnalyzer := worker.NewBugAnalyzer(
		deps.KafkaConsumer,
		deps.BugRepo,
		deps.KafkaProducer,
		deps.AIAnalyzer,
		deps.Logger,
	)

	log.Info("bug analyzer worker started")

// Why not start the worker normally (bugAnalyzer.Start())?
// 1) bugAnalyzer.Start() blocks the main thread; running it in a goroutine keeps main free.
// 2) log the actual reason of worker stop (user exit) or worker start error

workerErrors := make(chan error, 1)

// Start worker in background (remain main thread free to process shutdown)  --> which use kakfa
go func() {
	if err := bugAnalyzer.Start(context.Background()); err != nil {
		workerErrors <- err
	}
}()

// Listen for OS shutdown signals
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

// Wait for either shutdown signal or worker error (also log it)
select {
case <-sigChan:
	log.Info("shutdown signal received")
case err := <-workerErrors:
	log.Fatal("worker error", zap.Error(err))
}

log.Info("worker stopped")
}
