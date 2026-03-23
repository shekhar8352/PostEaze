package apiv1

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shekhar8352/PostEaze/entities/repositories"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/utils"
)

// GetProfileAnalyticsHandler godoc
// @Summary      Get Profile Analytics
// @Description  Get Instagram profile analytics for a channel within a date range
// @Tags         Analytics
// @Param        channelId path int true "Channel ID"
// @Param        start_date query string false "Start date (YYYY-MM-DD)"
// @Param        end_date query string false "End date (YYYY-MM-DD)"
// @Param        limit query int false "Page size (default 100, max 500; 0 = all)"
// @Param        offset query int false "Offset"
// @Param        include_raw query string false "1 to include raw Meta JSON per row"
// @Router       /channels/{channelId}/analytics/profile [get]
// @Security     BearerAuth
func GetProfileAnalyticsHandler(c *gin.Context) {
	channelID, err := strconv.ParseInt(c.Param("channelId"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid channel ID")
		return
	}

	startDate, endDate := utils.ParseDateRange(c)
	limit, offset := utils.ParseAnalyticsPagination(c)
	includeRaw := utils.AnalyticsIncludeRawQuery(c)

	f := repositories.ProfileAnalyticsListFilters{Limit: limit, Offset: offset}
	total, err := repositories.CountProfileAnalyticsInRange(c.Request.Context(), channelID, startDate, endDate)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to count profile analytics: "+err.Error())
		return
	}

	rows, err := repositories.GetProfileAnalyticsByDateRange(c.Request.Context(), channelID, startDate, endDate, f)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to fetch profile analytics: "+err.Error())
		return
	}

	resp := modelsv1.ProfileAnalyticsResponse{
		Meta: modelsv1.DateRangeMeta{
			StartDate: utils.FormatAnalyticsDate(startDate),
			EndDate:   utils.FormatAnalyticsDate(endDate),
		},
		Pagination: modelsv1.PaginationMeta{
			Limit:  limit,
			Offset: offset,
			Total:  total,
		},
		Analytics: utils.MapProfileAnalyticsItems(rows, includeRaw),
	}
	utils.SendSuccess(c, resp, "Profile analytics retrieved successfully")
}

// GetPostAnalyticsHandler godoc
// @Summary      Get Post Analytics
// @Description  Get Instagram post analytics for a channel within a date range
// @Tags         Analytics
// @Param        channelId path int true "Channel ID"
// @Param        start_date query string false "Start date (YYYY-MM-DD)"
// @Param        end_date query string false "End date (YYYY-MM-DD)"
// @Param        post_type query string false "Filter by posts.post_type"
// @Param        limit query int false "Page size (default 100, max 500; 0 = all)"
// @Param        offset query int false "Offset"
// @Param        include_raw query string false "1 to include raw Meta JSON per row"
// @Router       /channels/{channelId}/analytics/posts [get]
// @Security     BearerAuth
func GetPostAnalyticsHandler(c *gin.Context) {
	channelID, err := strconv.ParseInt(c.Param("channelId"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid channel ID")
		return
	}

	startDate, endDate := utils.ParseDateRange(c)
	limit, offset := utils.ParseAnalyticsPagination(c)
	postType := utils.OptionalPostTypeQuery(c)
	includeRaw := utils.AnalyticsIncludeRawQuery(c)

	f := repositories.PostAnalyticsListFilters{
		PostType: postType,
		Limit:    limit,
		Offset:   offset,
	}
	total, err := repositories.CountPostAnalyticsInRange(c.Request.Context(), channelID, startDate, endDate, postType)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to count post analytics: "+err.Error())
		return
	}

	rows, err := repositories.GetPostAnalyticsByDateRange(c.Request.Context(), channelID, startDate, endDate, f)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to fetch post analytics: "+err.Error())
		return
	}

	resp := modelsv1.PostAnalyticsResponse{
		Meta: modelsv1.DateRangeMeta{
			StartDate: utils.FormatAnalyticsDate(startDate),
			EndDate:   utils.FormatAnalyticsDate(endDate),
		},
		Pagination: modelsv1.PaginationMeta{
			Limit:  limit,
			Offset: offset,
			Total:  total,
		},
		Analytics: utils.MapPostAnalyticsItems(rows, includeRaw),
	}
	utils.SendSuccess(c, resp, "Post analytics retrieved successfully")
}

