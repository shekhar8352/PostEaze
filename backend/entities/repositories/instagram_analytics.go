package repositories

import (
	"context"

	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/utils/database"
)

// jsonbOrEmpty avoids pq sending an empty string for nil/empty []byte, which PostgreSQL
// rejects for JSONB ("invalid input syntax for type json").
func jsonbOrEmpty(b []byte) []byte {
	if len(b) == 0 {
		return []byte("{}")
	}
	return b
}

// UpsertInstagramPostAnalytics inserts or updates Instagram post analytics
func UpsertInstagramPostAnalytics(ctx context.Context, analytics *entities.InstagramPostAnalytics) error {
	db := database.GetDB()
	query := `
		INSERT INTO instagram_post_analytics 
		(channel_id, post_id, date, impressions, reach, likes, comments, saves, shares, video_views, profile_visits, follows, views, total_interactions, engagement_rate, plays, raw)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
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
			views = EXCLUDED.views,
			total_interactions = EXCLUDED.total_interactions,
			engagement_rate = EXCLUDED.engagement_rate,
			plays = EXCLUDED.plays,
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
		analytics.Views,
		analytics.TotalInteractions,
		analytics.EngagementRate,
		analytics.Plays,
		jsonbOrEmpty(analytics.Raw),
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
		jsonbOrEmpty(analytics.Raw),
	).Scan(&analytics.ID, &analytics.CreatedAt)
}

// UpsertInstagramProfileAnalytics inserts or updates Instagram profile analytics
func UpsertInstagramProfileAnalytics(ctx context.Context, analytics *entities.InstagramProfileAnalytics) error {
	db := database.GetDB()
	query := `
		INSERT INTO instagram_profile_analytics 
		(channel_id, date, follower_count, impressions, profile_views, reach, website_clicks, email_clicks, views, accounts_engaged, total_interactions, bio_link_clicks, phone_call_clicks, text_message_clicks, get_directions_clicks, raw)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		ON CONFLICT (channel_id, date) 
		DO UPDATE SET
			follower_count = COALESCE(EXCLUDED.follower_count, instagram_profile_analytics.follower_count),
			impressions = COALESCE(EXCLUDED.impressions, instagram_profile_analytics.impressions),
			profile_views = COALESCE(EXCLUDED.profile_views, instagram_profile_analytics.profile_views),
			reach = COALESCE(EXCLUDED.reach, instagram_profile_analytics.reach),
			website_clicks = COALESCE(EXCLUDED.website_clicks, instagram_profile_analytics.website_clicks),
			email_clicks = COALESCE(EXCLUDED.email_clicks, instagram_profile_analytics.email_clicks),
			views = COALESCE(EXCLUDED.views, instagram_profile_analytics.views),
			accounts_engaged = COALESCE(EXCLUDED.accounts_engaged, instagram_profile_analytics.accounts_engaged),
			total_interactions = COALESCE(EXCLUDED.total_interactions, instagram_profile_analytics.total_interactions),
			bio_link_clicks = COALESCE(EXCLUDED.bio_link_clicks, instagram_profile_analytics.bio_link_clicks),
			phone_call_clicks = COALESCE(EXCLUDED.phone_call_clicks, instagram_profile_analytics.phone_call_clicks),
			text_message_clicks = COALESCE(EXCLUDED.text_message_clicks, instagram_profile_analytics.text_message_clicks),
			get_directions_clicks = COALESCE(EXCLUDED.get_directions_clicks, instagram_profile_analytics.get_directions_clicks),
			raw = COALESCE(EXCLUDED.raw, instagram_profile_analytics.raw)
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
		analytics.Views,
		analytics.AccountsEngaged,
		analytics.TotalInteractions,
		analytics.BioLinkClicks,
		analytics.PhoneCallClicks,
		analytics.TextMessageClicks,
		analytics.GetDirectionsClicks,
		jsonbOrEmpty(analytics.Raw),
	).Scan(&analytics.ID, &analytics.CreatedAt)
}

// UpsertInstagramAudienceSnapshot stores or replaces the audience insight blob for a channel on a calendar day (UTC).
func UpsertInstagramAudienceSnapshot(ctx context.Context, s *entities.InstagramAudienceSnapshot) error {
	db := database.GetDB()
	query := `
		INSERT INTO instagram_audience_snapshots (channel_id, snapshot_date, raw)
		VALUES ($1, $2, $3)
		ON CONFLICT (channel_id, snapshot_date)
		DO UPDATE SET raw = EXCLUDED.raw
		RETURNING id, created_at
	`
	return db.QueryRowContext(ctx, query,
		s.ChannelID,
		s.SnapshotDate,
		jsonbOrEmpty(s.Raw),
	).Scan(&s.ID, &s.CreatedAt)
}
