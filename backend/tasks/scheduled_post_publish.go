package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/shekhar8352/PostEaze/scheduledpublish"
)

// EnqueueInstagramScheduledPostPublish schedules an Asynq job to publish at runAt (UTC).
func EnqueueInstagramScheduledPostPublish(scheduledPostID int64, runAt time.Time) (*asynq.TaskInfo, error) {
	if GetClient() == nil {
		return nil, fmt.Errorf("asynq client is not initialized")
	}
	p := ScheduledPostInstagramPublishPayload{ScheduledPostID: scheduledPostID}
	body, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	task := asynq.NewTask(TypeInstagramScheduledPostPublish, body)
	return EnqueueTask(task, QueueMedium, asynq.ProcessAt(runAt))
}

// HandleInstagramScheduledPostPublish runs deferred Instagram publishing (immediate Meta flow).
func HandleInstagramScheduledPostPublish(ctx context.Context, t *asynq.Task) error {
	var p ScheduledPostInstagramPublishPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal: %v: %w", err, asynq.SkipRetry)
	}
	if p.ScheduledPostID <= 0 {
		return fmt.Errorf("invalid scheduled_post_id: %w", asynq.SkipRetry)
	}
	return scheduledpublish.ExecuteInstagramPublishAtScheduledTime(ctx, p.ScheduledPostID)
}