// GetAnalyticsOverviewHandler godoc
// @Summary      Get Analytics Overview
// @Tags         Analytics
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

	resp := modelsv1.OverviewResponse{
		Meta: modelsv1.DateRangeMeta{
			StartDate: utils.FormatAnalyticsDate(startDate),
			EndDate:   utils.FormatAnalyticsDate(endDate),
		},
		Overview: utils.MapAggregatedProfileOverview(aggregated),
	}
	utils.SendSuccess(c, resp, "Analytics overview retrieved successfully")
}

// GetTopPostsHandler godoc
// @Summary      Get Top Posts
// @Param        sort query string false "engagement|reach|impressions|plays" default(engagement)
// @Param        post_type query string false "Filter by post_type"
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
			if limit > utils.AnalyticsMaxPageSize {
				limit = utils.AnalyticsMaxPageSize
			}
		}
	}

	startDate, endDate := utils.ParseDateRange(c)
	sort := utils.ParseTopPostsSort(c)
	postType := utils.OptionalPostTypeQuery(c)

	topPosts, err := repositories.GetTopPosts(c.Request.Context(), channelID, limit, startDate, endDate, sort, postType)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to fetch top posts: "+err.Error())
		return
	}

	resp := modelsv1.TopPostsResponse{
		Meta: modelsv1.DateRangeMeta{
			StartDate: utils.FormatAnalyticsDate(startDate),
			EndDate:   utils.FormatAnalyticsDate(endDate),
		},
		Sort:     string(sort),
		Limit:    limit,
		TopPosts: utils.MapTopPostItems(topPosts),
	}
	utils.SendSuccess(c, resp, "Top posts retrieved successfully")
}

// GetPostsOverviewHandler godoc
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

	resp := modelsv1.PostsOverviewResponse{
		Meta: modelsv1.DateRangeMeta{
			StartDate: utils.FormatAnalyticsDate(startDate),
			EndDate:   utils.FormatAnalyticsDate(endDate),
		},
		Overview: utils.MapPostsOverview(overview),
	}
	utils.SendSuccess(c, resp, "Posts overview retrieved successfully")
}

// GetPostInsightsHandler godoc
// @Router       /channels/{channelId}/analytics/posts/{postId} [get]
// @Security     BearerAuth
func GetPostInsightsHandler(c *gin.Context) {
	channelID, err := strconv.ParseInt(c.Param("channelId"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid channel ID")
		return
	}
	postID, err := strconv.ParseInt(c.Param("postId"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid post ID")
		return
	}

	ok, err := repositories.PostBelongsToChannel(c.Request.Context(), postID, channelID)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to verify post: "+err.Error())
		return
	}
	if !ok {
		utils.SendError(c, http.StatusNotFound, "Post not found for this channel")
		return
	}

	limit := 7
	if lStr := c.Query("days"); lStr != "" {
		if l, err := strconv.Atoi(lStr); err == nil && l > 0 {
			limit = l
		}
	}

	includeRaw := utils.AnalyticsIncludeRawQuery(c)

	analytics, err := repositories.GetPostDetailedAnalytics(c.Request.Context(), postID, limit)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to fetch post insights: "+err.Error())
		return
	}

	resp := modelsv1.PostInsightsResponse{
		Analytics: utils.MapPostAnalyticsItems(analytics, includeRaw),
	}
	utils.SendSuccess(c, resp, "Post insights retrieved successfully")
}

// GetChannelDashboardHandler godoc
// @Summary      Get Channel Dashboard
// @Router       /channels/{channelId}/analytics/dashboard [get]
// @Security     BearerAuth
func GetChannelDashboardHandler(c *gin.Context) {
	channelID, err := strconv.ParseInt(c.Param("channelId"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid channel ID")
		return
	}

	startDate, endDate := utils.ParseDateRange(c)

	overview, err := repositories.GetAggregatedProfileAnalytics(c.Request.Context(), channelID, startDate, endDate)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to fetch overview: "+err.Error())
		return
	}

	topPosts, err := repositories.GetTopPosts(c.Request.Context(), channelID, 5, startDate, endDate, repositories.TopPostsSortEngagement, nil)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to fetch top posts: "+err.Error())
		return
	}

	postsOverview, err := repositories.GetPostsOverview(c.Request.Context(), channelID, startDate, endDate)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to fetch posts overview: "+err.Error())
		return
	}

	resp := modelsv1.DashboardResponse{
		Meta: modelsv1.DateRangeMeta{
			StartDate: utils.FormatAnalyticsDate(startDate),
			EndDate:   utils.FormatAnalyticsDate(endDate),
		},
		Overview:      utils.MapAggregatedProfileOverview(overview),
		TopPosts:      utils.MapTopPostItems(topPosts),
		PostsOverview: utils.MapPostsOverview(postsOverview),
	}
	utils.SendSuccess(c, resp, "Dashboard retrieved successfully")
}

