package tasks

// Task Types
const (
	TypeEmailDelivery          = "email:deliver"
	TypeLogMessage             = "log:message"
	TypeInstagramComment       = "instagram:comment"
	TypeInstagramMention       = "instagram:mention"
	TypeInstagramStoryInsight  = "instagram:story_insight"
	TypeSyncInstagramProfiles  = "instagram:sync_profiles"
	TypeSyncInstagramPosts     = "instagram:sync_posts"
	TypePeriodSnapshot         = "analytics:period_snapshot"
	TypeSyncInstagramAnalytics = "instagram:sync_analytics"
)

// Queue Names
const (
	QueueFast   = "fast"
	QueueMedium = "medium"
	QueueSlow   = "slow"
)

// Payload Structs
type EmailDeliveryPayload struct {
	UserID  string
	Subject string
	Body    string
}

type LogMessagePayload struct {
	Message string
}
