package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// NewPostgresConnection initializes and returns a database connection
func NewPostgresConnection(dbURL string, log *zap.Logger) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		log.Error("failed to connect to database", zap.Error(err))
		return nil, err
	}

	// Test connection
	if err := db.Ping(); err != nil {
		log.Error("database ping failed", zap.Error(err))
		return nil, err
	}

	log.Info("database connected")
	return db, nil
}
