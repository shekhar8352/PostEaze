package controllers

import (
	"fmt"
	"net/http"
	"posteaze-backend/services"
	"posteaze-backend/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	Service *services.AuthService
}

func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{Service: service}
}

func (a *AuthController) SignupHandler(c *gin.Context) {
	var body services.SignupParams
	if err := c.ShouldBindJSON(&body); err != nil {
		fmt.Println(err)
		utils.SendError(c, http.StatusBadRequest, "Invalid signup data")
		return
	}

	user, err := a.Service.Signup(c.Request.Context(), body)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccess(c, user)
}

func (a *AuthController) LoginHandler(c *gin.Context) {
	var body services.LoginParams
	if err := c.ShouldBindJSON(&body); err != nil || body.Email == "" || body.Password == "" {
		utils.SendError(c, http.StatusBadRequest, "Email and password are required")
		utils.Logger.Info("Error binding JSON: ", err)
		return
	}

	user, err := a.Service.Login(c.Request.Context(), body)
	if err != nil {
		utils.SendError(c, http.StatusUnauthorized, err.Error())
		utils.Logger.Info("Error logging in user: ", err)
		return
	}

	utils.SendSuccess(c, user)
	utils.Logger.Info("Logged in user successfully: ", user)
}
