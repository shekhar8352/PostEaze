package apiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	businessv1 "github.com/shekhar8352/PostEaze/business/v1"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/utils"
)


// AuthenticateWithFirebaseHandler godoc
// @Summary      Authenticate user with Firebase
// @Description  Authenticate user using Firebase token and local ID, create user if not exists
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body modelsv1.FirebaseAuthParams true "Firebase Authentication Request"
// @Success      200 {object} modelsv1.SuccessResponse{data=map[string]interface{}}
// @Failure      400 {object} modelsv1.ErrorResponse
// @Failure      401 {object} modelsv1.ErrorResponse
// @Failure      500 {object} modelsv1.ErrorResponse
// @Router       /auth/authenticate [post]
func AuthenticateWithFirebaseHandler(c *gin.Context) {
	var body modelsv1.FirebaseAuthParams
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid authentication data")
		utils.Logger.Info(c.Request.Context(), "Error binding JSON: ", err)
		return
	}

	result, err := businessv1.AuthenticateWithFirebase(c.Request.Context(), body)
	if err != nil {
		if err.Error() == "invalid Firebase token" {
			utils.SendError(c, http.StatusUnauthorized, err.Error())
		} else if err.Error() == "email is required for this platform" {
			utils.SendError(c, http.StatusBadRequest, err.Error())
		} else {
			utils.SendError(c, http.StatusInternalServerError, err.Error())
		}
		utils.Logger.Info(c.Request.Context(), "Error authenticating user with Firebase: ", err)
		return
	}

	utils.SendSuccess(c, result, "Authenticated successfully")
	utils.Logger.Info(c.Request.Context(), "Authenticated user with Firebase successfully")
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