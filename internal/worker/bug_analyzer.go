package worker

import (
	"context"

	"bug_triage/internal/kafka"
	"bug_triage/internal/repository"

	"go.uber.org/zap"
)

// BugAnalyzer consumes bug_created events and performs AI analysis
type BugAnalyzer struct {
	consumer   *kafka.Consumer
	bugRepo    repository.BugRepository
	producer   *kafka.Producer
	aiAnalyzer AIAnalyzer
	logger     *zap.Logger
}

// AIAnalyzer simulates AI bug classification
type AIAnalyzer interface {
	AnalyzeBug(ctx context.Context, title, description string) (priority, category string, err error)
}

func NewBugAnalyzer(
	consumer *kafka.Consumer,
	bugRepo repository.BugRepository,
	producer *kafka.Producer,
	aiAnalyzer AIAnalyzer,
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
		priority, category, err := ba.aiAnalyzer.AnalyzeBug(ctx, event.Title, event.Description)
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
			ba.logger.Error("failed to update bug analysis", zap.Error(err))
			return err
		}

		// Publish bug_analyzed event
		analyzedEvent := &kafka.BugAnalyzedEvent{
			BugID:    event.BugID,
			Priority: priority,
			Category: category,
		}

		if err := ba.producer.PublishBugAnalyzedEvent(ctx, analyzedEvent); err != nil {
			ba.logger.Error("failed to publish bug_analyzed event", zap.Error(err))
			// Don't return error here - analysis was successful, just notification failed
		}

		return nil
	}

	return ba.consumer.StartConsuming(ctx, handler)
}
