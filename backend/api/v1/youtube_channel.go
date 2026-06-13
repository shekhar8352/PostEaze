package apiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	businessv1 "github.com/shekhar8352/PostEaze/business/v1"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/utils"
)

func CreateYouTubeChannelHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	var req modelsv1.CreateYouTubeChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := businessv1.CreateYouTubeChannel(c.Request.Context(), req, userIDStr.(string))
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendSuccess(c, resp, "YouTube channel connected")
}
