package modelsv1

import "time"

// ScheduledMediaItem is one URL in the media JSON payload.
type ScheduledMediaItem struct {
	URL  string `json:"url" binding:"required"`
	Kind string `json:"kind" binding:"required"` // image | video
}

// CreateScheduledPostRequest is the body for POST /scheduled-posts.
type CreateScheduledPostRequest struct {
	ChannelIDs []int64              `json:"channel_ids" binding:"required,min=1"`
	Platforms  []string             `json:"platforms" binding:"required,min=1"`
	ScheduledAt time.Time         `json:"scheduled_at" binding:"required"`
	PostType   string              `json:"post_type" binding:"required"` // image | video | carousel
	Caption    string              `json:"caption"`
	Media      ScheduledMediaPayload `json:"media" binding:"required"`
}

// ScheduledMediaPayload wraps items for storage and validation.
type ScheduledMediaPayload struct {
	Items []ScheduledMediaItem `json:"items" binding:"required,min=1"`
}

// ChannelScheduleResult is one channel’s outcome after Meta calls.
type ChannelScheduleResult struct {
	ChannelID    int64  `json:"channel_id"`
	Success      bool   `json:"success"`
	CreationID   string `json:"creation_id,omitempty"`
	PublishedID  string `json:"published_media_id,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// CreateScheduledPostResponse is returned after scheduling (may be multi-status).
type CreateScheduledPostResponse struct {
	ScheduledPostID int64                   `json:"scheduled_post_id"`
	OverallStatus   string                  `json:"overall_status"` // scheduled | partial_failure | failed
	Results         []ChannelScheduleResult `json:"results"`
}

// ScheduledPostListItem is a calendar-friendly row.
type ScheduledPostListItem struct {
	ID           int64    `json:"id"`
	ChannelIDs   []int64  `json:"channel_ids"`
	Platforms    []string `json:"platforms"`
	ScheduledAt  string   `json:"scheduled_at"`
	Status       string   `json:"status"`
	PostType     string   `json:"post_type"`
	Caption      string   `json:"caption,omitempty"`
	Media        any      `json:"media"`
	ProviderState any     `json:"provider_state,omitempty"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
}

// ListScheduledPostsResponse for GET list.
type ListScheduledPostsResponse struct {
	Posts []ScheduledPostListItem `json:"posts"`
}

// ListScheduledPostsQuery binds query params.
type ListScheduledPostsQuery struct {
	From      string `form:"from" binding:"required"` // ISO date or RFC3339
	To        string `form:"to" binding:"required"`
	ChannelID *int64 `form:"channel_id"`
}
