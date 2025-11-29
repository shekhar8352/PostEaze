package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

// RegisterHandlers registers task handlers to the mux
func RegisterHandlers(mux *asynq.ServeMux) {
	mux.HandleFunc(TypeEmailDelivery, HandleEmailDeliveryTask)
	mux.HandleFunc(TypeLogMessage, HandleLogMessageTask)
	mux.HandleFunc(TypeInstagramComment, HandleInstagramCommentTask)
	mux.HandleFunc(TypeInstagramMention, HandleInstagramMentionTask)
	mux.HandleFunc(TypeInstagramStoryInsight, HandleInstagramStoryInsightTask)
}

// HandleEmailDeliveryTask handles email delivery tasks
func HandleEmailDeliveryTask(ctx context.Context, t *asynq.Task) error {
	var p EmailDeliveryPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Sending Email to User: %s, Subject: %s", p.UserID, p.Subject)
	// Logic to send email would go here
	return nil
}

// HandleLogMessageTask handles log message tasks
func HandleLogMessageTask(ctx context.Context, t *asynq.Task) error {
	var p LogMessagePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Log Task: %s", p.Message)
	return nil
}

// HandleInstagramCommentTask handles Instagram comment events
func HandleInstagramCommentTask(ctx context.Context, t *asynq.Task) error {
	var change map[string]interface{}
	if err := json.Unmarshal(t.Payload(), &change); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Processing Instagram Comment: %v", change)
	// TODO: Implement comment processing logic (e.g., save to DB, notify user)
	return nil
}

// HandleInstagramMentionTask handles Instagram mention events
func HandleInstagramMentionTask(ctx context.Context, t *asynq.Task) error {
	var change map[string]interface{}
	if err := json.Unmarshal(t.Payload(), &change); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Processing Instagram Mention: %v", change)
	// TODO: Implement mention processing logic
	return nil
}

// HandleInstagramStoryInsightTask handles Instagram story insight events
func HandleInstagramStoryInsightTask(ctx context.Context, t *asynq.Task) error {
	var change map[string]interface{}
	if err := json.Unmarshal(t.Payload(), &change); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Processing Instagram Story Insight: %v", change)
	// TODO: Implement story insight processing logic
	return nil
}
