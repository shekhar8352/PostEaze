package apiv1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	businessv1 "github.com/shekhar8352/PostEaze/business/v1"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/utils"
)

// ListScheduledPostsHandler godoc
// @Summary      List scheduled posts (calendar range)
// @Tags         ScheduledPosts
// @Security     BearerAuth
// @Param        from query string true "Start (RFC3339 or YYYY-MM-DD)"
// @Param        to query string true "End (RFC3339 or YYYY-MM-DD)"
// @Param        channel_id query int false "Filter by channel"
// @Success      200 {object} map[string]interface{}
// @Router       /scheduled-posts [get]
func ListScheduledPostsHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	var q modelsv1.ListScheduledPostsQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	resp, code, err := businessv1.ListScheduledPostsForCalendar(c.Request.Context(), userIDStr.(string), q)
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, resp, "Scheduled posts retrieved")
}

// CreateScheduledPostHandler godoc
// @Summary      Create and submit scheduled post(s) to Instagram
// @Tags         ScheduledPosts
// @Security     BearerAuth
// @Param        body body modelsv1.CreateScheduledPostRequest true "Request"
// @Success      200 {object} map[string]interface{}
// @Success      207 {object} map[string]interface{}
// @Failure      501 {object} map[string]interface{}
// @Router       /scheduled-posts [post]
func CreateScheduledPostHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	var req modelsv1.CreateScheduledPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	resp, code, err := businessv1.CreateScheduledPost(c.Request.Context(), userIDStr.(string), &req)
	if err != nil {
		if code == 501 {
			utils.SendError(c, http.StatusNotImplemented, err.Error())
			return
		}
		utils.SendError(c, code, err.Error())
		return
	}
	switch resp.OverallStatus {
	case "failed":
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": "error",
			"msg":    "scheduled post failed",
			"data":   resp,
		})
	case "partial_failure":
		c.JSON(http.StatusMultiStatus, gin.H{
			"status": "warning",
			"msg":    "scheduled post partially processed",
			"data":   resp,
		})
	default:
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"msg":    "scheduled post processed",
			"data":   resp,
		})
	}
}

// GetScheduledPostHandler godoc
// @Summary      Get one scheduled post
// @Tags         ScheduledPosts
// @Security     BearerAuth
// @Param        id path int true "Scheduled post ID"
// @Router       /scheduled-posts/{id} [get]
func GetScheduledPostHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid id")
		return
	}
	item, code, err := businessv1.GetScheduledPost(c.Request.Context(), userIDStr.(string), id)
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, item, "OK")
}

// CancelScheduledPostHandler godoc
// @Summary      Cancel a scheduled post (local + allowed statuses)
// @Tags         ScheduledPosts
// @Security     BearerAuth
// @Param        id path int true "Scheduled post ID"
// @Router       /scheduled-posts/{id} [delete]
func CancelScheduledPostHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid id")
		return
	}
	code, err := businessv1.CancelScheduledPost(c.Request.Context(), userIDStr.(string), id)
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, gin.H{"id": id}, "Cancelled")
}
