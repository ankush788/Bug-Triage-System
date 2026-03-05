package router

import (
	"bug_triage/internal/auth"
	"bug_triage/internal/handler"
	"bug_triage/internal/middleware"
	"bug_triage/internal/pkg"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SetupRouter initializes and returns a configured Gin router
func SetupRouter(
	authHandler *handler.AuthHandler,
	bugHandler *handler.BugHandler,
	jwtManager *auth.JWTManager,
	rateLimiter *pkg.RateLimiter,
	log *zap.Logger,
) *gin.Engine {
	router := gin.Default()

	// Health check (no auth required)
	router.GET("/health", authHandler.Health)

	// Auth routes (no auth required)
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
	}

	// Bug routes (auth required, rate limited)
	bugGroup := router.Group("/bugs")
	bugGroup.Use(
		middleware.JWTMiddleware(jwtManager, log),
		middleware.RateLimitMiddleware(rateLimiter, log),
	)
	{
		bugGroup.POST("", bugHandler.CreateBug)
		bugGroup.GET("/:id", bugHandler.GetBug)
		bugGroup.GET("", bugHandler.ListBugs)
		bugGroup.PATCH("/:id/status", bugHandler.UpdateBugStatus)
	}

	log.Info("routes registered")
	return router
}
