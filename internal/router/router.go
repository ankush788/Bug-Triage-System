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

	// Instrument requests for Prometheus
	router.Use(middleware.MetricsMiddleware())

	// Health check (no auth required) 
	// (it is private api so, public will not access it) --> no threat of ddos attack
	
	router.GET("/health", authHandler.Health)

// Auth routes (no auth required but rate limited)
    authGroup := router.Group("/auth")
    // apply rate limiter to slow down brute-force and signup abuse
    authGroup.Use(
        middleware.RateLimitMiddleware(rateLimiter, log),
    )
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