// GetPeriodComparisonHandler godoc
// @Router       /channels/{channelId}/analytics/comparison [get]
// @Security     BearerAuth
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

	prevStart, prevEnd := utils.PreviousPeriodInclusive(startDate, endDate)
	prevPeriod, err := repositories.GetAggregatedProfileAnalytics(c.Request.Context(), channelID, prevStart, prevEnd)
	if err != nil {
		resp := modelsv1.ComparisonResponse{
			Meta: modelsv1.DateRangeMeta{
				StartDate: utils.FormatAnalyticsDate(startDate),
				EndDate:   utils.FormatAnalyticsDate(endDate),
			},
			Current:  utils.MapAggregatedProfileOverview(currentPeriod),
			Previous: nil,
			Period: modelsv1.PeriodBounds{
				Start:     utils.FormatAnalyticsDate(startDate),
				End:       utils.FormatAnalyticsDate(endDate),
				PrevStart: utils.FormatAnalyticsDate(prevStart),
				PrevEnd:   utils.FormatAnalyticsDate(prevEnd),
			},
		}
		utils.SendSuccess(c, resp, "Comparison retrieved (no previous data)")
		return
	}

	cur := utils.MapAggregatedProfileOverview(currentPeriod)
	prev := utils.MapAggregatedProfileOverview(prevPeriod)
	resp := modelsv1.ComparisonResponse{
		Meta: modelsv1.DateRangeMeta{
			StartDate: utils.FormatAnalyticsDate(startDate),
			EndDate:   utils.FormatAnalyticsDate(endDate),
		},
		Current:  cur,
		Previous: &prev,
		Period: modelsv1.PeriodBounds{
			Start:     utils.FormatAnalyticsDate(startDate),
			End:       utils.FormatAnalyticsDate(endDate),
			PrevStart: utils.FormatAnalyticsDate(prevStart),
			PrevEnd:   utils.FormatAnalyticsDate(prevEnd),
		},
	}
	utils.SendSuccess(c, resp, "Comparison retrieved successfully")
}

// GetStoryAnalyticsHandler godoc
// @Summary      List story analytics
// @Param        post_id query int false "Filter by internal post id"
// @Param        limit query int false "Page size (default 100, max 500)"
// @Param        offset query int false "Offset"
// @Param        include_raw query string false "1 for raw JSON"
// @Router       /channels/{channelId}/analytics/stories [get]
// @Security     BearerAuth
func GetStoryAnalyticsHandler(c *gin.Context) {
	channelID, err := strconv.ParseInt(c.Param("channelId"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid channel ID")
		return
	}

	startDate, endDate := utils.ParseDateRange(c)
	limit, offset := utils.ParseAnalyticsPagination(c)
	if limit == 0 {
		limit = utils.AnalyticsDefaultPageSize
	}
	includeRaw := utils.AnalyticsIncludeRawQuery(c)

	var postID *int64
	if ps := strings.TrimSpace(c.Query("post_id")); ps != "" {
		if pid, err := strconv.ParseInt(ps, 10, 64); err == nil {
			postID = &pid
		}
	}

	total, err := repositories.CountStoryAnalyticsInRange(c.Request.Context(), channelID, startDate, endDate, postID)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to count story analytics: "+err.Error())
		return
	}

	rows, err := repositories.GetStoryAnalyticsJoined(c.Request.Context(), channelID, startDate, endDate, postID, limit, offset)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to fetch story analytics: "+err.Error())
		return
	}

	resp := modelsv1.StoryAnalyticsResponse{
		Meta: modelsv1.DateRangeMeta{
			StartDate: utils.FormatAnalyticsDate(startDate),
			EndDate:   utils.FormatAnalyticsDate(endDate),
		},
		Pagination: modelsv1.PaginationMeta{
			Limit:  limit,
			Offset: offset,
			Total:  total,
		},
		Stories: utils.MapStoryAnalyticsItems(rows, includeRaw),
	}
	utils.SendSuccess(c, resp, "Story analytics retrieved successfully")
}
