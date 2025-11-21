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
