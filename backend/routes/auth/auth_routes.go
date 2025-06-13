package auth

import (
	"posteaze-backend/controllers"
	"posteaze-backend/services"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")

	authService := services.NewAuthService()
	authController := controllers.NewAuthController(authService)

	auth.POST("/signup", authController.SignupHandler)
	auth.POST("/login", authController.LoginHandler)
	auth.POST("/refresh", authController.RefreshTokenHandler)
	auth.POST("/logout", authController.LogoutHandler)
}
