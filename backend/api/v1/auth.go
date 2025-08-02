package apiv1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	businessv1 "github.com/shekhar8352/PostEaze/business/v1"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/utils"
)

func SignupHandler(c *gin.Context) {
	var body modelsv1.SignupParams
	if err := c.ShouldBindJSON(&body); err != nil {
		fmt.Println(err)
		utils.SendError(c, http.StatusBadRequest, "Invalid signup data")
		return
	}

	user, err := businessv1.Signup(c.Request.Context(), body)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccess(c, user)
}

func LoginHandler(c *gin.Context) {
	var body modelsv1.LoginParams
	if err := c.ShouldBindJSON(&body); err != nil || body.Email == "" || body.Password == "" {
		utils.SendError(c, http.StatusBadRequest, "Email and password are required")
		utils.Logger.Info(c.Request.Context(), "Error binding JSON: ", err)
		return
	}

	user, err := businessv1.Login(c.Request.Context(), body)
	if err != nil {
		utils.SendError(c, http.StatusUnauthorized, err.Error())
		utils.Logger.Info(c.Request.Context(), "Error logging in user: ", err)
		return
	}

	utils.SendSuccess(c, user)
	utils.Logger.Info(c.Request.Context(), "Logged in user successfully: %s", user)
}

func RefreshTokenHandler(c *gin.Context) {
	var body modelsv1.RefreshTokenParams
	if err := c.ShouldBindJSON(&body); err != nil || body.RefreshToken == "" {
		utils.SendError(c, http.StatusBadRequest, "Refresh token is required")
		utils.Logger.Info(c.Request.Context(), "Error binding JSON: ", err)
		return
	}

	user, err := businessv1.RefreshToken(c.Request.Context(), body.RefreshToken)
	if err != nil {
		utils.SendError(c, http.StatusUnauthorized, err.Error())
		utils.Logger.Info(c.Request.Context(), "Error refreshing token: ", err)
		return
	}

	utils.SendSuccess(c, user)
	utils.Logger.Info(c.Request.Context(), "Refreshed token successfully: ", user)
}

func LogoutHandler(c *gin.Context) {
	err := businessv1.Logout(c.Request.Context(), c.GetHeader("Authorization"))
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		utils.Logger.Info(c.Request.Context(), "Error logging out user: ", err)
		return
	}

	utils.SendSuccess(c, "Logged out successfully")
}
