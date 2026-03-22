package utils

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/entities/repositories"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
)

const (
	AnalyticsMaxPageSize     = 500
	AnalyticsDefaultPageSize = 100
)

// ParseAnalyticsPagination reads limit and offset from query (limit 0 = no cap; default 100; max 500).
func ParseAnalyticsPagination(c *gin.Context) (limit, offset int) {
	limit = AnalyticsDefaultPageSize
	offset = 0
	if ls := c.Query("limit"); ls != "" {
		if l, err := strconv.Atoi(ls); err == nil {
			if l == 0 {
				limit = 0
			} else if l > 0 {
				limit = l
				if limit > AnalyticsMaxPageSize {
					limit = AnalyticsMaxPageSize
				}
			}
		}
	}
	if os := c.Query("offset"); os != "" {
		if o, err := strconv.Atoi(os); err == nil && o >= 0 {
			offset = o
		}
	}
	return limit, offset
}

// AnalyticsIncludeRawQuery is true when the client asks for raw Meta JSON blobs.
func AnalyticsIncludeRawQuery(c *gin.Context) bool {
	v := strings.TrimSpace(c.Query("include_raw"))
	return v == "1" || strings.EqualFold(v, "true") || strings.EqualFold(v, "yes")
}

// OptionalPostTypeQuery returns post_type query or nil when empty.
func OptionalPostTypeQuery(c *gin.Context) *string {
	v := strings.TrimSpace(c.Query("post_type"))
	if v == "" {
		return nil
	}
	return &v
}

// FormatAnalyticsDate formats t as YYYY-MM-DD in UTC.
func FormatAnalyticsDate(t time.Time) string {
	return t.UTC().Format("2006-01-02")
}

// FormatAnalyticsTimePtr formats a timestamp as RFC3339 UTC, or nil.
func FormatAnalyticsTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.UTC().Format(time.RFC3339)
	return &s
}

func analyticsRawJSON(b []byte, include bool) json.RawMessage {
	if !include || len(b) == 0 {
		return nil
	}
	return json.RawMessage(b)
}

// MapProfileAnalyticsItems maps entity rows to API DTOs.
func MapProfileAnalyticsItems(items []entities.InstagramProfileAnalytics, includeRaw bool) []modelsv1.ProfileAnalyticsItem {
	out := make([]modelsv1.ProfileAnalyticsItem, 0, len(items))
	for _, a := range items {
		out = append(out, modelsv1.ProfileAnalyticsItem{
			ID:                  a.ID,
			Date:                FormatAnalyticsDate(a.Date),
			FollowerCount:       a.FollowerCount,
			Impressions:         a.Impressions,
			ProfileViews:        a.ProfileViews,
			Reach:               a.Reach,
			WebsiteClicks:       a.WebsiteClicks,
			EmailClicks:         a.EmailClicks,
			Views:               a.Views,
			AccountsEngaged:     a.AccountsEngaged,
			TotalInteractions:   a.TotalInteractions,
			BioLinkClicks:       a.BioLinkClicks,
			PhoneCallClicks:     a.PhoneCallClicks,
			TextMessageClicks:   a.TextMessageClicks,
			GetDirectionsClicks: a.GetDirectionsClicks,
			Raw:                 analyticsRawJSON(a.Raw, includeRaw),
		})
	}
	return out
}

// MapPostAnalyticsItems maps entity rows to API DTOs.
func MapPostAnalyticsItems(items []entities.InstagramPostAnalytics, includeRaw bool) []modelsv1.PostAnalyticsItem {
	out := make([]modelsv1.PostAnalyticsItem, 0, len(items))
	for _, a := range items {
		out = append(out, modelsv1.PostAnalyticsItem{
			ID:                a.ID,
			PostID:            a.PostID,
			Date:              FormatAnalyticsDate(a.Date),
			Impressions:       a.Impressions,
			Reach:             a.Reach,
			Likes:             a.Likes,
			Comments:          a.Comments,
			Saves:             a.Saves,
			Shares:            a.Shares,
			VideoViews:        a.VideoViews,
			ProfileVisits:     a.ProfileVisits,
			Follows:           a.Follows,
			Views:             a.Views,
			TotalInteractions: a.TotalInteractions,
			EngagementRate:    a.EngagementRate,
			Plays:             a.Plays,
			Raw:               analyticsRawJSON(a.Raw, includeRaw),
		})
	}
	return out
}

