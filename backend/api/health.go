package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	appdb "github.com/shekhar8352/PostEaze/utils/database"
	redisutil "github.com/shekhar8352/PostEaze/utils/redis"
)

const healthCheckTimeout = 3 * time.Second

func registerHealthRoutes(api *gin.RouterGroup) {
	api.GET("/health", healthLivenessHandler)
	api.GET("/health/postgres", healthPostgresHandler)
	api.GET("/health/redis", healthRedisHandler)
	api.GET("/health/ready", healthReadyHandler)
}

func healthLivenessHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func healthPostgresHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), healthCheckTimeout)
	defer cancel()

	if err := appdb.Ping(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "error",
			"service": "postgres",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "postgres",
	})
}

func healthRedisHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), healthCheckTimeout)
	defer cancel()

	if err := redisutil.Ping(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "error",
			"service": "redis",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "redis",
	})
}

func healthReadyHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), healthCheckTimeout)
	defer cancel()

	pgErr := appdb.Ping(ctx)
	redisErr := redisutil.Ping(ctx)

	checks := gin.H{
		"postgres": serviceStatus(pgErr),
		"redis":    serviceStatus(redisErr),
	}

	if pgErr != nil || redisErr != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "not_ready",
			"checks":  checks,
			"message": "one or more dependencies are unavailable",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
		"checks": checks,
	})
}

func serviceStatus(err error) gin.H {
	if err == nil {
		return gin.H{"status": "up"}
	}
	return gin.H{
		"status": "down",
		"error":  err.Error(),
	}
}
