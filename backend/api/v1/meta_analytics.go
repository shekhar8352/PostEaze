package apiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	businessv1 "github.com/shekhar8352/PostEaze/business/v1"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/utils"
)

// SyncMetaAnalyticsHandler godoc
// @Summary      Sync Meta analytics
// @Description  Fetches Instagram and/or Facebook Page insights for channels owned by the user.
// @Tags         Meta
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body modelsv1.SyncMetaAnalyticsRequest false "Optional channel_ids filter"
// @Success      200  {object}  modelsv1.SyncMetaAnalyticsResponse
// @Router       /meta/analytics/sync [post]
func SyncMetaAnalyticsHandler(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}

	var req modelsv1.SyncMetaAnalyticsRequest
	_ = c.ShouldBindJSON(&req)

	resp, err := businessv1.SyncMetaAnalyticsForUser(c.Request.Context(), userIDStr.(string), req.ChannelIDs)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccess(c, resp, "Meta analytics sync completed")
}
