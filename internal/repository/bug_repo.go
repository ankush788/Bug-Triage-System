package repository

import (
	"context"

	"bug_triage/internal/models"
)

// BugRepository defines persistence behaviour for bugs.

type BugRepository interface {
    Create(ctx context.Context, b *models.Bug) error
    GetByID(ctx context.Context, id int64) (*models.Bug, error)
    List(ctx context.Context, limit, offset int) ([]*models.Bug, error)
    UpdateStatus(ctx context.Context, id int64, status string) error
    UpdateAnalysis(ctx context.Context, id int64, priority, category string) error
}