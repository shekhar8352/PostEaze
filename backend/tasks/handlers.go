package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/shekhar8352/PostEaze/utils"
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
	utils.Logger.Info(ctx, "Sending Email to User: ", p.UserID, ", Subject: ", p.Subject)
	// Logic to send email would go here
	return nil
}

// HandleLogMessageTask handles log message tasks
func HandleLogMessageTask(ctx context.Context, t *asynq.Task) error {
	var p LogMessagePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	utils.Logger.Info(ctx, "Log Task: ", p.Message)
	return nil
}

// HandleInstagramCommentTask handles Instagram comment events
func HandleInstagramCommentTask(ctx context.Context, t *asynq.Task) error {
	var change map[string]interface{}
	if err := json.Unmarshal(t.Payload(), &change); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	utils.Logger.Info(ctx, "Processing Instagram Comment: ", change)
	// TODO: Implement comment processing logic (e.g., save to DB, notify user)
	return nil
}

// HandleInstagramMentionTask handles Instagram mention events
func HandleInstagramMentionTask(ctx context.Context, t *asynq.Task) error {
	var change map[string]interface{}
	if err := json.Unmarshal(t.Payload(), &change); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	utils.Logger.Info(ctx, "Processing Instagram Mention: ", change)
	// TODO: Implement mention processing logic
	return nil
}

// HandleInstagramStoryInsightTask handles Instagram story insight events
func HandleInstagramStoryInsightTask(ctx context.Context, t *asynq.Task) error {
	var change map[string]interface{}
	if err := json.Unmarshal(t.Payload(), &change); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	utils.Logger.Info(ctx, "Processing Instagram Story Insight: ", change)
	// TODO: Implement story insight processing logic
	return nil
}
