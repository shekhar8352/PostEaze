package publishing

import (
	"context"
	"time"
)

// Platform identifiers for multi-platform expansion.
const (
	PlatformInstagram = "instagram"
	PlatformYouTube   = "youtube"
)

// SchedulePayload is platform-agnostic input for a single channel publish attempt.
type SchedulePayload struct {
	PostType     string // image | video | carousel
	Caption      string
	PublishNow   bool // if true, omit Meta scheduled_publish_time (immediate feed post)
	ScheduledAt  time.Time // UTC; used for Meta scheduled_publish_time when !PublishNow
	ImageURL     string
	VideoURL     string
	CarouselURLs []string
}

// Publisher schedules or publishes content for one connected account.
type Publisher interface {
	Platform() string
	Schedule(ctx context.Context, accessToken, igUserID string, payload SchedulePayload) (creationID, publishedMediaID string, err error)
}
