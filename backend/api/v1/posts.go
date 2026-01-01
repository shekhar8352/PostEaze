package apiv1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shekhar8352/PostEaze/entities/repositories"
	"github.com/shekhar8352/PostEaze/utils"
)

// GetPostsHandler godoc
// @Summary      Get Posts
// @Description  Get posts with filtering options
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Param        channel_ids query string false "Comma separated channel IDs"
// @Param        provider query string false "Provider name (e.g. instagram)"
// @Param        limit query int false "Limit" default(20)
// @Param        offset query int false "Offset"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /posts [get]
// @Security     BearerAuth
func GetPostsHandler(c *gin.Context) {
	filters := repositories.PostFilters{}

	// Parse channel_ids (Assuming single ID for now based on PostFilters struct, but can extend later)
	// The new schema supports channel_ids array, but PostFilters has ChannelID *int64.
	// We will support single channel filter for now as per repository support.
	if cidStr := c.Query("channel_id"); cidStr != "" {
		if cid, err := strconv.ParseInt(cidStr, 10, 64); err == nil {
			filters.ChannelID = &cid
		}
	}

	if provider := c.Query("provider"); provider != "" {
		filters.Provider = &provider
	}

	filters.Limit = 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			filters.Limit = l
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			filters.Offset = o
		}
	}

	posts, err := repositories.GetPosts(c.Request.Context(), filters)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to fetch posts: "+err.Error())
		return
	}

	utils.SendSuccess(c, gin.H{
		"posts": posts,
		"count": len(posts),
	}, "Posts retrieved successfully")
}
