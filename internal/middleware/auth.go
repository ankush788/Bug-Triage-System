package middleware

import (
	"net/http"
	"strings"

	"bug_triage/internal/auth"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// JWTMiddleware validates JWT tokens from Authorization header
func JWTMiddleware(jwtManager *auth.JWTManager, logger *zap.Logger) gin.HandlerFunc {
	
	return func(c *gin.Context) {
		// logger.Debug("this is auth middleware")
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}
        
        
		parts := strings.Split(authHeader, " ") // Authorization: Bearer <token> --> ["Bearer", "<token>"]
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]

		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			logger.Debug("invalid token", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Store user info in context
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("claims", claims)

		c.Next()
	}
}


// GetUserID gets the authenticated user ID from context
func GetUserID(c *gin.Context) (int64, bool) {
	val, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	userID, ok := val.(int64)
	return userID, ok
}

// GetEmail gets the authenticated user email from context
func GetEmail(c *gin.Context) (string, bool) {
	val, exists := c.Get("email")
	if !exists {
		return "", false
	}
	email, ok := val.(string)
	return email, ok
}



// OptionalJWTMiddleware validates JWT tokens but doesn't fail if missing
// it is used when we create endpoint which can be used by signed-in and guest both type of user
// func OptionalJWTMiddleware(jwtManager *auth.JWTManager, logger *zap.Logger) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" {
// 			// Token optional, continue without user context
// 			c.Next()
// 			return
// 		}

// 		parts := strings.Split(authHeader, " ")
// 		if len(parts) != 2 || parts[0] != "Bearer" {
// 			// Invalid format, continue (don't fail)
// 			c.Next()
// 			return
// 		}

// 		token := parts[1]

// 		claims, err := jwtManager.ValidateToken(token)
// 		if err != nil {
// 			logger.Debug("invalid optional token", zap.Error(err))
// 			// Don't fail, just continue
// 			c.Next()
// 			return
// 		}

// 		// Store user info in context
// 		c.Set("user_id", claims.UserID)
// 		c.Set("email", claims.Email)
// 		c.Set("claims", claims)

// 		c.Next()
// 	}
// }
