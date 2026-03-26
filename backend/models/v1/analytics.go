package modelsv1

import "encoding/json"

// PaginationMeta accompanies paginated analytics lists.
type PaginationMeta struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

// DateRangeMeta describes the requested analytics window (UTC calendar dates).
type DateRangeMeta struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// ProfileAnalyticsItem is one day of account-level metrics.
type ProfileAnalyticsItem struct {
	ID                  int64           `json:"id"`
	Date                string          `json:"date"`
	FollowerCount       *int            `json:"follower_count,omitempty"`
	Impressions         *int            `json:"impressions,omitempty"`
	ProfileViews        *int            `json:"profile_views,omitempty"`
	Reach               *int            `json:"reach,omitempty"`
	WebsiteClicks       *int            `json:"website_clicks,omitempty"`
	EmailClicks         *int            `json:"email_clicks,omitempty"`
	Views               *int            `json:"views,omitempty"`
	AccountsEngaged     *int            `json:"accounts_engaged,omitempty"`
	TotalInteractions   *int            `json:"total_interactions,omitempty"`
	BioLinkClicks       *int            `json:"bio_link_clicks,omitempty"`
	PhoneCallClicks     *int            `json:"phone_call_clicks,omitempty"`
	TextMessageClicks   *int            `json:"text_message_clicks,omitempty"`
	GetDirectionsClicks *int            `json:"get_directions_clicks,omitempty"`
	Raw                 json.RawMessage `json:"raw,omitempty"`
}

// ProfileAnalyticsResponse is GET .../analytics/profile.
type ProfileAnalyticsResponse struct {
	Meta       DateRangeMeta          `json:"meta"`
	Pagination PaginationMeta         `json:"pagination"`
	Analytics  []ProfileAnalyticsItem `json:"analytics"`
}

// PostAnalyticsItem is one snapshot row for a feed post/reel/carousel.
type PostAnalyticsItem struct {
	ID                int64           `json:"id"`
	PostID            int64           `json:"post_id"`
	Date              string          `json:"date"`
	Impressions       *int            `json:"impressions,omitempty"`
	Reach             *int            `json:"reach,omitempty"`
	Likes             *int            `json:"likes,omitempty"`
	Comments          *int            `json:"comments,omitempty"`
	Saves             *int            `json:"saves,omitempty"`
	Shares            *int            `json:"shares,omitempty"`
	VideoViews        *int            `json:"video_views,omitempty"`
	ProfileVisits     *int            `json:"profile_visits,omitempty"`
	Follows           *int            `json:"follows,omitempty"`
	Views             *int            `json:"views,omitempty"`
	TotalInteractions *int            `json:"total_interactions,omitempty"`
	EngagementRate    *float64        `json:"engagement_rate,omitempty"`
	Plays             *int            `json:"plays,omitempty"`
	Raw               json.RawMessage `json:"raw,omitempty"`
}

// PostAnalyticsResponse is GET .../analytics/posts.
type PostAnalyticsResponse struct {
	Meta       DateRangeMeta       `json:"meta"`
	Pagination PaginationMeta      `json:"pagination"`
	Analytics  []PostAnalyticsItem `json:"analytics"`
}

// AggregatedProfileOverview is account-level rollups for a range.
type AggregatedProfileOverview struct {
	TotalReach           int     `json:"total_reach"`
	TotalImpressions     int     `json:"total_impressions"`
	TotalProfileViews    int     `json:"total_profile_views"`
	TotalWebsiteClicks   int     `json:"total_website_clicks"`
	TotalViews           int     `json:"total_views"`
	TotalAccountsEngaged int     `json:"total_accounts_engaged"`
	TotalInteractions    int     `json:"total_interactions"`
	AverageReach         float64 `json:"average_reach"`
	FollowerGrowth       int     `json:"follower_growth"`
	StartFollowers       int     `json:"start_followers"`
	EndFollowers         int     `json:"end_followers"`
}

// OverviewResponse is GET .../analytics/overview.
type OverviewResponse struct {
	Meta     DateRangeMeta             `json:"meta"`
	Overview AggregatedProfileOverview `json:"overview"`
}

