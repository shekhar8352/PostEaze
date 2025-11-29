package apiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shekhar8352/PostEaze/utils"

	businessv1 "github.com/shekhar8352/PostEaze/business/v1"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
)

// GetUserByIdHandler gdoc
// @Summary      Get user by ID
// @Description  Get user details by ID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user_id  path      string  true  "User ID"
// @Security     ApiKeyAuth
// @Success      200 {object} modelsv1.User
// @Failure      500 {object} map[string]interface{}
// @Router       /users/{user_id} [get]
func GetUserByIdHandler(c *gin.Context) {
	userID := c.Param("user_id")
	user, err := businessv1.GetUserById(c.Request.Context(), userID)
	if err != nil {
		utils.Logger.Error(c.Request.Context(), "Error fetching user by ID: %v", err)
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Logger.Info(c.Request.Context(), "User fetched successfully")
	utils.SendSuccess(c, user, "User fetched successfully")
}

// UpdateUserHandler gdoc
// @Summary      Update user
// @Description  Update user details
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request body modelsv1.UpdateUserParams true "Update user data"
// @Security     ApiKeyAuth
// @Success      200 {object} modelsv1.User
// @Failure      500 {object} map[string]interface{}
// @Router       /users/{user_id} [patch]
func UpdateUserHandler(c *gin.Context) {
	var body modelsv1.UpdateUserParams
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid update user data")
		utils.Logger.Info(c.Request.Context(), "Error binding JSON: ", err)
		return
	}

	user, err := businessv1.UpdateUser(c.Request.Context(), body)
	if err != nil {
		utils.Logger.Error(c.Request.Context(), "Error updating user: %v", err)
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Logger.Info(c.Request.Context(), "User updated successfully")
	utils.SendSuccess(c, user, "User updated successfully")
}