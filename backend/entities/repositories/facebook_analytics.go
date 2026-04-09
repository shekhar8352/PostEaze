package repositories

import (
	"context"

	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/utils/database"
)

// UpsertFacebookPageAnalytics inserts or updates a daily Facebook page analytics row.
func UpsertFacebookPageAnalytics(ctx context.Context, a *entities.FacebookPageAnalytics) error {
	db := database.GetDB()
	query := `
		INSERT INTO facebook_page_analytics (
			channel_id, date, followers, likes, reach, impressions, page_engaged_users, raw
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (channel_id, date) DO UPDATE SET
			followers = EXCLUDED.followers,
			likes = EXCLUDED.likes,
			reach = EXCLUDED.reach,
			impressions = EXCLUDED.impressions,
			page_engaged_users = EXCLUDED.page_engaged_users,
			raw = EXCLUDED.raw
		RETURNING id, created_at
	`
	return db.QueryRowContext(ctx, query,
		a.ChannelID,
		a.Date,
		a.Followers,
		a.Likes,
		a.Reach,
		a.Impressions,
		a.PageEngagedUsers,
		jsonbOrEmpty(a.Raw),
	).Scan(&a.ID, &a.CreatedAt)
}

// UpsertFacebookPostAnalytics inserts or updates Facebook post analytics for a calendar day.
func UpsertFacebookPostAnalytics(ctx context.Context, a *entities.FacebookPostAnalytics) error {
	db := database.GetDB()
	query := `
		INSERT INTO facebook_post_analytics (
			channel_id, post_id, date, impressions, reach, likes, comments, shares, photo_views, video_views, raw
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (post_id, date) DO UPDATE SET
			impressions = EXCLUDED.impressions,
			reach = EXCLUDED.reach,
			likes = EXCLUDED.likes,
			comments = EXCLUDED.comments,
			shares = EXCLUDED.shares,
			photo_views = EXCLUDED.photo_views,
			video_views = EXCLUDED.video_views,
			raw = EXCLUDED.raw
		RETURNING id, created_at
	`
	return db.QueryRowContext(ctx, query,
		a.ChannelID,
		a.PostID,
		a.Date,
		a.Impressions,
		a.Reach,
		a.Likes,
		a.Comments,
		a.Shares,
		a.PhotoViews,
		a.VideoViews,
		jsonbOrEmpty(a.Raw),
	).Scan(&a.ID, &a.CreatedAt)
}
