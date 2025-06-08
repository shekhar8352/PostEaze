package controllers

import (
	"net/http"
	"posteaze-backend/services"
	"posteaze-backend/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	Service *services.AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{
		Service: &services.AuthService{},
	}
}

func (a *AuthController) SignupHandler(c *gin.Context) {
	var body services.SignupParams
	if err := c.ShouldBindJSON(&body); err != nil {
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
		return
	}

	user, err := a.Service.Login(c.Request.Context(), body)
	if err != nil {
		utils.SendError(c, http.StatusUnauthorized, err.Error())
		return
	}

	utils.SendSuccess(c, user)
}
