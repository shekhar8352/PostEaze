package routes

import (
	"github.com/gin-gonic/gin"
	"posteaze-backend/routes/auth"
)

func RegisterRoutes(router *gin.Engine) {
	// API Group
	api := router.Group("/api")

	// Register Auth Routes
	auth.RegisterAuthRoutes(api)
}