// TopPostItem is a ranked post with latest metrics in range.
type TopPostItem struct {
	PostID       int64   `json:"post_id"`
	PostType     string  `json:"post_type"`
	Caption      string  `json:"caption"`
	PublishedAt  *string `json:"published_at,omitempty"`
	ThumbnailURL *string `json:"thumbnail_url,omitempty"`
	Permalink    *string `json:"permalink,omitempty"`
	Impressions  int     `json:"impressions"`
	Reach       int     `json:"reach"`
	Likes       int     `json:"likes"`
	Comments    int     `json:"comments"`
	Saves       int     `json:"saves"`
	Shares      int     `json:"shares"`
	Plays       int     `json:"plays"`
	Engagement  int     `json:"engagement"`
}

// TopPostsResponse is GET .../analytics/top-posts.
type TopPostsResponse struct {
	Meta     DateRangeMeta `json:"meta"`
	Sort     string        `json:"sort"`
	Limit    int           `json:"limit"`
	TopPosts []TopPostItem `json:"top_posts"`
}

// PostsOverviewResponse wraps posts activity summary.
type PostsOverviewResponse struct {
	Meta     DateRangeMeta `json:"meta"`
	Overview PostsOverview `json:"overview"`
}

// PostsOverview mirrors repositories.PostsOverview with JSON tags.
type PostsOverview struct {
	NewPosts         int `json:"new_posts"`
	TotalPosts       int `json:"total_posts"`
	NewLikes         int `json:"new_likes"`
	NewComments      int `json:"new_comments"`
	NewSaves         int `json:"new_saves"`
	NewShares        int `json:"new_shares"`
	TotalLikes       int `json:"total_likes"`
	TotalComments    int `json:"total_comments"`
	TotalSaves       int `json:"total_saves"`
	TotalShares      int `json:"total_shares"`
	TotalReach       int `json:"total_reach"`
	TotalImpressions int `json:"total_impressions"`
}

// PostInsightsResponse is GET .../analytics/posts/:postId.
type PostInsightsResponse struct {
	Analytics []PostAnalyticsItem `json:"analytics"`
}

// PeriodBounds describes current and previous windows.
type PeriodBounds struct {
	Start     string `json:"start"`
	End       string `json:"end"`
	PrevStart string `json:"prev_start"`
	PrevEnd   string `json:"prev_end"`
}

// ComparisonResponse is GET .../analytics/comparison.
type ComparisonResponse struct {
	Meta     DateRangeMeta              `json:"meta"`
	Current  AggregatedProfileOverview  `json:"current"`
	Previous *AggregatedProfileOverview `json:"previous,omitempty"`
	Period   PeriodBounds               `json:"period"`
}

// DashboardResponse aggregates overview, top posts, and posts summary.
type DashboardResponse struct {
	Meta          DateRangeMeta             `json:"meta"`
	Overview      AggregatedProfileOverview `json:"overview"`
	TopPosts      []TopPostItem             `json:"top_posts"`
	PostsOverview PostsOverview             `json:"posts_overview"`
}

// StoryAnalyticsItem is one story metrics row with post context.
type StoryAnalyticsItem struct {
	ID           int64           `json:"id"`
	PostID       int64           `json:"post_id"`
	Date         string          `json:"date"`
	Impressions  *int            `json:"impressions,omitempty"`
	Reach        *int            `json:"reach,omitempty"`
	Exits        *int            `json:"exits,omitempty"`
	Replies      *int            `json:"replies,omitempty"`
	TapsForward  *int            `json:"taps_forward,omitempty"`
	TapsBackward *int            `json:"taps_backward,omitempty"`
	TapsExit     *int            `json:"taps_exit,omitempty"`
	Caption      string          `json:"caption"`
	PostType     string          `json:"post_type"`
	PublishedAt  *string         `json:"published_at,omitempty"`
	Raw          json.RawMessage `json:"raw,omitempty"`
}

// StoryAnalyticsResponse is GET .../analytics/stories.
type StoryAnalyticsResponse struct {
	Meta       DateRangeMeta        `json:"meta"`
	Pagination PaginationMeta       `json:"pagination"`
	Stories    []StoryAnalyticsItem `json:"stories"`
}
