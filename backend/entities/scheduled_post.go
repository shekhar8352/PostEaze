package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// ScheduledPostStatus is stored as text in DB.
type ScheduledPostStatus string

const (
	ScheduledStatusPending    ScheduledPostStatus = "pending"
	ScheduledStatusSubmitting ScheduledPostStatus = "submitting"
	ScheduledStatusScheduled  ScheduledPostStatus = "scheduled"
	ScheduledStatusFailed     ScheduledPostStatus = "failed"
	ScheduledStatusCancelled  ScheduledPostStatus = "cancelled"
	ScheduledStatusPublished  ScheduledPostStatus = "published"
)

// ScheduledPost is outbound scheduled content (not yet merged into analytics `posts`).
type ScheduledPost struct {
	ID             int64               `db:"id"`
	OwnerUserID    uuid.UUID           `db:"owner_user_id"`
	ChannelIDs     pq.Int64Array       `db:"channel_ids"`
	Platforms      pq.StringArray      `db:"platforms"`
	ScheduledAt    time.Time           `db:"scheduled_at"`
	Status         string              `db:"status"`
	PostType       string              `db:"post_type"`
	Caption        *string             `db:"caption"`
	Media          []byte              `db:"media"`           // JSONB
	ProviderState  []byte              `db:"provider_state"` // JSONB per-channel results
	CreatedAt      time.Time           `db:"created_at"`
	UpdatedAt      time.Time           `db:"updated_at"`
}
