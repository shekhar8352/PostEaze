package apiv1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	businessv1 "github.com/shekhar8352/PostEaze/business/v1"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/utils"
)

// SignupHandler godoc
// @Summary      User Signup
// @Description  Create a new user account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body modelsv1.SignupParams true "Signup Request"
// @Success      200 {object} modelsv1.SuccessResponse
// @Failure      400 {object} modelsv1.ErrorResponse
// @Failure      500 {object} modelsv1.ErrorResponse
// @Router       /auth/signup [post]
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

	utils.SendSuccess(c, user, "Signed up successfully")
}

// LoginHandler godoc
// @Summary      User Login
// @Description  Authenticate user and return tokens
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body modelsv1.LoginParams true "Login Request"
// @Success      200 {object} modelsv1.SuccessResponse
// @Failure      400 {object} modelsv1.ErrorResponse
// @Failure      401 {object} modelsv1.ErrorResponse
// @Router       /auth/login [post]
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

	utils.SendSuccess(c, user, "Logged in successfully")
	utils.Logger.Info(c.Request.Context(), "Logged in user successfully: %s", user)
}

// RefreshTokenHandler godoc
// @Summary      Refresh Token
// @Description  Generate new access token using refresh token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        request body modelsv1.RefreshTokenParams true "Refresh Token Request"
// @Success      200 {object} modelsv1.SuccessResponse
// @Failure      400 {object} modelsv1.ErrorResponse
// @Failure      401 {object} modelsv1.ErrorResponse
// @Router       /auth/refresh [post]
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

	utils.SendSuccess(c, user, "Refreshed token successfully")
	utils.Logger.Info(c.Request.Context(), "Refreshed token successfully: ", user)
}

// LogoutHandler godoc
// @Summary      User Logout
// @Description  Invalidate user session and token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200 {object} modelsv1.SuccessResponse
// @Failure      500 {object} modelsv1.ErrorResponse
// @Router       /auth/logout [post]
func LogoutHandler(c *gin.Context) {
	err := businessv1.Logout(c.Request.Context(), c.GetHeader("Authorization"))
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		utils.Logger.Info(c.Request.Context(), "Error logging out user: ", err)
		return
	}

	utils.SendSuccess(c, nil, "Logged out successfully")
}
