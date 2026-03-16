package repository

import (
	"context"
	"errors"
	"time"

	errortype "bug_triage/internal/error"
	"bug_triage/internal/metrics"
	"bug_triage/internal/models"

	"gorm.io/gorm"
)

// UserRepository defines methods for working with user persistence.
//
// like the bug repository, lookups return ErrNotFound when there is no matching row.
// we create different repo structure/files for different domain (bug , user) because their storing data or
//acessing db pattern may  be differnet domain to domain

type UserRepository interface {
    Create(ctx context.Context, u *models.User) error
    GetByEmail(ctx context.Context, email string) (*models.User, error)
    GetByID(ctx context.Context, id int64) (*models.User, error)
}

// PostgresUserRepo is a Postgres implementation of UserRepository.
type PostgresUserRepo struct {
	db *gorm.DB
}

func NewPostgresUserRepo(db *gorm.DB) *PostgresUserRepo {
	return &PostgresUserRepo{db: db}
}

func (r *PostgresUserRepo) Create(ctx context.Context, u *models.User) error {
	start := time.Now()
	err := r.db.WithContext(ctx).Create(u).Error
	duration := time.Since(start).Seconds()

	metrics.DBQueryDuration.WithLabelValues("insert", "users").Observe(duration)

	return err
}

func (r *PostgresUserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	start := time.Now()
	var u models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	duration := time.Since(start).Seconds()

	metrics.DBQueryDuration.WithLabelValues("select", "users").Observe(duration)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errortype.ErrNotFound
	}
	return &u, err
}

func (r *PostgresUserRepo) GetByID(ctx context.Context, id int64) (*models.User, error) {
	start := time.Now()
	var u models.User
	err := r.db.WithContext(ctx).First(&u, id).Error
	duration := time.Since(start).Seconds()

	metrics.DBQueryDuration.WithLabelValues("select", "users").Observe(duration)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errortype.ErrNotFound
	}
	return &u, err
}