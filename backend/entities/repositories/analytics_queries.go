package repositories

import (
	"context"
	"time"

	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/utils/database"
)

// GetProfileAnalyticsByDateRange retrieves profile analytics for a channel within a date range
func GetProfileAnalyticsByDateRange(ctx context.Context, channelID int64, startDate, endDate time.Time) ([]entities.InstagramProfileAnalytics, error) {
	db := database.GetDB()
	query := `
		SELECT id, channel_id, date, follower_count, impressions, profile_views, reach, website_clicks, email_clicks, raw, created_at
		FROM instagram_profile_analytics
		WHERE channel_id = $1 AND date >= $2 AND date <= $3
		ORDER BY date DESC
	`

	rows, err := db.QueryContext(ctx, query, channelID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analytics []entities.InstagramProfileAnalytics
	for rows.Next() {
		var a entities.InstagramProfileAnalytics
		err := rows.Scan(
			&a.ID,
			&a.ChannelID,
			&a.Date,
			&a.FollowerCount,
			&a.Impressions,
			&a.ProfileViews,
			&a.Reach,
			&a.WebsiteClicks,
			&a.EmailClicks,
			&a.Raw,
			&a.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		analytics = append(analytics, a)
	}

	return analytics, rows.Err()
}

// GetPostAnalyticsByDateRange retrieves post analytics for a channel within a date range
func GetPostAnalyticsByDateRange(ctx context.Context, channelID int64, startDate, endDate time.Time) ([]entities.InstagramPostAnalytics, error) {
	db := database.GetDB()
	query := `
		SELECT pa.id, pa.channel_id, pa.post_id, pa.date, pa.impressions, pa.reach, pa.likes, pa.comments, 
		       pa.saves, pa.shares, pa.video_views, pa.profile_visits, pa.follows, pa.raw, pa.created_at
		FROM instagram_post_analytics pa
		WHERE pa.channel_id = $1 AND pa.date >= $2 AND pa.date <= $3
		ORDER BY pa.date DESC
	`

	rows, err := db.QueryContext(ctx, query, channelID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analytics []entities.InstagramPostAnalytics
	for rows.Next() {
		var a entities.InstagramPostAnalytics
		err := rows.Scan(
			&a.ID,
			&a.ChannelID,
			&a.PostID,
			&a.Date,
			&a.Impressions,
			&a.Reach,
			&a.Likes,
			&a.Comments,
			&a.Saves,
			&a.Shares,
			&a.VideoViews,
			&a.ProfileVisits,
			&a.Follows,
			&a.Raw,
			&a.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		analytics = append(analytics, a)
	}

	return analytics, rows.Err()
}

// AggregatedProfileAnalytics represents aggregated profile metrics
type AggregatedProfileAnalytics struct {
	TotalReach         int
	TotalImpressions   int
	TotalProfileViews  int
	TotalWebsiteClicks int
	AverageReach       float64
	FollowerGrowth     int
	StartFollowers     int
	EndFollowers       int
}

// GetAggregatedProfileAnalytics returns aggregated profile analytics for a date range
func GetAggregatedProfileAnalytics(ctx context.Context, channelID int64, startDate, endDate time.Time) (*AggregatedProfileAnalytics, error) {
	db := database.GetDB()
	query := `
		SELECT 
			COALESCE(SUM(reach), 0) as total_reach,
			COALESCE(SUM(impressions), 0) as total_impressions,
			COALESCE(SUM(profile_views), 0) as total_profile_views,
			COALESCE(SUM(website_clicks), 0) as total_website_clicks,
			COALESCE(AVG(reach), 0) as average_reach,
			(SELECT follower_count FROM instagram_profile_analytics WHERE channel_id = $1 AND date >= $2 ORDER BY date ASC LIMIT 1) as start_followers,
			(SELECT follower_count FROM instagram_profile_analytics WHERE channel_id = $1 AND date <= $3 ORDER BY date DESC LIMIT 1) as end_followers
		FROM instagram_profile_analytics
		WHERE channel_id = $1 AND date >= $2 AND date <= $3
	`

	var agg AggregatedProfileAnalytics
	var startFollowers, endFollowers *int

	err := db.QueryRowContext(ctx, query, channelID, startDate, endDate).Scan(
		&agg.TotalReach,
		&agg.TotalImpressions,
		&agg.TotalProfileViews,
		&agg.TotalWebsiteClicks,
		&agg.AverageReach,
		&startFollowers,
		&endFollowers,
	)

	if err != nil {
		return nil, err
	}

	if startFollowers != nil {
		agg.StartFollowers = *startFollowers
	}
	if endFollowers != nil {
		agg.EndFollowers = *endFollowers
	}
	agg.FollowerGrowth = agg.EndFollowers - agg.StartFollowers

	return &agg, nil
}

// TopPost represents a post with its analytics
type TopPost struct {
	PostID      int64
	PostType    string
	Caption     string
	PublishedAt time.Time
	Impressions int
	Reach       int
	Likes       int
	Comments    int
	Saves       int
	Shares      int
	Engagement  int
}

// GetTopPosts returns top performing posts by engagement
func GetTopPosts(ctx context.Context, channelID int64, limit int, startDate, endDate time.Time) ([]TopPost, error) {
	db := database.GetDB()
	query := `
		SELECT 
			p.id,
			COALESCE(p.post_type, 'post') as post_type,
			COALESCE(p.caption, '') as caption,
			p.published_at,
			COALESCE(pa.impressions, 0) as impressions,
			COALESCE(pa.reach, 0) as reach,
			COALESCE(pa.likes, 0) as likes,
			COALESCE(pa.comments, 0) as comments,
			COALESCE(pa.saves, 0) as saves,
			COALESCE(pa.shares, 0) as shares,
			COALESCE(pa.likes, 0) + COALESCE(pa.comments, 0) + COALESCE(pa.saves, 0) + COALESCE(pa.shares, 0) as engagement
		FROM posts p
		LEFT JOIN instagram_post_analytics pa ON p.id = pa.post_id
		WHERE p.channel_id = $1 
			AND p.published_at >= $2 
			AND p.published_at <= $3
		ORDER BY engagement DESC
		LIMIT $4
	`

	rows, err := db.QueryContext(ctx, query, channelID, startDate, endDate, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topPosts []TopPost
	for rows.Next() {
		var tp TopPost
		err := rows.Scan(
			&tp.PostID,
			&tp.PostType,
			&tp.Caption,
			&tp.PublishedAt,
			&tp.Impressions,
			&tp.Reach,
			&tp.Likes,
			&tp.Comments,
			&tp.Saves,
			&tp.Shares,
			&tp.Engagement,
		)
		if err != nil {
			return nil, err
		}
		topPosts = append(topPosts, tp)
	}

	return topPosts, rows.Err()
}

// PostsOverview represents aggregated posts activity metrics
type PostsOverview struct {
	NewPosts         int
	TotalPosts       int
	NewLikes         int
	NewComments      int
	NewSaves         int
	NewShares        int
	TotalLikes       int
	TotalComments    int
	TotalSaves       int
	TotalShares      int
	TotalReach       int
	TotalImpressions int
}

// GetPostsOverview returns aggregated posts activity for a date range
func GetPostsOverview(ctx context.Context, channelID int64, startDate, endDate time.Time) (*PostsOverview, error) {
	db := database.GetDB()

	// Get new posts count
	newPostsQuery := `
		SELECT COUNT(*) 
		FROM posts 
		WHERE channel_id = $1 AND published_at >= $2 AND published_at <= $3
	`

	// Get total posts count
	totalPostsQuery := `
		SELECT COUNT(*) 
		FROM posts 
		WHERE channel_id = $1
	`

	// Get aggregated analytics for the date range from instagram_post_analytics
	// We need to get the LATEST analytics for each post in the date range
	analyticsQuery := `
		SELECT 
			COALESCE(SUM(COALESCE(ipa.likes, 0)), 0) as total_likes,
			COALESCE(SUM(COALESCE(ipa.comments, 0)), 0) as total_comments,
			COALESCE(SUM(COALESCE(ipa.saves, 0)), 0) as total_saves,
			COALESCE(SUM(COALESCE(ipa.shares, 0)), 0) as total_shares,
			COALESCE(SUM(COALESCE(ipa.reach, 0)), 0) as total_reach,
			COALESCE(SUM(COALESCE(ipa.impressions, 0)), 0) as total_impressions
		FROM (
			SELECT DISTINCT ON (post_id) 
				post_id, likes, comments, saves, shares, reach, impressions
			FROM instagram_post_analytics
			WHERE channel_id = $1 AND date >= $2 AND date <= $3
			ORDER BY post_id, date DESC
		) ipa
	`

	// Get previous period analytics for comparison
	duration := endDate.Sub(startDate)
	prevStartDate := startDate.Add(-duration - 24*time.Hour)
	prevEndDate := startDate.Add(-time.Second)

	prevAnalyticsQuery := `
		SELECT 
			COALESCE(SUM(COALESCE(ipa.likes, 0)), 0) as prev_likes,
			COALESCE(SUM(COALESCE(ipa.comments, 0)), 0) as prev_comments,
			COALESCE(SUM(COALESCE(ipa.saves, 0)), 0) as prev_saves,
			COALESCE(SUM(COALESCE(ipa.shares, 0)), 0) as prev_shares
		FROM (
			SELECT DISTINCT ON (post_id) 
				post_id, likes, comments, saves, shares
			FROM instagram_post_analytics
			WHERE channel_id = $1 AND date >= $2 AND date <= $3
			ORDER BY post_id, date DESC
		) ipa
	`

	var overview PostsOverview
	var prevLikes, prevComments, prevSaves, prevShares int

	// Execute queries
	err := db.QueryRowContext(ctx, newPostsQuery, channelID, startDate, endDate).Scan(&overview.NewPosts)
	if err != nil {
		return nil, err
	}

	err = db.QueryRowContext(ctx, totalPostsQuery, channelID).Scan(&overview.TotalPosts)
	if err != nil {
		return nil, err
	}

	err = db.QueryRowContext(ctx, analyticsQuery, channelID, startDate, endDate).Scan(
		&overview.TotalLikes,
		&overview.TotalComments,
		&overview.TotalSaves,
		&overview.TotalShares,
		&overview.TotalReach,
		&overview.TotalImpressions,
	)
	if err != nil {
		return nil, err
	}

	err = db.QueryRowContext(ctx, prevAnalyticsQuery, channelID, prevStartDate, prevEndDate).Scan(
		&prevLikes,
		&prevComments,
		&prevSaves,
		&prevShares,
	)
	if err != nil {
		return nil, err
	}

	// Calculate new engagement (difference from previous period)
	overview.NewLikes = overview.TotalLikes - prevLikes
	overview.NewComments = overview.TotalComments - prevComments
	overview.NewSaves = overview.TotalSaves - prevSaves
	overview.NewShares = overview.TotalShares - prevShares

	return &overview, nil
}
