package migration

import (
	"errors"

	"go.uber.org/zap"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Run executes database migrations using the provided database URL
func Run(dbURL string, log *zap.Logger) error {
	
	m, err := migrate.New("file://migrations", dbURL) // create migration instance using migrations folder and db connection
	if err != nil {
		log.Error("failed to create migrate instance", zap.Error(err))
		return err
	}


	// schema table ke version ko check karta hai and uske aage ke version ko  run karta hai 
	// (if version "1" present in table then start running migration file from "2" and so on)
	if err := m.Up(); err != nil && err != migrate.ErrNoChange { // run all pending up migrations

	// if err := m.Down(); err != nil && err != migrate.ErrNoChange { for downgrade we have to use this 

		var de migrate.ErrDirty

		if errors.As(err, &de) { // check if error occurred because database is in dirty state
			ver := uint(de.Version) // get version where migration stopped
            
			log.Warn("database in dirty state, forcing version", zap.Uint("version", ver)) 

			if ferr := m.Force(int(ver)); ferr != nil { // force migration version to clear dirty state (roll back to previous clear state)
				log.Error("failed to force version", zap.Error(ferr))
				return ferr   // exist if it not come to clean state 
			}

			if rer := m.Up(); rer != nil && rer != migrate.ErrNoChange { // retry running migrations
				log.Error("migration up failed after force", zap.Error(rer))
				return rer
			}

			log.Info("database migrations applied after clearing dirty flag") // log successful recovery
			return nil
		}

		log.Error("migration up failed", zap.Error(err)) // log general migration error
		return err
	}

	log.Info("database migrations applied") // log successful migration execution
	return nil
}