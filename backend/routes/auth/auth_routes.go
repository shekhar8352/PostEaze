package auth

import (
	"github.com/gin-gonic/gin"
	"posteaze-backend/controllers"
)

func RegisterAuthRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")

	authController := controllers.NewAuthController()

	auth.POST("/signup", authController.SignupHandler)
	auth.POST("/login", authController.LoginHandler)
}
