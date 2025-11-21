package tasks

// Task Types
const (
	TypeEmailDelivery = "email:deliver"
	TypeLogMessage    = "log:message"
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
