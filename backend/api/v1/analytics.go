package apiv1

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shekhar8352/PostEaze/entities/repositories"
	"github.com/shekhar8352/PostEaze/utils"
)

// GetProfileAnalyticsHandler godoc
// @Summary      Get Profile Analytics
// @Description  Get Instagram profile analytics for a channel within a date range
// @Tags         Analytics
// @Accept       json
// @Produce      json
// @Param        channelId path int true "Channel ID"
// @Param        start_date query string false "Start date (YYYY-MM-DD)" default(7 days ago)
// @Param        end_date query string false "End date (YYYY-MM-DD)" default(today)
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /channels/{channelId}/analytics/profile [get]
// @Security     BearerAuth
func GetProfileAnalyticsHandler(c *gin.Context) {
	channelID, err := strconv.ParseInt(c.Param("channelId"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid channel ID")
		return
	}

	// Parse date range
	startDate, endDate := utils.ParseDateRange(c)

	analytics, err := repositories.GetProfileAnalyticsByDateRange(c.Request.Context(), channelID, startDate, endDate)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to fetch profile analytics: "+err.Error())
		return
	}

	utils.SendSuccess(c, gin.H{
		"analytics":  analytics,
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
	}, "Profile analytics retrieved successfully")
}

// GetPostAnalyticsHandler godoc
// @Summary      Get Post Analytics
// @Description  Get Instagram post analytics for a channel within a date range
// @Tags         Analytics
// @Accept       json
// @Produce      json
// @Param        channelId path int true "Channel ID"
// @Param        start_date query string false "Start date (YYYY-MM-DD)"
// @Param        end_date query string false "End date (YYYY-MM-DD)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /channels/{channelId}/analytics/posts [get]
// @Security     BearerAuth
func GetPostAnalyticsHandler(c *gin.Context) {
	channelID, err := strconv.ParseInt(c.Param("channelId"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid channel ID")
		return
	}

	startDate, endDate := utils.ParseDateRange(c)

	analytics, err := repositories.GetPostAnalyticsByDateRange(c.Request.Context(), channelID, startDate, endDate)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to fetch post analytics: "+err.Error())
		return
	}

	utils.SendSuccess(c, gin.H{
		"analytics":  analytics,
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
	}, "Post analytics retrieved successfully")
}

// GetAnalyticsOverviewHandler godoc
// @Summary      Get Analytics Overview
// @Description  Get aggregated analytics overview for a channel
// @Tags         Analytics
// @Accept       json
// @Produce      json
// @Param        channelId path int true "Channel ID"
// @Param        start_date query string false "Start date (YYYY-MM-DD)"
// @Param        end_date query string false "End date (YYYY-MM-DD)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /channels/{channelId}/analytics/overview [get]
// @Security     BearerAuth
func GetAnalyticsOverviewHandler(c *gin.Context) {
	channelID, err := strconv.ParseInt(c.Param("channelId"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid channel ID")
		return
	}

	startDate, endDate := utils.ParseDateRange(c)

	aggregated, err := repositories.GetAggregatedProfileAnalytics(c.Request.Context(), channelID, startDate, endDate)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to fetch aggregated analytics: "+err.Error())
		return
	}

	utils.SendSuccess(c, gin.H{
		"overview":   aggregated,
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
	}, "Analytics overview retrieved successfully")
}

// GetTopPostsHandler godoc
// @Summary      Get Top Posts
// @Description  Get top performing posts by engagement
// @Tags         Analytics
// @Accept       json
// @Produce      json
// @Param        channelId path int true "Channel ID"
// @Param        limit query int false "Number of posts to return" default(10)
// @Param        start_date query string false "Start date (YYYY-MM-DD)"
// @Param        end_date query string false "End date (YYYY-MM-DD)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /channels/{channelId}/analytics/top-posts [get]
// @Security     BearerAuth
func GetTopPostsHandler(c *gin.Context) {
	channelID, err := strconv.ParseInt(c.Param("channelId"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid channel ID")
		return
	}

	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	startDate, endDate := utils.ParseDateRange(c)

	topPosts, err := repositories.GetTopPosts(c.Request.Context(), channelID, limit, startDate, endDate)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to fetch top posts: "+err.Error())
		return
	}

	utils.SendSuccess(c, gin.H{
		"top_posts":  topPosts,
		"limit":      limit,
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
	}, "Top posts retrieved successfully")
}

