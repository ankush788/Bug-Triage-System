package migration

import (
	"go.uber.org/zap"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Run applies all up migrations located in the migrations directory using
// the provided database URL. It logs progress and returns any error
// encountered. If there are no new changes, the function returns nil.
func Run(dbURL string, log *zap.Logger) error {
    log.Info("running database migrations", zap.String("db_url", dbURL))

    m, err := migrate.New("file://migrations", dbURL)
    if err != nil {
        log.Error("failed to create migrate instance", zap.Error(err))
        return err
    }

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Error("migration up failed", zap.Error(err))
        return err
    }

    log.Info("database migrations applied")
    return nil
}
