package entities

import (
	"time"
)

// Post represents a social media post from any provider
type Post struct {
	ID             int64      `db:"id"`
	ChannelID      int64      `db:"channel_id"`
	Provider       string     `db:"provider"`
	ProviderPostID string     `db:"provider_post_id"`
	Source         string     `db:"source"` // 'native' | 'posteaze'
	PostType       *string    `db:"post_type"`
	Caption        *string    `db:"caption"`
	Media          []byte     `db:"media"` // JSONB
	PublishedAt    *time.Time `db:"published_at"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
}