// GetPostsOverviewHandler godoc
// @Summary      Get Posts Overview
// @Description  Get overview of posts activity including new posts, likes, comments, etc.
// @Tags         Analytics
// @Accept       json
// @Produce      json
// @Param        channelId path int true "Channel ID"
// @Param        start_date query string false "Start date (YYYY-MM-DD)"
// @Param        end_date query string false "End date (YYYY-MM-DD)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /channels/{channelId}/analytics/posts-overview [get]
// @Security     BearerAuth
func GetPostsOverviewHandler(c *gin.Context) {
	channelID, err := strconv.ParseInt(c.Param("channelId"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid channel ID")
		return
	}

	startDate, endDate := utils.ParseDateRange(c)

	overview, err := repositories.GetPostsOverview(c.Request.Context(), channelID, startDate, endDate)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to fetch posts overview: "+err.Error())
		return
	}

	utils.SendSuccess(c, gin.H{
		"overview":   overview,
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
	}, "Posts overview retrieved successfully")
}

// GetPostInsightsHandler godoc
// @Summary      Get Single Post Insights
// @Description  Get detailed historical analytics for a single post
// @Tags         Analytics
// @Param        channelId path int true "Channel ID"
// @Param        postId path int true "Post ID"
// @Param        days query int false "Days of history" default(7)
// @Router       /channels/{channelId}/analytics/posts/{postId} [get]
func GetPostInsightsHandler(c *gin.Context) {
	_, err := strconv.ParseInt(c.Param("channelId"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid channel ID")
		return
	}
	postID, err := strconv.ParseInt(c.Param("postId"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid post ID")
		return
	}

	limit := 7
	if lStr := c.Query("days"); lStr != "" {
		if l, err := strconv.Atoi(lStr); err == nil && l > 0 {
			limit = l
		}
	}

	analytics, err := repositories.GetPostDetailedAnalytics(c.Request.Context(), postID, limit)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to fetch post insights: "+err.Error())
		return
	}

	utils.SendSuccess(c, gin.H{
		"analytics": analytics,
	}, "Post insights retrieved successfully")
}

// GetChannelDashboardHandler godoc
// @Summary      Get Channel Dashboard
// @Description  Get comprehensive channel dashboard metrics
// @Router       /channels/{channelId}/analytics/dashboard [get]
func GetChannelDashboardHandler(c *gin.Context) {
	GetAnalyticsOverviewHandler(c)
}

// GetPeriodComparisonHandler godoc
// @Summary      Get Period Comparison
// @Description  Compare analytics between selected period and previous period
// @Param        channelId path int true "Channel ID"
// @Param        start_date query string false "Start date"
// @Param        end_date query string false "End date"
// @Router       /channels/{channelId}/analytics/comparison [get]
func GetPeriodComparisonHandler(c *gin.Context) {
	channelID, err := strconv.ParseInt(c.Param("channelId"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid channel ID")
		return
	}
	startDate, endDate := utils.ParseDateRange(c)

	currentPeriod, err := repositories.GetAggregatedProfileAnalytics(c.Request.Context(), channelID, startDate, endDate)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed: "+err.Error())
		return
	}

	duration := endDate.Sub(startDate)
	prevStart := startDate.Add(-duration - 24*time.Hour) // simplistic logic
	prevEnd := startDate.Add(-time.Hour * 24)

	prevPeriod, err := repositories.GetAggregatedProfileAnalytics(c.Request.Context(), channelID, prevStart, prevEnd)
	if err != nil {
		// Just log error and return current? Or fail?
		// Return current with empty comparison
		utils.SendSuccess(c, gin.H{"current": currentPeriod, "previous": nil}, "Comparison retrieved (no previous data)")
		return
	}

	utils.SendSuccess(c, gin.H{
		"current":  currentPeriod,
		"previous": prevPeriod,
		"period": gin.H{
			"start": startDate, "end": endDate,
			"prev_start": prevStart, "prev_end": prevEnd,
		},
	}, "Comparison retrieved successfully")
}
