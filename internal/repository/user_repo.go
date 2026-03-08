package repository

import (
	"context"
	"errors"

	"bug_triage/internal/models"

	"gorm.io/gorm"
)

// UserRepository defines methods for working with user persistence.

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
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *PostgresUserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &u, err
}

func (r *PostgresUserRepo) GetByID(ctx context.Context, id int64) (*models.User, error) {
	var u models.User
	err := r.db.WithContext(ctx).First(&u, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &u, err
}