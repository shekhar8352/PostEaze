package entities

import (
	"time"
)

// FacebookPageAnalytics is daily aggregated Page insights for a Facebook channel.
type FacebookPageAnalytics struct {
	ID               int64     `db:"id"`
	ChannelID        int64     `db:"channel_id"`
	Date             time.Time `db:"date"`
	Followers        *int      `db:"followers"`
	Likes            *int      `db:"likes"`
	Reach            *int      `db:"reach"`
	Impressions      *int      `db:"impressions"`
	PageEngagedUsers *int      `db:"page_engaged_users"`
	Raw              []byte    `db:"raw"`
	CreatedAt        time.Time `db:"created_at"`
}

// FacebookPostAnalytics is daily analytics for a Facebook Page post.
type FacebookPostAnalytics struct {
	ID          int64     `db:"id"`
	ChannelID   int64     `db:"channel_id"`
	PostID      int64     `db:"post_id"`
	Date        time.Time `db:"date"`
	Impressions *int      `db:"impressions"`
	Reach       *int      `db:"reach"`
	Likes       *int      `db:"likes"`
	Comments    *int      `db:"comments"`
	Shares      *int      `db:"shares"`
	PhotoViews  *int      `db:"photo_views"`
	VideoViews  *int      `db:"video_views"`
	Raw         []byte    `db:"raw"`
	CreatedAt   time.Time `db:"created_at"`
}
