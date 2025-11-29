package apiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	businessv1 "github.com/shekhar8352/PostEaze/business/v1"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/utils"
)

// CreateInstagramChannelHandler godoc
// @Summary      Create Instagram Channel
// @Description  Creates an Instagram channel by exchanging authorization code for token
// @Tags         Channels
// @Accept       json
// @Produce      json
// @Param        request body modelsv1.CreateInstagramChannelRequest true "Create Instagram Channel Request"
// @Success      200  {object}  modelsv1.CreateInstagramChannelResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /channels/instagram/create [post]
func CreateInstagramChannelHandler(c *gin.Context) {
	var req modelsv1.CreateInstagramChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := businessv1.CreateInstagramChannel(c.Request.Context(), req)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccess(c, resp, "Instagram channel created successfully")
}
