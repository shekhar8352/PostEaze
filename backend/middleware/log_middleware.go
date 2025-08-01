package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shekhar8352/PostEaze/utils"
)

// GinLoggingMiddleware is a Gin-compatible logging middleware
func GinLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Record start time
		startTime := time.Now()

		// Add log ID to the request context
		ctx := utils.AddLogID(c.Request.Context())

		// Update the request context in Gin
		c.Request = c.Request.WithContext(ctx)

		// Get the log ID for consistent logging
		logID := utils.GetLogID(ctx)

		// Log the incoming request
		utils.Logger.Info(ctx, "Started %s %s | IP: %s | User-Agent: %s",
			c.Request.Method, c.Request.URL.Path, getClientIP(c), c.Request.UserAgent())

		// Continue to next handler
		c.Next()

		// Calculate duration
		duration := time.Since(startTime)

		// Log request completion with status and duration
		utils.Logger.Info(ctx, "Completed %s %s | Status: %d | Duration: %v | LogID: %s",
			c.Request.Method, c.Request.URL.Path, c.Writer.Status(), duration, logID)
	}
}

// getClientIP extracts the real client IP from Gin context
func getClientIP(c *gin.Context) string {
	// Gin has a built-in method for this, but we can use our custom logic
	// Check X-Forwarded-For header (for load balancers/proxies)
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		return xff
	}

	// Check X-Real-IP header
	if xri := c.GetHeader("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to ClientIP (Gin's method)
	return c.ClientIP()
}
