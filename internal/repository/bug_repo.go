package repository

import (
	"context"
	"errors"

	errortype "bug_triage/internal/err"
	"bug_triage/internal/models"

	"gorm.io/gorm"
)

// BugRepository defines  operations/behaviour for bugs.

// we create different repo structure/files for different domain (bug , user) because their storing data or
//acessing db pattern may  be differnet domain to domain

// BugRepository defines  operations/behaviour for bugs.
//
// we create different repo structure/files for different domain (bug , user) because their storing data or
//acessing db pattern may  be differnet domain to domain
//
// methods that look up a single entity return ErrNotFound when the row does not exist.
type BugRepository interface {
    Create(ctx context.Context, b *models.Bug) error
    GetByID(ctx context.Context, id int64) (*models.Bug, error)
    List(ctx context.Context, limit, offset int) ([]*models.Bug, error)
    UpdateStatus(ctx context.Context, id int64, status string) error
    UpdateAnalysis(ctx context.Context, id int64, priority, category string) error
}

// PostgresBugRepo is a Postgres implementation of BugRepository.
// can create different more struct for different implmentation like Mongo DB
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
		return nil, errortype.ErrNotFound
	}
	return &b, err
}

func (r *PostgresBugRepo) List(ctx context.Context, limit, offset int) ([]*models.Bug, error) {
	var bugs []*models.Bug
	err := r.db.WithContext(ctx).Order("created_at DESC").Limit(limit).Offset(offset).Find(&bugs).Error
	return bugs, err
}

func (r *PostgresBugRepo) UpdateStatus(ctx context.Context, id int64, status string) error {
	res := r.db.WithContext(ctx).Model(&models.Bug{}).Where("id = ?", id).Update("status", status)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errortype.ErrNotFound
	}
	return nil
}

func (r *PostgresBugRepo) UpdateAnalysis(ctx context.Context, id int64, priority, category string) error {
	res := r.db.WithContext(ctx).Model(&models.Bug{}).Where("id = ?", id).Updates(map[string]interface{}{
		"priority": priority,
		"category": category,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errortype.ErrNotFound
	}
	return nil
}