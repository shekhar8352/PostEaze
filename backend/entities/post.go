package entities

import (
	"time"

	"github.com/lib/pq"
)

// Post represents a social media post from any provider
type Post struct {
	ID         int64         `db:"id"`
	ChannelIDs pq.Int64Array `db:"channel_ids"`
	OwnerID    *string       `db:"owner_id"` // UUID as string, nullable

	// Provider tracking
	Providers       pq.StringArray `db:"providers"`         // e.g., ['instagram']
	ProviderPostIDs []byte         `db:"provider_post_ids"` // JSONB: {"instagram": "post_id"}

	Source      string     `db:"source"` // 'native' | 'posteaze'
	PostType    *string    `db:"post_type"`
	Caption     *string    `db:"caption"`
	Media       []byte     `db:"media"` // JSONB
	PublishedAt *time.Time `db:"published_at"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}
