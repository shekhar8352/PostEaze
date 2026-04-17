package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/utils/database"
)

// ProfileAnalyticsListFilters optional pagination for profile series.
type ProfileAnalyticsListFilters struct {
	Limit  int
	Offset int
}

// PostAnalyticsListFilters optional post_type filter and pagination.
type PostAnalyticsListFilters struct {
	PostType *string
	Limit    int
	Offset   int
}

// CountProfileAnalyticsInRange returns total rows matching the date filter.
func CountProfileAnalyticsInRange(ctx context.Context, channelID int64, startDate, endDate time.Time) (int, error) {
	db := database.GetDB()
	var n int
	err := db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM instagram_profile_analytics
		WHERE channel_id = $1 AND date >= $2::date AND date <= $3::date
	`, channelID, startDate, endDate).Scan(&n)
	return n, err
}

// GetProfileAnalyticsByDateRange retrieves profile analytics for a channel within a date range.
func GetProfileAnalyticsByDateRange(ctx context.Context, channelID int64, startDate, endDate time.Time, f ProfileAnalyticsListFilters) ([]entities.InstagramProfileAnalytics, error) {
	db := database.GetDB()
	query := `
		SELECT id, channel_id, date, follower_count, impressions, profile_views, reach, website_clicks, email_clicks,
		       views, accounts_engaged, total_interactions, bio_link_clicks, phone_call_clicks, text_message_clicks,
		       get_directions_clicks, raw, created_at
		FROM instagram_profile_analytics
		WHERE channel_id = $1 AND date >= $2::date AND date <= $3::date
		ORDER BY date DESC
	`
	args := []interface{}{channelID, startDate, endDate}
	if f.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
		args = append(args, f.Limit, f.Offset)
	}

	rows, err := db.QueryContext(ctx, query, args...)
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
			&a.Views,
			&a.AccountsEngaged,
			&a.TotalInteractions,
			&a.BioLinkClicks,
			&a.PhoneCallClicks,
			&a.TextMessageClicks,
			&a.GetDirectionsClicks,
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

// CountPostAnalyticsInRange counts post analytics rows after optional post_type filter.
func CountPostAnalyticsInRange(ctx context.Context, channelID int64, startDate, endDate time.Time, postType *string) (int, error) {
	db := database.GetDB()
	base := `
		SELECT COUNT(*) FROM instagram_post_analytics pa
		JOIN posts p ON p.id = pa.post_id AND $1 = ANY(p.channel_ids)
		WHERE pa.channel_id = $1 AND pa.date >= $2::date AND pa.date <= $3::date
	`
	args := []interface{}{channelID, startDate, endDate}
	if postType != nil && *postType != "" {
		base += ` AND COALESCE(p.post_type, '') = $4`
		args = append(args, *postType)
	}
	var n int
	err := db.QueryRowContext(ctx, base, args...).Scan(&n)
	return n, err
}

// GetPostAnalyticsByDateRange retrieves post analytics for a channel within a date range.
func GetPostAnalyticsByDateRange(ctx context.Context, channelID int64, startDate, endDate time.Time, f PostAnalyticsListFilters) ([]entities.InstagramPostAnalytics, error) {
	db := database.GetDB()
	query := `
		SELECT pa.id, pa.channel_id, pa.post_id, pa.date, pa.impressions, pa.reach, pa.likes, pa.comments,
		       pa.saves, pa.shares, pa.video_views, pa.profile_visits, pa.follows, pa.views, pa.total_interactions,
		       pa.engagement_rate, pa.plays, pa.raw, pa.created_at
		FROM instagram_post_analytics pa
		JOIN posts p ON p.id = pa.post_id AND $1 = ANY(p.channel_ids)
		WHERE pa.channel_id = $1 AND pa.date >= $2::date AND pa.date <= $3::date
	`
	args := []interface{}{channelID, startDate, endDate}
	argN := 4
	if f.PostType != nil && *f.PostType != "" {
		query += fmt.Sprintf(" AND COALESCE(p.post_type, '') = $%d", argN)
		args = append(args, *f.PostType)
		argN++
	}
	query += " ORDER BY pa.date DESC, pa.id DESC"
	if f.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argN, argN+1)
		args = append(args, f.Limit, f.Offset)
	}

	rows, err := db.QueryContext(ctx, query, args...)
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
			&a.Views,
			&a.TotalInteractions,
			&a.EngagementRate,
			&a.Plays,
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
	TotalReach           int
	TotalImpressions     int
	TotalProfileViews    int
	TotalWebsiteClicks   int
	TotalViews           int
	TotalAccountsEngaged int
	TotalInteractions    int
	AverageReach         float64
	FollowerGrowth       int
	StartFollowers       int
	EndFollowers         int
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
			COALESCE(SUM(views), 0) as total_views,
			COALESCE(SUM(accounts_engaged), 0) as total_accounts_engaged,
			COALESCE(SUM(total_interactions), 0) as total_interactions,
			COALESCE(AVG(reach), 0) as average_reach,
			(SELECT follower_count FROM instagram_profile_analytics WHERE channel_id = $1 AND date >= $2 ORDER BY date ASC LIMIT 1) as start_followers,
			(SELECT follower_count FROM instagram_profile_analytics WHERE channel_id = $1 AND date <= $3 ORDER BY date DESC LIMIT 1) as end_followers
		FROM instagram_profile_analytics
		WHERE channel_id = $1 AND date >= $2::date AND date <= $3::date
	`

	var agg AggregatedProfileAnalytics
	var startFollowers, endFollowers *int

	err := db.QueryRowContext(ctx, query, channelID, startDate, endDate).Scan(
		&agg.TotalReach,
		&agg.TotalImpressions,
		&agg.TotalProfileViews,
		&agg.TotalWebsiteClicks,
		&agg.TotalViews,
		&agg.TotalAccountsEngaged,
		&agg.TotalInteractions,
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

// GetPostDetailedAnalytics returns historical analytics for a single post
func GetPostDetailedAnalytics(ctx context.Context, postID int64, limit int) ([]entities.InstagramPostAnalytics, error) {
	db := database.GetDB()
	query := `
		SELECT id, channel_id, post_id, date, impressions, reach, likes, comments, saves, shares, video_views, profile_visits, follows, views, total_interactions, engagement_rate, plays, raw, created_at
		FROM instagram_post_analytics
		WHERE post_id = $1
		ORDER BY date DESC
		LIMIT $2
	`
	rows, err := db.QueryContext(ctx, query, postID, limit)
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
			&a.Views,
			&a.TotalInteractions,
			&a.EngagementRate,
			&a.Plays,
			&a.Raw,
			&a.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		analytics = append(analytics, a)
	}
	return analytics, nil
}

// TopPost represents a post with its analytics
type TopPost struct {
	PostID       int64
	PostType     string
	Caption      string
	PublishedAt  *time.Time
	ThumbnailURL *string
	Permalink    *string
	Impressions  int
	Reach        int
	Likes        int
	Comments     int
	Saves        int
	Shares       int
	Plays        int
	Engagement   int
}

// TopPostsSort controls ranking for GetTopPosts (whitelist).
type TopPostsSort string

const (
	TopPostsSortEngagement  TopPostsSort = "engagement"
	TopPostsSortReach       TopPostsSort = "reach"
	TopPostsSortImpressions TopPostsSort = "impressions"
	TopPostsSortPlays       TopPostsSort = "plays"
)

// GetTopPosts returns top performing posts using the latest analytics row per post in the date window.
func GetTopPosts(ctx context.Context, channelID int64, limit int, startDate, endDate time.Time, sort TopPostsSort, postType *string) ([]TopPost, error) {
	orderClause := "engagement DESC NULLS LAST"
	switch sort {
	case TopPostsSortReach:
		orderClause = "reach DESC NULLS LAST, engagement DESC NULLS LAST"
	case TopPostsSortImpressions:
		orderClause = "impressions DESC NULLS LAST, engagement DESC NULLS LAST"
	case TopPostsSortPlays:
		orderClause = "plays DESC NULLS LAST, engagement DESC NULLS LAST"
	case TopPostsSortEngagement, "":
		orderClause = "engagement DESC NULLS LAST"
	}

	db := database.GetDB()
	query := fmt.Sprintf(`
		WITH latest AS (
			SELECT DISTINCT ON (pa.post_id)
				pa.post_id,
				pa.impressions,
				pa.reach,
				pa.likes,
				pa.comments,
				pa.saves,
				pa.shares,
				pa.plays
			FROM instagram_post_analytics pa
			WHERE pa.channel_id = $1 AND pa.date >= $2::date AND pa.date <= $3::date
			ORDER BY pa.post_id, pa.date DESC
		)
		SELECT
			p.id,
			COALESCE(p.post_type, 'post') AS post_type,
			COALESCE(p.caption, '') AS caption,
			p.published_at,
			NULLIF(TRIM(COALESCE(p.media->0->>'url', p.media->0->>'media_url', '')), '') AS thumbnail_url,
			NULLIF(TRIM(COALESCE(p.media->0->>'permalink', '')), '') AS permalink,
			COALESCE(latest.impressions, 0) AS impressions,
			COALESCE(latest.reach, 0) AS reach,
			COALESCE(latest.likes, 0) AS likes,
			COALESCE(latest.comments, 0) AS comments,
			COALESCE(latest.saves, 0) AS saves,
			COALESCE(latest.shares, 0) AS shares,
			COALESCE(latest.plays, 0) AS plays,
			COALESCE(latest.likes, 0) + COALESCE(latest.comments, 0) + COALESCE(latest.saves, 0) + COALESCE(latest.shares, 0) AS engagement
		FROM posts p
		INNER JOIN latest ON latest.post_id = p.id
		WHERE $1 = ANY(p.channel_ids)
			AND p.published_at IS NOT NULL
			AND (p.published_at AT TIME ZONE 'UTC')::date >= $2::date
			AND (p.published_at AT TIME ZONE 'UTC')::date <= $3::date
			AND ($4::text IS NULL OR $4::text = '' OR COALESCE(p.post_type, '') = $4)
		ORDER BY %s
		LIMIT $5
	`, orderClause)

	rows, err := db.QueryContext(ctx, query, channelID, startDate, endDate, nullableString(postType), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topPosts []TopPost
	for rows.Next() {
		var tp TopPost
		var pubAt sql.NullTime
		var thumb, link sql.NullString
		err := rows.Scan(
			&tp.PostID,
			&tp.PostType,
			&tp.Caption,
			&pubAt,
			&thumb,
			&link,
			&tp.Impressions,
			&tp.Reach,
			&tp.Likes,
			&tp.Comments,
			&tp.Saves,
			&tp.Shares,
			&tp.Plays,
			&tp.Engagement,
		)
		if err != nil {
			return nil, err
		}
		if pubAt.Valid {
			t := pubAt.Time
			tp.PublishedAt = &t
		}
		if thumb.Valid && thumb.String != "" {
			s := thumb.String
			tp.ThumbnailURL = &s
		}
		if link.Valid && link.String != "" {
			s := link.String
			tp.Permalink = &s
		}
		topPosts = append(topPosts, tp)
	}

	return topPosts, rows.Err()
}

// DailyPostEngagementRow is summed day-over-day deltas for post-level metrics (UTC calendar day).
type DailyPostEngagementRow struct {
	Date     time.Time
	Likes    int64
	Comments int64
	Shares   int64
	Saves    int64
}

// GetDailyPostEngagementSeries returns per-day sums of (metric_today - metric_yesterday) per post,
// using all snapshots up to endDate so the first in-range day has a meaningful prior row when it exists.
// When there is no prior row for a post, that post contributes 0 for that day (avoids counting lifetime totals as “daily”).
func GetDailyPostEngagementSeries(ctx context.Context, channelID int64, startDate, endDate time.Time) ([]DailyPostEngagementRow, error) {
	db := database.GetDB()
	query := `
		WITH ordered AS (
			SELECT
				pa.post_id,
				pa.date::date AS d,
				COALESCE(pa.likes, 0)::bigint AS likes,
				COALESCE(pa.comments, 0)::bigint AS comments,
				COALESCE(pa.shares, 0)::bigint AS shares,
				COALESCE(pa.saves, 0)::bigint AS saves,
				LAG(COALESCE(pa.likes, 0)) OVER (PARTITION BY pa.post_id ORDER BY pa.date) AS prev_likes,
				LAG(COALESCE(pa.comments, 0)) OVER (PARTITION BY pa.post_id ORDER BY pa.date) AS prev_comments,
				LAG(COALESCE(pa.shares, 0)) OVER (PARTITION BY pa.post_id ORDER BY pa.date) AS prev_shares,
				LAG(COALESCE(pa.saves, 0)) OVER (PARTITION BY pa.post_id ORDER BY pa.date) AS prev_saves
			FROM instagram_post_analytics pa
			INNER JOIN posts p ON p.id = pa.post_id AND $1 = ANY(p.channel_ids)
			WHERE pa.channel_id = $1
				AND pa.date::date <= $3::date
		),
		deltas AS (
			SELECT
				d,
				CASE WHEN prev_likes IS NULL THEN 0 ELSE GREATEST(0, likes - prev_likes) END AS dl,
				CASE WHEN prev_comments IS NULL THEN 0 ELSE GREATEST(0, comments - prev_comments) END AS dc,
				CASE WHEN prev_shares IS NULL THEN 0 ELSE GREATEST(0, shares - prev_shares) END AS ds,
				CASE WHEN prev_saves IS NULL THEN 0 ELSE GREATEST(0, saves - prev_saves) END AS dz
			FROM ordered
			WHERE d >= $2::date AND d <= $3::date
		)
		SELECT d, COALESCE(SUM(dl), 0), COALESCE(SUM(dc), 0), COALESCE(SUM(ds), 0), COALESCE(SUM(dz), 0)
		FROM deltas
		GROUP BY d
		ORDER BY d ASC
	`
	rows, err := db.QueryContext(ctx, query, channelID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []DailyPostEngagementRow
	for rows.Next() {
		var r DailyPostEngagementRow
		if err := rows.Scan(&r.Date, &r.Likes, &r.Comments, &r.Shares, &r.Saves); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func nullableString(s *string) interface{} {
	if s == nil {
		return nil
	}
	return *s
}

func previousCalendarPeriod(startDate, endDate time.Time) (prevStart, prevEnd time.Time) {
	startDate = startDate.UTC()
	endDate = endDate.UTC()
	days := int(endDate.Sub(startDate).Hours()/24) + 1
	if days < 1 {
		days = 1
	}
	prevEnd = startDate.AddDate(0, 0, -1)
	prevStart = prevEnd.AddDate(0, 0, -(days - 1))
	y, m, d := prevStart.Date()
	prevStart = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	y2, m2, d2 := prevEnd.Date()
	prevEnd = time.Date(y2, m2, d2, 0, 0, 0, 0, time.UTC)
	return prevStart, prevEnd
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

	// Get new posts count (calendar-day inclusive range in UTC)
	newPostsQuery := `
		SELECT COUNT(*) 
		FROM posts 
		WHERE $1 = ANY(channel_ids)
			AND published_at IS NOT NULL
			AND (published_at AT TIME ZONE 'UTC')::date >= $2::date
			AND (published_at AT TIME ZONE 'UTC')::date <= $3::date
	`

	// Get total posts count
	totalPostsQuery := `
		SELECT COUNT(*) 
		FROM posts 
		WHERE $1 = ANY(channel_ids)
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
			SELECT DISTINCT ON (pa.post_id) 
				pa.post_id, pa.likes, pa.comments, pa.saves, pa.shares, pa.reach, pa.impressions
			FROM instagram_post_analytics pa
			JOIN posts p ON p.id = pa.post_id
			WHERE pa.channel_id = $1 AND pa.date >= $2::date AND pa.date <= $3::date
			ORDER BY pa.post_id, pa.date DESC
		) ipa
	`

	prevStartDate, prevEndDate := previousCalendarPeriod(startDate, endDate)

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
			WHERE channel_id = $1 AND date >= $2::date AND date <= $3::date
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

	err = db.QueryRowContext(ctx, prevAnalyticsQuery, channelID, prevStartDate.UTC(), prevEndDate.UTC()).Scan(
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

// StoryAnalyticsJoinedRow is one story snapshot joined with post metadata.
type StoryAnalyticsJoinedRow struct {
	entities.InstagramStoryAnalytics
	Caption     string
	PostType    string
	PublishedAt *time.Time
}

// CountStoryAnalyticsInRange returns rows matching filters.
func CountStoryAnalyticsInRange(ctx context.Context, channelID int64, startDate, endDate time.Time, postID *int64) (int, error) {
	db := database.GetDB()
	q := `
		SELECT COUNT(*) FROM instagram_story_analytics sa
		JOIN posts p ON p.id = sa.post_id AND $1 = ANY(p.channel_ids)
		WHERE sa.channel_id = $1 AND sa.date >= $2::date AND sa.date <= $3::date
	`
	args := []interface{}{channelID, startDate, endDate}
	if postID != nil {
		q += ` AND sa.post_id = $4`
		args = append(args, *postID)
	}
	var n int
	err := db.QueryRowContext(ctx, q, args...).Scan(&n)
	return n, err
}

// GetStoryAnalyticsJoined returns story analytics with post caption/type, paginated.
func GetStoryAnalyticsJoined(ctx context.Context, channelID int64, startDate, endDate time.Time, postID *int64, limit, offset int) ([]StoryAnalyticsJoinedRow, error) {
	db := database.GetDB()
	q := `
		SELECT sa.id, sa.channel_id, sa.post_id, sa.date, sa.impressions, sa.reach, sa.exits, sa.forwards,
		       sa.replies, sa.taps_forward, sa.taps_backward, sa.taps_exit, sa.raw, sa.created_at,
		       COALESCE(p.caption, ''), COALESCE(p.post_type, ''), p.published_at
		FROM instagram_story_analytics sa
		JOIN posts p ON p.id = sa.post_id AND $1 = ANY(p.channel_ids)
		WHERE sa.channel_id = $1 AND sa.date >= $2::date AND sa.date <= $3::date
	`
	args := []interface{}{channelID, startDate, endDate}
	argN := 4
	if postID != nil {
		q += fmt.Sprintf(" AND sa.post_id = $%d", argN)
		args = append(args, *postID)
		argN++
	}
	q += fmt.Sprintf(" ORDER BY sa.date DESC, sa.id DESC LIMIT $%d OFFSET $%d", argN, argN+1)
	args = append(args, limit, offset)

	rows, err := db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []StoryAnalyticsJoinedRow
	for rows.Next() {
		var r StoryAnalyticsJoinedRow
		var pubAt sql.NullTime
		err := rows.Scan(
			&r.ID,
			&r.ChannelID,
			&r.PostID,
			&r.Date,
			&r.Impressions,
			&r.Reach,
			&r.Exits,
			&r.Forwards,
			&r.Replies,
			&r.TapsForward,
			&r.TapsBackward,
			&r.TapsExit,
			&r.Raw,
			&r.CreatedAt,
			&r.Caption,
			&r.PostType,
			&pubAt,
		)
		if err != nil {
			return nil, err
		}
		if pubAt.Valid {
			t := pubAt.Time
			r.PublishedAt = &t
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// UpsertAnalyticsSnapshot upserts an analytics snapshot
func UpsertAnalyticsSnapshot(ctx context.Context, snapshot *entities.AnalyticsSnapshot) error {
	db := database.GetDB()
	query := `
		INSERT INTO analytics_snapshots 
		(entity_type, entity_id, period_type, start_date, end_date, metrics, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		ON CONFLICT (entity_type, entity_id, period_type, start_date) 
		DO UPDATE SET
			end_date = EXCLUDED.end_date,
			metrics = EXCLUDED.metrics,
			updated_at = NOW()
		RETURNING id, created_at, updated_at
	`
	return db.QueryRowContext(ctx, query,
		snapshot.EntityType,
		snapshot.EntityID,
		snapshot.PeriodType,
		snapshot.StartDate,
		snapshot.EndDate,
		snapshot.Metrics,
		time.Now(),
	).Scan(&snapshot.ID, &snapshot.CreatedAt, &snapshot.UpdatedAt)
}

// GetAnalyticsSnapshot retrieves a specific snapshot
func GetAnalyticsSnapshot(ctx context.Context, entityType, periodType string, entityID int64, startDate time.Time) (*entities.AnalyticsSnapshot, error) {
	db := database.GetDB()
	query := `
		SELECT id, entity_type, entity_id, period_type, start_date, end_date, metrics, created_at, updated_at
		FROM analytics_snapshots
		WHERE entity_type = $1 AND entity_id = $2 AND period_type = $3 AND start_date = $4
	`
	var s entities.AnalyticsSnapshot
	err := db.QueryRowContext(ctx, query, entityType, entityID, periodType, startDate).Scan(
		&s.ID, &s.EntityType, &s.EntityID, &s.PeriodType, &s.StartDate, &s.EndDate, &s.Metrics, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// GetLatestAnalyticsSnapshot retrieves the latest snapshot for an entity
func GetLatestAnalyticsSnapshot(ctx context.Context, entityType, periodType string, entityID int64) (*entities.AnalyticsSnapshot, error) {
	db := database.GetDB()
	query := `
		SELECT id, entity_type, entity_id, period_type, start_date, end_date, metrics, created_at, updated_at
		FROM analytics_snapshots
		WHERE entity_type = $1 AND entity_id = $2 AND period_type = $3
		ORDER BY start_date DESC
		LIMIT 1
	`
	var s entities.AnalyticsSnapshot
	err := db.QueryRowContext(ctx, query, entityType, entityID, periodType).Scan(
		&s.ID, &s.EntityType, &s.EntityID, &s.PeriodType, &s.StartDate, &s.EndDate, &s.Metrics, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// AudienceSnapshotListFilters optional pagination for audience snapshot history.
type AudienceSnapshotListFilters struct {
	Limit  int
	Offset int
}

// CountAudienceSnapshotsInRange counts instagram_audience_snapshots rows for a channel in [startDate, endDate].
func CountAudienceSnapshotsInRange(ctx context.Context, channelID int64, startDate, endDate time.Time) (int, error) {
	db := database.GetDB()
	var n int
	err := db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM instagram_audience_snapshots
		WHERE channel_id = $1 AND snapshot_date >= $2::date AND snapshot_date <= $3::date
	`, channelID, startDate, endDate).Scan(&n)
	return n, err
}

// GetAudienceSnapshotsByDateRange returns audience snapshots newest-first.
func GetAudienceSnapshotsByDateRange(ctx context.Context, channelID int64, startDate, endDate time.Time, f AudienceSnapshotListFilters) ([]entities.InstagramAudienceSnapshot, error) {
	db := database.GetDB()
	query := `
		SELECT id, channel_id, snapshot_date, raw, created_at
		FROM instagram_audience_snapshots
		WHERE channel_id = $1 AND snapshot_date >= $2::date AND snapshot_date <= $3::date
		ORDER BY snapshot_date DESC
	`
	args := []interface{}{channelID, startDate, endDate}
	if f.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
		args = append(args, f.Limit, f.Offset)
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []entities.InstagramAudienceSnapshot
	for rows.Next() {
		var s entities.InstagramAudienceSnapshot
		if err := rows.Scan(&s.ID, &s.ChannelID, &s.SnapshotDate, &s.Raw, &s.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}
