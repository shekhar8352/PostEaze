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

// GetPageDetailsHandler godoc
// @Summary      Get Instagram Page Details
// @Description  Fetches Instagram page details for a specific channel using the stored access token. Requires authentication.
// @Tags         Channels
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        channel_id query int true "Channel ID"
// @Success      200  {object}  modelsv1.GetPageDetailsResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /channels/details [get]
func GetPageDetailsHandler(c *gin.Context) {
	// Extract user_id from JWT token (set by AuthMiddleware)
	userIDStr, exists := c.Get("user_id")
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}

	// Get channel_id from query parameter
	var req modelsv1.GetPageDetailsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid or missing channel_id")
		return
	}

	resp, err := businessv1.GetPageDetails(c.Request.Context(), req.ChannelID, userIDStr.(string))
	if err != nil {
		if err.Error() == "unauthorized: channel does not belong to user" {
			utils.SendError(c, http.StatusForbidden, err.Error())
			return
		}
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccess(c, resp, "Page details retrieved successfully")
}

// SubscribeWebhooksHandler godoc
// @Summary      Subscribe to Meta Webhooks
// @Description  Subscribes an Instagram channel to Meta webhooks (comments, mentions, story_insights). Requires authentication.
// @Tags         Channels
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body modelsv1.SubscribeWebhooksRequest true "Subscribe Webhooks Request"
// @Success      200  {object}  modelsv1.SubscribeWebhooksResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /channels/instagram/subscribe-webhooks [post]
func SubscribeWebhooksHandler(c *gin.Context) {
	// Extract user_id from JWT token (set by AuthMiddleware)
	userIDStr, exists := c.Get("user_id")
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}

	// Bind request body
	var req modelsv1.SubscribeWebhooksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	resp, err := businessv1.SubscribeToWebhooks(c.Request.Context(), req.ChannelID, userIDStr.(string), req.Fields)
	if err != nil {
		if err.Error() == "unauthorized: channel does not belong to user" {
			utils.SendError(c, http.StatusForbidden, err.Error())
			return
		}
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccess(c, resp, "Successfully subscribed to webhooks")
}

// CreateFacebookChannelHandler godoc
// @Summary      Create Facebook Page channel
// @Description  Exchanges OAuth code, resolves the Page, and stores the Page access token for analytics.
// @Tags         Channels
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body modelsv1.CreateFacebookChannelRequest true "Create Facebook channel"
// @Success      200  {object}  modelsv1.CreateFacebookChannelResponse
// @Router       /channels/facebook/create [post]
func CreateFacebookChannelHandler(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}

	var req modelsv1.CreateFacebookChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	tokenFlow := req.PageID != "" && req.PageAccessToken != ""
	codeFlow := req.Code != "" && req.RedirectURI != "" && req.PageID != ""
	if !tokenFlow && !codeFlow {
		utils.SendError(c, http.StatusBadRequest, "provide page_id and page_access_token (after Meta callback), or code, redirect_uri, and page_id")
		return
	}

	var resp *modelsv1.CreateFacebookChannelResponse
	var err error
	if tokenFlow {
		resp, err = businessv1.CreateFacebookChannelFromPageToken(c.Request.Context(), req, userIDStr.(string))
	} else {
		resp, err = businessv1.CreateFacebookChannel(c.Request.Context(), req, userIDStr.(string))
	}
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccess(c, resp, "Facebook Page connected successfully")
}
