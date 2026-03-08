package database

import (
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewPostgresConnection initializes and returns a database connection
func NewPostgresConnection(dbURL string, log *zap.Logger) (*gorm.DB, error) {

	// Use GORM's default logger
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // show SQL queries
	})
	if err != nil {
		log.Error("failed to connect to database", zap.Error(err))
		return nil, err
	}

	// Get underlying sql.DB to verify connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Error("failed to get underlying sql.DB", zap.Error(err))
		return nil, err
	}

	// Ping database
	if err := sqlDB.Ping(); err != nil {
		log.Error("database ping failed", zap.Error(err))
		return nil, err
	}

	log.Info("database connected")

	return db, nil
}