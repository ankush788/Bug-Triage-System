package database

import (
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// zapGormWriter implements gorm logger.Writer interface
type zapGormWriter struct {
	log *zap.Logger
}

func (w *zapGormWriter) Printf(format string, args ...interface{}) {
	w.log.Sugar().Infof(format, args...)
}



// NewPostgresConnection initializes and returns a database connection
func NewPostgresConnection(dbURL string, log *zap.Logger) (*gorm.DB, error) {
	
	// Create a custom logger for GORM
	//  you are bridging GORM logs → Zap logger
	// if we not create this gormLogger using Zap logger and do --> db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	  // --> it will create "a gorm default logger" internally 
	  // issue -> 
	  // 1) for query logs ex:- db.Create(&user) , db.Where("email = ?", email).First(&user) =>
	  // the logs are generated in that format which is not in zap format and diffcult to find 
	gormLogger := logger.New(    
		&zapGormWriter{log: log},
		logger.Config{
			SlowThreshold:             0,   // Log all queries
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Error("failed to connect to database", zap.Error(err))
		return nil, err
	}

	// Test connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Error("failed to get underlying sql.DB", zap.Error(err))
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		log.Error("database ping failed", zap.Error(err))
		return nil, err
	}

	log.Info("database connected")
	return db, nil
}