// MapAggregatedProfileOverview maps repository aggregate to API DTO.
func MapAggregatedProfileOverview(a *repositories.AggregatedProfileAnalytics) modelsv1.AggregatedProfileOverview {
	if a == nil {
		return modelsv1.AggregatedProfileOverview{}
	}
	return modelsv1.AggregatedProfileOverview{
		TotalReach:           a.TotalReach,
		TotalImpressions:     a.TotalImpressions,
		TotalProfileViews:    a.TotalProfileViews,
		TotalWebsiteClicks:   a.TotalWebsiteClicks,
		TotalViews:           a.TotalViews,
		TotalAccountsEngaged: a.TotalAccountsEngaged,
		TotalInteractions:    a.TotalInteractions,
		AverageReach:         a.AverageReach,
		FollowerGrowth:       a.FollowerGrowth,
		StartFollowers:       a.StartFollowers,
		EndFollowers:         a.EndFollowers,
	}
}

// MapTopPostItems maps repository top posts to API DTOs.
func MapTopPostItems(items []repositories.TopPost) []modelsv1.TopPostItem {
	out := make([]modelsv1.TopPostItem, 0, len(items))
	for _, tp := range items {
		out = append(out, modelsv1.TopPostItem{
			PostID:      tp.PostID,
			PostType:    tp.PostType,
			Caption:     tp.Caption,
			PublishedAt: FormatAnalyticsTimePtr(tp.PublishedAt),
			Impressions: tp.Impressions,
			Reach:       tp.Reach,
			Likes:       tp.Likes,
			Comments:    tp.Comments,
			Saves:       tp.Saves,
			Shares:      tp.Shares,
			Plays:       tp.Plays,
			Engagement:  tp.Engagement,
		})
	}
	return out
}

// MapPostsOverview maps repository overview to API DTO.
func MapPostsOverview(o *repositories.PostsOverview) modelsv1.PostsOverview {
	if o == nil {
		return modelsv1.PostsOverview{}
	}
	return modelsv1.PostsOverview{
		NewPosts:         o.NewPosts,
		TotalPosts:       o.TotalPosts,
		NewLikes:         o.NewLikes,
		NewComments:      o.NewComments,
		NewSaves:         o.NewSaves,
		NewShares:        o.NewShares,
		TotalLikes:       o.TotalLikes,
		TotalComments:    o.TotalComments,
		TotalSaves:       o.TotalSaves,
		TotalShares:      o.TotalShares,
		TotalReach:       o.TotalReach,
		TotalImpressions: o.TotalImpressions,
	}
}

// ParseTopPostsSort reads sort query (engagement|reach|impressions|plays).
func ParseTopPostsSort(c *gin.Context) repositories.TopPostsSort {
	s := strings.ToLower(strings.TrimSpace(c.DefaultQuery("sort", string(repositories.TopPostsSortEngagement))))
	switch s {
	case string(repositories.TopPostsSortReach):
		return repositories.TopPostsSortReach
	case string(repositories.TopPostsSortImpressions):
		return repositories.TopPostsSortImpressions
	case string(repositories.TopPostsSortPlays):
		return repositories.TopPostsSortPlays
	default:
		return repositories.TopPostsSortEngagement
	}
}

// MapStoryAnalyticsItems maps joined story rows to API DTOs.
func MapStoryAnalyticsItems(rows []repositories.StoryAnalyticsJoinedRow, includeRaw bool) []modelsv1.StoryAnalyticsItem {
	out := make([]modelsv1.StoryAnalyticsItem, 0, len(rows))
	for _, r := range rows {
		s := r.InstagramStoryAnalytics
		out = append(out, modelsv1.StoryAnalyticsItem{
			ID:           s.ID,
			PostID:       s.PostID,
			Date:         FormatAnalyticsDate(s.Date),
			Impressions:  s.Impressions,
			Reach:        s.Reach,
			Exits:        s.Exits,
			Replies:      s.Replies,
			TapsForward:  s.TapsForward,
			TapsBackward: s.TapsBackward,
			TapsExit:     s.TapsExit,
			Caption:      r.Caption,
			PostType:     r.PostType,
			PublishedAt:  FormatAnalyticsTimePtr(r.PublishedAt),
			Raw:          analyticsRawJSON(s.Raw, includeRaw),
		})
	}
	return out
}
