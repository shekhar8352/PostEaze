package apiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shekhar8352/PostEaze/utils"

	businessv1 "github.com/shekhar8352/PostEaze/business/v1"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
)

func CreateTeamHandler(c *gin.Context) {
	var team modelsv1.Team
	if err := c.ShouldBindJSON(&team); err != nil {
		utils.Logger.Info(c.Request.Context(), "Error binding JSON: ", err)
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	createdTeam, err := businessv1.CreateTeam(c.Request.Context(), team)
	if err != nil {
		utils.Logger.Error(c.Request.Context(), "Error creating team: %v", err)
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Logger.Info(c.Request.Context(), "Team created successfully")
	utils.SendSuccess(c, createdTeam, "Team created successfully")
}