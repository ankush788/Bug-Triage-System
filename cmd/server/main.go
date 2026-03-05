package main

import (
	"os"
	"os/signal"
	"syscall"

	"bug_triage/internal/bootstrap"
	"bug_triage/internal/config"
	"bug_triage/internal/logger"
	"bug_triage/internal/router"

	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger.Init()
	defer logger.Sync()

	log := logger.Log

	// Load configuration
	cfg := config.Load()
	log.Info("configuration loaded",
		zap.String("port", cfg.Port),
		zap.String("db_url", cfg.DBUrl),
	)

	// Initialize all dependencies
	deps, err := bootstrap.NewAppDependencies(cfg, log)
	if err != nil {
		log.Fatal("failed to initialize dependencies", zap.Error(err))
	}
	defer deps.Close()

	// Setup router --> return gin Engine pointer 
	httpRouter := router.SetupRouter(  // httprouter = gin engine
		deps.Handlers.AuthHandler,
		deps.Handlers.BugHandler,
		deps.Auth.JWTManager,
		deps.RateLimiter,
		log,
	)

	log.Info("server starting", zap.String("port", cfg.Port))

// Why not start the server normally (router.Run())?
// 1) router.Run() blocks the main thread; running it in a goroutine keeps main free.
// 2) log the actual reason of server stop ( user exist) or server start error


serverErrors := make(chan error, 1)

// Start HTTP server in background (remain main thread free to proces down things)
go func() {
	if err := httpRouter.Run(":" + cfg.Port); err != nil {
		serverErrors <- err
	}
}()

// Listen for OS shutdown signals
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

// Wait for either shutdown signal or server error (aslo log it)
select {
case <-sigChan:
	log.Info("shutdown signal received")
case err := <-serverErrors:
	log.Fatal("server error", zap.Error(err))
}

log.Info("server stopped")
}
