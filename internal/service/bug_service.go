package service

import (
	"context"
	"errors"

	"bug_triage/internal/cache"
	"bug_triage/internal/dto"
	"bug_triage/internal/kafka"
	"bug_triage/internal/models"
	"bug_triage/internal/repository"

	"go.uber.org/zap"
)

// BugService handles bug-related business logic
type BugService struct {
	bugRepo  repository.BugRepository
	producer *kafka.Producer
	logger   *zap.Logger
	bugCache *cache.BugCache
}

func NewBugService(
	bugRepo repository.BugRepository,
	producer *kafka.Producer,
	logger *zap.Logger,
	bugCache *cache.BugCache,
) *BugService {
	return &BugService{
		bugRepo:  bugRepo,
		producer: producer,
		logger:   logger,
		bugCache: bugCache,
	}
}

// didn't use redis in ListBUg :-
// frequent write/update option make redis overhead beside of reducing latancy

// CreateBug creates a new bug report and publishes event
func (s *BugService) CreateBug(ctx context.Context, req *dto.CreateBugRequest, reporterID int64) (*models.Bug, error) {
	
	bug := &models.Bug{
		Title:       req.Title,
		Description: req.Description,
		ReporterID:  reporterID,
		Status:      "OPEN",
		Priority:    "UNKNOWN",
		Category:    "UNCLASSIFIED",
	}

	// Save to database
	if err := s.bugRepo.Create(ctx, bug); err != nil {
		s.logger.Error("failed to create bug", zap.Error(err))
		return nil, err
	}

	s.logger.Info("bug created", zap.Int64("bug_id", bug.ID))

	// Cache the bug
	if err := s.bugCache.Set(ctx, bug.ID, bug); err != nil {
		s.logger.Warn("failed to cache bug", zap.Error(err))
	}

	// Publish event for processing
	event := &kafka.BugCreatedEvent{
		BugID:       bug.ID,
		Title:       bug.Title,
		Description: bug.Description,
		ReporterID:  bug.ReporterID,
	}

	if err := s.producer.PublishBugCreatedEvent(ctx, event); err != nil {
		s.logger.Error("failed to publish bug_created event", zap.Error(err))
		// Don't fail the request - bug was created successfully, just async processing failed
	}

	return bug, nil
}

// GetBug retrieves a bug by ID
func (s *BugService) GetBug(ctx context.Context, bugID int64) (*models.Bug, error) {
	// Try to get from cache
	bug, err := s.bugCache.Get(ctx, bugID)
	if err == nil {
		s.logger.Debug("bug retrieved from cache", zap.Int64("bug_id", bugID))
		return bug, nil
	}

	// Get from database
	bug, err = s.bugRepo.GetByID(ctx, bugID)
	if err != nil {
		s.logger.Error("failed to get bug", zap.Error(err))
		return nil, err
	}
	if bug == nil {
		return nil, errors.New("bug not found")
	}

	// Cache the result
	if err := s.bugCache.Set(ctx, bug.ID, bug); err != nil {
		s.logger.Warn("failed to cache bug", zap.Error(err))
	}

	return bug, nil
}

// ListBugs retrieves a paginated list of bugs
func (s *BugService) ListBugs(ctx context.Context, limit, offset int) ([]*models.Bug, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	bugs, err := s.bugRepo.List(ctx, limit, offset)
	if err != nil {
		s.logger.Error("failed to list bugs", zap.Error(err))
		return nil, err
	}

	return bugs, nil
}

// UpdateBugStatus updates the status of a bug
func (s *BugService) UpdateBugStatus(ctx context.Context, bugID int64, status string) error {
	// Validate status
	validStatuses := map[string]bool{
		"OPEN":       true,
		"IN_PROGRESS": true,
		"RESOLVED":   true,
		"CLOSED":     true,
	}
	if !validStatuses[status] {
		return errors.New("invalid status")
	}

	if err := s.bugRepo.UpdateStatus(ctx, bugID, status); err != nil {
		s.logger.Error("failed to update bug status", zap.Error(err))
		return err
	}

	// Invalidate cache
	if err := s.bugCache.Delete(ctx, bugID); err != nil {
		s.logger.Warn("failed to invalidate cache", zap.Error(err))
	}

	s.logger.Info("bug status updated",
		zap.Int64("bug_id", bugID),
		zap.String("status", status),
	)

	return nil
}
