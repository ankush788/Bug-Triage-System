package repository

import (
	"context"
	"database/sql"

	"bug_triage/internal/models"

	"github.com/jmoiron/sqlx"
)

// PostgresUserRepo is a Postgres implementation of UserRepository.

type PostgresUserRepo struct {
    db *sqlx.DB
}

func NewPostgresUserRepo(db *sqlx.DB) *PostgresUserRepo {
    return &PostgresUserRepo{db: db}
}

func (r *PostgresUserRepo) Create(ctx context.Context, u *models.User) error {
    query := `INSERT INTO users (email, password_hash, created_at) VALUES ($1, $2, NOW()) RETURNING id, created_at`;
    return r.db.QueryRowxContext(ctx, query, u.Email, u.PasswordHash).Scan(&u.ID, &u.CreatedAt)
}

func (r *PostgresUserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    var u models.User
    err := r.db.GetContext(ctx, &u, "SELECT id, email, password_hash, created_at FROM users WHERE email=$1", email)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    return &u, err
}

func (r *PostgresUserRepo) GetByID(ctx context.Context, id int64) (*models.User, error) {
    var u models.User
    err := r.db.GetContext(ctx, &u, "SELECT id, email, password_hash, created_at FROM users WHERE id=$1", id)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    return &u, err
}