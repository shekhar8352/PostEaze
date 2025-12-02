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
// @Description  Creates an Instagram channel by exchanging authorization code for token. Requires authentication via JWT token. Accepts optional metadata dictionary.
// @Tags         Channels
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body modelsv1.CreateInstagramChannelRequest true "Create Instagram Channel Request"
// @Success      200  {object}  modelsv1.CreateInstagramChannelResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /channels/instagram/create [post]
func CreateInstagramChannelHandler(c *gin.Context) {
	var req modelsv1.CreateInstagramChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Extract user_id from JWT token (set by AuthMiddleware)
	userIDStr, exists := c.Get("user_id")
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}

	resp, err := businessv1.CreateInstagramChannel(c.Request.Context(), req, userIDStr.(string))
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccess(c, resp, "Instagram channel created successfully")
}

// GetChannelsHandler godoc
// @Summary      Get User Channels
// @Description  Retrieves all channels for the authenticated user. Optionally filter by provider (e.g., instagram, facebook).
// @Tags         Channels
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        provider query string false "Filter by provider (instagram, facebook, etc.)"
// @Success      200  {object}  modelsv1.GetChannelsResponse
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /channels [get]
func GetChannelsHandler(c *gin.Context) {
	// Extract user_id from JWT token (set by AuthMiddleware)
	userIDStr, exists := c.Get("user_id")
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}

	// Get optional provider filter from query params
	var req modelsv1.GetChannelsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := businessv1.GetChannels(c.Request.Context(), userIDStr.(string), req.Provider)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccess(c, resp, "Channels retrieved successfully")
}
