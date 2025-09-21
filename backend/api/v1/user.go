package apiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shekhar8352/PostEaze/utils"

	businessv1 "github.com/shekhar8352/PostEaze/business/v1"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
)

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