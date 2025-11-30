package tasks

// Task Types
const (
	TypeEmailDelivery         = "email:deliver"
	TypeLogMessage            = "log:message"
	TypeInstagramComment      = "instagram:comment"
	TypeInstagramMention      = "instagram:mention"
	TypeInstagramStoryInsight = "instagram:story_insight"
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
