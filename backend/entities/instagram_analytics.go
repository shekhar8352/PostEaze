package entities

import (
	"time"
)

// InstagramPostAnalytics represents daily analytics for an Instagram post
type InstagramPostAnalytics struct {
	ID                int64     `db:"id"`
	ChannelID         int64     `db:"channel_id"`
	PostID            int64     `db:"post_id"`
	Date              time.Time `db:"date"`
	Impressions       *int      `db:"impressions"`
	Reach             *int      `db:"reach"`
	Likes             *int      `db:"likes"`
	Comments          *int      `db:"comments"`
	Saves             *int      `db:"saves"`
	Shares            *int      `db:"shares"`
	VideoViews        *int      `db:"video_views"`
	ProfileVisits     *int      `db:"profile_visits"`
	Follows           *int      `db:"follows"`
	Views             *int      `db:"views"`              // Usage started after July 2024
	TotalInteractions *int      `db:"total_interactions"` // Sum of likes, comments, shares, saves
	EngagementRate    *float64  `db:"engagement_rate"`    // Calculated field
	Plays             *int      `db:"plays"`              // For Reels/Videos
	Raw               []byte    `db:"raw"`                // JSONB
	CreatedAt         time.Time `db:"created_at"`
}

// InstagramStoryAnalytics represents analytics for an Instagram story
type InstagramStoryAnalytics struct {
	ID           int64     `db:"id"`
	ChannelID    int64     `db:"channel_id"`
	PostID       int64     `db:"post_id"`
	Date         time.Time `db:"date"`
	Impressions  *int      `db:"impressions"`
	Reach        *int      `db:"reach"`
	Exits        *int      `db:"exits"`
	Forwards     *int      `db:"forwards"`
	Replies      *int      `db:"replies"`
	TapsForward  *int      `db:"taps_forward"`
	TapsBackward *int      `db:"taps_backward"`
	TapsExit     *int      `db:"taps_exit"`
	Raw          []byte    `db:"raw"` // JSONB
	CreatedAt    time.Time `db:"created_at"`
}

// InstagramProfileAnalytics represents daily analytics for an Instagram profile
type InstagramProfileAnalytics struct {
	ID                  int64     `db:"id"`
	ChannelID           int64     `db:"channel_id"`
	Date                time.Time `db:"date"`
	FollowerCount       *int      `db:"follower_count"`
	Impressions         *int      `db:"impressions"`
	ProfileViews        *int      `db:"profile_views"`
	Reach               *int      `db:"reach"`
	WebsiteClicks       *int      `db:"website_clicks"`
	EmailClicks         *int      `db:"email_clicks"`
	Views               *int      `db:"views"`
	AccountsEngaged     *int      `db:"accounts_engaged"`
	TotalInteractions   *int      `db:"total_interactions"`
	BioLinkClicks       *int      `db:"bio_link_clicks"`
	PhoneCallClicks     *int      `db:"phone_call_clicks"`
	TextMessageClicks   *int      `db:"text_message_clicks"`
	GetDirectionsClicks *int      `db:"get_directions_clicks"`
	Raw                 []byte    `db:"raw"` // JSONB
	CreatedAt           time.Time `db:"created_at"`
}

// InstagramAudienceSnapshot stores Meta lifetime audience_* insight payloads (merged) per UTC day.
type InstagramAudienceSnapshot struct {
	ID           int64     `db:"id"`
	ChannelID    int64     `db:"channel_id"`
	SnapshotDate time.Time `db:"snapshot_date"`
	Raw          []byte    `db:"raw"` // JSONB: { "data": [ ... InsightData ... ] }
	CreatedAt    time.Time `db:"created_at"`
}
