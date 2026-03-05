package repository

import (
	"context"
	"database/sql"

	"bug_triage/internal/models"

	"github.com/jmoiron/sqlx"
)

// PostgresBugRepo is a Postgres implementation of BugRepository.
type PostgresBugRepo struct {
	db *sqlx.DB
}

func NewPostgresBugRepo(db *sqlx.DB) *PostgresBugRepo {
	return &PostgresBugRepo{db: db}
}

func (r *PostgresBugRepo) Create(ctx context.Context, b *models.Bug) error {
	query := `INSERT INTO bugs (title, description, reporter_id, status, priority, category, created_at) 
	         VALUES ($1, $2, $3, $4, $5, $6, NOW()) 
	         RETURNING id, created_at`

	return r.db.QueryRowxContext(ctx, query,
		b.Title,
		b.Description,
		b.ReporterID,
		b.Status,
		b.Priority,
		b.Category,
	).Scan(&b.ID, &b.CreatedAt)
}

func (r *PostgresBugRepo) GetByID(ctx context.Context, id int64) (*models.Bug, error) {
	var b models.Bug
	err := r.db.GetContext(ctx, &b,
		`SELECT id, title, description, reporter_id, status, priority, category, created_at 
		 FROM bugs WHERE id=$1`,
		id,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &b, err
}

func (r *PostgresBugRepo) List(ctx context.Context, limit, offset int) ([]*models.Bug, error) {
	var bugs []*models.Bug
	err := r.db.SelectContext(ctx, &bugs,
		`SELECT id, title, description, reporter_id, status, priority, category, created_at 
		 FROM bugs ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err == sql.ErrNoRows {
		return []*models.Bug{}, nil
	}
	return bugs, err
}

func (r *PostgresBugRepo) UpdateStatus(ctx context.Context, id int64, status string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE bugs SET status=$1 WHERE id=$2`,
		status, id,
	)
	return err
}

func (r *PostgresBugRepo) UpdateAnalysis(ctx context.Context, id int64, priority, category string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE bugs SET priority=$1, category=$2 WHERE id=$3`,
		priority, category, id,
	)
	return err
}
