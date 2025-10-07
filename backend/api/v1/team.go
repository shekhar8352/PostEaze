package apiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shekhar8352/PostEaze/utils"

	businessv1 "github.com/shekhar8352/PostEaze/business/v1"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
)

// CreateTeamHandler gdoc
// @Summary      Create team
// @Description  Create a team with name and owner ID
// @Tags         Teams
// @Accept       json
// @Produce      json
// @Param        team body modelsv1.Team true "Team details"
// @Success      200 {object} modelsv1.Team
// @Failure      400 {object} modelsv1.ErrorResponse
// @Failure      500 {object} modelsv1.ErrorResponse
// @Router       /teams [post]
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

// GetAllTeamsHandler godoc
// @Summary      Get all teams
// @Description  Get a list of all teams
// @Tags         Teams
// @Accept       json
// @Produce      json
// @Success      200 {array} modelsv1.Team
// @Failure      500 {object} modelsv1.ErrorResponse
// @Router       /teams/all [get]
func GetAllTeamsHandler(c *gin.Context) {
	teams, err := businessv1.GetAllTeams(c.Request.Context())
	if err != nil {
		utils.Logger.Error(c.Request.Context(), "Error fetching teams: %v", err)
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Logger.Info(c.Request.Context(), "Teams fetched successfully")
	utils.SendSuccess(c, teams, "Teams fetched successfully")
}

// GetTeamByIDHandler godoc
// @Summary      Get team by ID
// @Description  Get team details by ID
// @Tags         Teams
// @Accept       json
// @Produce      json
// @Param        team_id path      string  true  "Team ID"
// @Success      200 {object} modelsv1.Team
// @Failure      500 {object} modelsv1.ErrorResponse
// @Router       /teams/{team_id} [get]
func GetTeamByIDHandler(c *gin.Context) {
	teamID := c.Param("team_id")
	team, err := businessv1.GetTeamByID(c.Request.Context(), teamID)
	if err != nil {
		utils.Logger.Error(c.Request.Context(), "Error fetching team by ID: %v", err)
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Logger.Info(c.Request.Context(), "Team fetched successfully")
	utils.SendSuccess(c, team, "Team fetched successfully")
}

// GetTeamByOwnerIDHandler godoc
// @Summary      Get team by owner ID
// @Description  Get team details by owner ID
// @Tags         Teams
// @Accept       json
// @Produce      json
// @Param        owner_id path      string  true  "Owner ID"
// @Success      200 {object} modelsv1.Team
// @Failure      500 {object} modelsv1.ErrorResponse
// @Router       /teams/owner/{owner_id} [get]
func GetTeamByOwnerIDHandler(c *gin.Context) {
	ownerID := c.Param("owner_id")
	team, err := businessv1.GetTeamByOwnerID(c.Request.Context(), ownerID)
	if err != nil {
		utils.Logger.Error(c.Request.Context(), "Error fetching team by owner ID: %v", err)
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Logger.Info(c.Request.Context(), "Team fetched successfully")
	utils.SendSuccess(c, team, "Team fetched successfully")
}

// UpdateTeamHandler godoc
// @Summary      Update team
// @Description  Update team details by ID
// @Tags         Teams
// @Accept       json
// @Produce      json
// @Param        team_id  path      string  true  "Team ID"
// @Param        team body      modelsv1.Team true "Team details"
// @Success      200 {object} modelsv1.Team
// @Failure      400 {object} modelsv1.ErrorResponse
// @Failure      500 {object} modelsv1.ErrorResponse
// @Router       /teams/{team_id} [patch]
func UpdateTeamHandler(c *gin.Context) {
	teamID := c.Param("team_id")
	var team modelsv1.Team
	if err := c.ShouldBindJSON(&team); err != nil {
		utils.Logger.Info(c.Request.Context(), "Error binding JSON: ", err)
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	updatedTeam, err := businessv1.UpdateTeam(c.Request.Context(), team, teamID)
	if err != nil {
		utils.Logger.Error(c.Request.Context(), "Error updating team: %v", err)
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Logger.Info(c.Request.Context(), "Team updated successfully")
	utils.SendSuccess(c, updatedTeam, "Team updated successfully")
}