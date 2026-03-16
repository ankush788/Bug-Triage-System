package worker

import (
	"context"
	"errors"
	"time"

	"bug_triage/internal/aianalyzer"
	errortype "bug_triage/internal/error"
	"bug_triage/internal/kafka"
	"bug_triage/internal/metrics"
	"bug_triage/internal/repository"

	"go.uber.org/zap"
)

// BugAnalyzer consumes bug_created events and performs AI analysis
type BugAnalyzer struct {
	consumer   *kafka.Consumer
	bugRepo    repository.BugRepository
	producer   *kafka.Producer
	aiAnalyzer aianalyzer.Analyzer
	logger     *zap.Logger
}

func NewBugAnalyzer(
	consumer *kafka.Consumer,
	bugRepo repository.BugRepository,
	producer *kafka.Producer,
	aiAnalyzer aianalyzer.Analyzer,
	logger *zap.Logger,
) *BugAnalyzer {
	return &BugAnalyzer{
		consumer:   consumer,
		bugRepo:    bugRepo,
		producer:   producer,
		aiAnalyzer: aiAnalyzer,
		logger:     logger,
	}
}

// Start begins processing bug_created events
func (ba *BugAnalyzer) Start(ctx context.Context) error {
	//inside it because it need context 
	handler := func(ctx context.Context, message []byte) error {
		event, err := kafka.ParseBugCreatedEvent(message)
		if err != nil {
			ba.logger.Error("failed to parse bug_created event", zap.Error(err))
			return err
		}

		// Perform AI analysis
		start := time.Now()
		priority, category, err := ba.aiAnalyzer.AnalyzeBug(ba.logger, ctx, event.Title, event.Description)
		duration := time.Since(start).Seconds()
		metrics.AILatency.Observe(duration)

		if err != nil {
			ba.logger.Error("ai analysis failed", zap.Error(err), zap.Int64("bug_id", event.BugID))
			return err
		}

		ba.logger.Info("bug analysis completed",
			zap.Int64("bug_id", event.BugID),
			zap.String("priority", priority),
			zap.String("category", category),
		)

		// Update bug with analysis results
		if err := ba.bugRepo.UpdateAnalysis(ctx, event.BugID, priority, category); err != nil {
			if errors.Is(err, errortype.ErrNotFound) {
				ba.logger.Warn("bug disappeared before analysis", zap.Int64("bug_id", event.BugID))
				// not a fatal error, just swallow it
				return nil
			}
			ba.logger.Error("failed to update bug analysis", zap.Error(err))
			return err
		}

		return nil
	}

	return ba.consumer.StartConsuming(ctx, handler)
}
