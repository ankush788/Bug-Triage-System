package repository

import (
	"context"
	"errors"

	"bug_triage/internal/models"

	"gorm.io/gorm"
)

// PostgresBugRepo is a Postgres implementation of BugRepository.
type PostgresBugRepo struct {
	db *gorm.DB
}

func NewPostgresBugRepo(db *gorm.DB) *PostgresBugRepo {
	return &PostgresBugRepo{db: db}
}

func (r *PostgresBugRepo) Create(ctx context.Context, b *models.Bug) error {
	return r.db.WithContext(ctx).Create(b).Error
}

func (r *PostgresBugRepo) GetByID(ctx context.Context, id int64) (*models.Bug, error) {
	var b models.Bug
	err := r.db.WithContext(ctx).First(&b, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &b, err
}

func (r *PostgresBugRepo) List(ctx context.Context, limit, offset int) ([]*models.Bug, error) {
	var bugs []*models.Bug
	err := r.db.WithContext(ctx).Order("created_at DESC").Limit(limit).Offset(offset).Find(&bugs).Error
	return bugs, err
}

func (r *PostgresBugRepo) UpdateStatus(ctx context.Context, id int64, status string) error {
	return r.db.WithContext(ctx).Model(&models.Bug{}).Where("id = ?", id).Update("status", status).Error
}

func (r *PostgresBugRepo) UpdateAnalysis(ctx context.Context, id int64, priority, category string) error {
	return r.db.WithContext(ctx).Model(&models.Bug{}).Where("id = ?", id).Updates(map[string]interface{}{
		"priority": priority,
		"category": category,
	}).Error
}
