package repositories

import (
	"context"

	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/utils/database"
)

// UpsertInstagramPostAnalytics inserts or updates Instagram post analytics
func UpsertInstagramPostAnalytics(ctx context.Context, analytics *entities.InstagramPostAnalytics) error {
	db := database.GetDB()
	query := `
		INSERT INTO instagram_post_analytics 
		(channel_id, post_id, date, impressions, reach, likes, comments, saves, shares, video_views, profile_visits, follows, raw)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (post_id, date) 
		DO UPDATE SET
			impressions = EXCLUDED.impressions,
			reach = EXCLUDED.reach,
			likes = EXCLUDED.likes,
			comments = EXCLUDED.comments,
			saves = EXCLUDED.saves,
			shares = EXCLUDED.shares,
			video_views = EXCLUDED.video_views,
			profile_visits = EXCLUDED.profile_visits,
			follows = EXCLUDED.follows,
			raw = EXCLUDED.raw
		RETURNING id, created_at
	`
	return db.QueryRowContext(ctx, query,
		analytics.ChannelID,
		analytics.PostID,
		analytics.Date,
		analytics.Impressions,
		analytics.Reach,
		analytics.Likes,
		analytics.Comments,
		analytics.Saves,
		analytics.Shares,
		analytics.VideoViews,
		analytics.ProfileVisits,
		analytics.Follows,
		analytics.Raw,
	).Scan(&analytics.ID, &analytics.CreatedAt)
}

// UpsertInstagramStoryAnalytics inserts or updates Instagram story analytics
func UpsertInstagramStoryAnalytics(ctx context.Context, analytics *entities.InstagramStoryAnalytics) error {
	db := database.GetDB()
	query := `
		INSERT INTO instagram_story_analytics 
		(channel_id, post_id, date, impressions, reach, exits, forwards, replies, taps_forward, taps_backward, taps_exit, raw)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (post_id, date) 
		DO UPDATE SET
			impressions = EXCLUDED.impressions,
			reach = EXCLUDED.reach,
			exits = EXCLUDED.exits,
			forwards = EXCLUDED.forwards,
			replies = EXCLUDED.replies,
			taps_forward = EXCLUDED.taps_forward,
			taps_backward = EXCLUDED.taps_backward,
			taps_exit = EXCLUDED.taps_exit,
			raw = EXCLUDED.raw
		RETURNING id, created_at
	`
	return db.QueryRowContext(ctx, query,
		analytics.ChannelID,
		analytics.PostID,
		analytics.Date,
		analytics.Impressions,
		analytics.Reach,
		analytics.Exits,
		analytics.Forwards,
		analytics.Replies,
		analytics.TapsForward,
		analytics.TapsBackward,
		analytics.TapsExit,
		analytics.Raw,
	).Scan(&analytics.ID, &analytics.CreatedAt)
}

// UpsertInstagramProfileAnalytics inserts or updates Instagram profile analytics
func UpsertInstagramProfileAnalytics(ctx context.Context, analytics *entities.InstagramProfileAnalytics) error {
	db := database.GetDB()
	query := `
		INSERT INTO instagram_profile_analytics 
		(channel_id, date, follower_count, impressions, profile_views, reach, website_clicks, email_clicks, raw)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (channel_id, date) 
		DO UPDATE SET
			follower_count = EXCLUDED.follower_count,
			impressions = EXCLUDED.impressions,
			profile_views = EXCLUDED.profile_views,
			reach = EXCLUDED.reach,
			website_clicks = EXCLUDED.website_clicks,
			email_clicks = EXCLUDED.email_clicks,
			raw = EXCLUDED.raw
		RETURNING id, created_at
	`
	return db.QueryRowContext(ctx, query,
		analytics.ChannelID,
		analytics.Date,
		analytics.FollowerCount,
		analytics.Impressions,
		analytics.ProfileViews,
		analytics.Reach,
		analytics.WebsiteClicks,
		analytics.EmailClicks,
		analytics.Raw,
	).Scan(&analytics.ID, &analytics.CreatedAt)
}
