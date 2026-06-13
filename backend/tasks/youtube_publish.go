package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/shekhar8352/PostEaze/scheduledpublish"
)

type ScheduledPostYouTubePublishPayload struct {
	ScheduledPostID int64 `json:"scheduled_post_id"`
}

func EnqueueYouTubeScheduledPostPublish(scheduledPostID int64, runAt time.Time) (*asynq.TaskInfo, error) {
	if GetClient() == nil {
		return nil, fmt.Errorf("asynq client is not initialized")
	}
	p := ScheduledPostYouTubePublishPayload{ScheduledPostID: scheduledPostID}
	body, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	task := asynq.NewTask(TypeYouTubeScheduledPostPublish, body)
	return EnqueueTask(task, QueueSlow, asynq.ProcessAt(runAt), asynq.Timeout(2*time.Hour))
}

func HandleYouTubeScheduledPostPublish(ctx context.Context, t *asynq.Task) error {
	var p ScheduledPostYouTubePublishPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal: %v: %w", err, asynq.SkipRetry)
	}
	if p.ScheduledPostID <= 0 {
		return fmt.Errorf("invalid scheduled_post_id: %w", asynq.SkipRetry)
	}
	return scheduledpublish.ExecuteYouTubePublish(ctx, p.ScheduledPostID)
}
