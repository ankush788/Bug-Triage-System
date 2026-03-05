package repository

import (
	"context"

	"bug_triage/internal/models"
)

// UserRepository defines methods for working with user persistence.

type UserRepository interface {
    Create(ctx context.Context, u *models.User) error
    GetByEmail(ctx context.Context, email string) (*models.User, error)
    GetByID(ctx context.Context, id int64) (*models.User, error)
}