package scheduledpublish

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/entities/repositories"
	"github.com/shekhar8352/PostEaze/mediapublish"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/services/publishing"
)

var ytPublisher = publishing.NewYouTubePublisher()

func ExecuteYouTubePublish(ctx context.Context, scheduledPostID int64) error {
	sp, err := repositories.GetScheduledPostByIDForJob(ctx, scheduledPostID)
	if err != nil {
		return err
	}
	if sp == nil {
		return fmt.Errorf("scheduled post %d not found: %w", scheduledPostID, asynq.SkipRetry)
	}

	switch sp.Status {
	case string(entities.ScheduledStatusCancelled), string(entities.ScheduledStatusPublished):
		return nil
	case string(entities.ScheduledStatusScheduled), string(entities.ScheduledStatusSubmitting), string(entities.ScheduledStatusPending):
		// proceed
	default:
		return fmt.Errorf("scheduled post %d has status %q: %w", scheduledPostID, sp.Status, asynq.SkipRetry)
	}

	if sp.PostType != "video" {
		return fmt.Errorf("youtube only supports video posts: %w", asynq.SkipRetry)
	}

	var media modelsv1.ScheduledMediaPayload
	if err := json.Unmarshal(sp.Media, &media); err != nil {
		return fmt.Errorf("%w: %v", asynq.SkipRetry, err)
	}
	if len(media.Items) != 1 || media.Items[0].Kind != "video" {
		return fmt.Errorf("youtube requires one video media item: %w", asynq.SkipRetry)
	}

	var version *entities.MediaVersion
	item := media.Items[0]
	if item.MediaVersionID != nil {
		version, err = publishing.LoadVersionForYouTube(ctx, sp.OwnerUserID, *item.MediaVersionID)
		if err != nil {
			return fmt.Errorf("%w: %v", asynq.SkipRetry, err)
		}
	} else if item.URL != "" {
		version = &entities.MediaVersion{
			StorageProvider: entities.StorageProviderBlob,
			BlobURL:         item.URL,
			ContentType:     "video/mp4",
			FileSize:        0,
			FileName:        "video",
		}
	} else {
		return fmt.Errorf("youtube publish requires media_version_id or url: %w", asynq.SkipRetry)
	}

	caption := ""
	if sp.Caption != nil {
		caption = *sp.Caption
	}

	stateMap := map[string]any{
		"youtube": map[string]any{"channels": map[string]any{}},
	}
	channelsMap, _ := stateMap["youtube"].(map[string]any)
	inner, _ := channelsMap["channels"].(map[string]any)

	_ = repositories.UpdateScheduledPostStatusAndProviderState(ctx, sp.ID, sp.OwnerUserID, string(entities.ScheduledStatusSubmitting), []byte(`{}`))

	successN := 0
	channelIDs := []int64(sp.ChannelIDs)
	for _, cid := range channelIDs {
		ch, err := repositories.GetChannelByID(ctx, cid)
		if err != nil || ch == nil || ch.Provider != "youtube" {
			inner[fmt.Sprintf("%d", cid)] = map[string]any{"error": "not a youtube channel"}
			continue
		}
		videoID, pubErr := ytPublisher.Publish(ctx, publishing.YouTubePublishInput{
			ChannelID: cid,
			OwnerID:   sp.OwnerUserID,
			Title:     caption,
			Caption:   caption,
			Version:   version,
		})
		if pubErr != nil {
			inner[fmt.Sprintf("%d", cid)] = map[string]any{"error": pubErr.Error()}
		} else {
			successN++
			inner[fmt.Sprintf("%d", cid)] = map[string]any{"video_id": videoID}
		}
	}

	stateBytes, _ := json.Marshal(stateMap)
	finalStatus := string(entities.ScheduledStatusFailed)
	if successN > 0 && successN == countYouTubeChannels(ctx, channelIDs) {
		finalStatus = string(entities.ScheduledStatusPublished)
	}
	_ = repositories.UpdateScheduledPostStatusAndProviderState(ctx, sp.ID, sp.OwnerUserID, finalStatus, stateBytes)
	if finalStatus == string(entities.ScheduledStatusPublished) {
		mediapublish.MarkLinkedMediaAssetsPublished(ctx, sp.OwnerUserID, sp.Media)
	}
	if OnPostFinalized != nil {
		_ = OnPostFinalized(ctx, sp.ID, finalStatus == string(entities.ScheduledStatusPublished))
	}
	if successN < countYouTubeChannels(ctx, channelIDs) {
		return fmt.Errorf("youtube publish incomplete for scheduled_post %d: %w", scheduledPostID, asynq.SkipRetry)
	}
	return nil
}

func countYouTubeChannels(ctx context.Context, channelIDs []int64) int {
	n := 0
	for _, cid := range channelIDs {
		ch, err := repositories.GetChannelByID(ctx, cid)
		if err == nil && ch != nil && ch.Provider == "youtube" {
			n++
		}
	}
	return n
}
