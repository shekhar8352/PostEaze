package scheduledpublish

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/hibiken/asynq"
	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/entities/repositories"
	"github.com/shekhar8352/PostEaze/mediapublish"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/provider/instagram"
	"github.com/shekhar8352/PostEaze/services/publishing"
	"github.com/shekhar8352/PostEaze/utils/encryption"
)

var igPublisher = publishing.NewInstagramPublisher(instagram.NewInstagramProvider())

// ExecuteInstagramPublishAtScheduledTime loads the row and performs immediate Instagram
// publishing (no Meta scheduled_publish_time). Intended to run from an Asynq task at scheduled_at.
func ExecuteInstagramPublishAtScheduledTime(ctx context.Context, scheduledPostID int64) error {
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
	case string(entities.ScheduledStatusScheduled):
		// proceed
	default:
		return fmt.Errorf("scheduled post %d has status %q, expected scheduled: %w", scheduledPostID, sp.Status, asynq.SkipRetry)
	}

	if sp.PostType == "reel" || sp.PostType == "story" {
		return fmt.Errorf("unsupported post_type %q: %w", sp.PostType, asynq.SkipRetry)
	}

	payload, err := schedulePayloadFromStoredPost(sp)
	if err != nil {
		return fmt.Errorf("%w: %v", asynq.SkipRetry, err)
	}

	stateMap := map[string]any{
		"publish_now": true,
		"queue":       "asynq",
		"instagram":   map[string]any{"channels": map[string]any{}},
	}
	channelsMap, _ := stateMap["instagram"].(map[string]any)
	inner, _ := channelsMap["channels"].(map[string]any)

	_ = repositories.UpdateScheduledPostStatusAndProviderState(ctx, sp.ID, sp.OwnerUserID, string(entities.ScheduledStatusSubmitting), []byte(`{}`))

	pubPayload := payload
	pubPayload.PublishNow = true
	pubPayload.ScheduledAt = time.Time{}

	successN := 0
	var channelIDs []int64
	for _, cid := range sp.ChannelIDs {
		ch, err := repositories.GetChannelByID(ctx, cid)
		if err == nil && ch != nil && ch.Provider == "instagram" {
			channelIDs = append(channelIDs, cid)
		}
	}
	for _, cid := range channelIDs {
		tok, err := repositories.GetLatestTokenByChannelID(ctx, cid)
		if err != nil {
			inner[fmt.Sprintf("%d", cid)] = map[string]any{"error": "token: " + err.Error()}
			continue
		}
		plain, err := encryption.Decrypt(tok.AccessToken)
		if err != nil {
			inner[fmt.Sprintf("%d", cid)] = map[string]any{"error": "decrypt token: " + err.Error()}
			continue
		}
		ch, err := repositories.GetChannelByID(ctx, cid)
		if err != nil {
			inner[fmt.Sprintf("%d", cid)] = map[string]any{"error": "channel: " + err.Error()}
			continue
		}
		igUserID, err := instagram.UserIDFromChannelMetadata(ch.Metadata)
		if err != nil {
			inner[fmt.Sprintf("%d", cid)] = map[string]any{"error": err.Error()}
			continue
		}
		creationID, publishedID, pubErr := igPublisher.Schedule(ctx, plain, igUserID, pubPayload)
		if pubErr != nil {
			inner[fmt.Sprintf("%d", cid)] = map[string]any{"error": pubErr.Error()}
		} else {
			successN++
			inner[fmt.Sprintf("%d", cid)] = map[string]any{
				"creation_id":          creationID,
				"published_media_id": publishedID,
			}
		}
	}

	stateBytes, _ := json.Marshal(stateMap)
	finalStatus := string(entities.ScheduledStatusFailed)
	if successN == len(channelIDs) {
		finalStatus = string(entities.ScheduledStatusPublished)
	}
	_ = repositories.UpdateScheduledPostStatusAndProviderState(ctx, sp.ID, sp.OwnerUserID, finalStatus, stateBytes)
	if finalStatus == string(entities.ScheduledStatusPublished) {
		mediapublish.MarkLinkedMediaAssetsPublished(ctx, sp.OwnerUserID, sp.Media)
	}
	// Propagate lifecycle to any linked Studio Piece (move to published/log failure).
	if OnPostFinalized != nil {
		_ = OnPostFinalized(ctx, sp.ID, finalStatus == string(entities.ScheduledStatusPublished))
	}
	if successN < len(channelIDs) {
		// Avoid Asynq retries re-publishing channels that already succeeded.
		return fmt.Errorf("instagram publish incomplete for scheduled_post %d (%d/%d channels ok): %w",
			scheduledPostID, successN, len(channelIDs), asynq.SkipRetry)
	}
	return nil
}

func schedulePayloadFromStoredPost(sp *entities.ScheduledPost) (publishing.SchedulePayload, error) {
	var out publishing.SchedulePayload
	out.PostType = sp.PostType
	if sp.Caption != nil {
		out.Caption = *sp.Caption
	}
	var media modelsv1.ScheduledMediaPayload
	if err := json.Unmarshal(sp.Media, &media); err != nil {
		return out, fmt.Errorf("media json: %w", err)
	}
	items := media.Items
	switch sp.PostType {
	case "image":
		if len(items) != 1 || items[0].Kind != "image" {
			return out, fmt.Errorf("image post requires exactly one media item with kind image")
		}
		if err := mustInstagramFetchableMediaURL(items[0].URL); err != nil {
			return out, err
		}
		out.ImageURL = items[0].URL
	case "video":
		if len(items) != 1 || items[0].Kind != "video" {
			return out, fmt.Errorf("video post requires exactly one media item with kind video")
		}
		if err := mustInstagramFetchableMediaURL(items[0].URL); err != nil {
			return out, err
		}
		out.VideoURL = items[0].URL
	case "carousel":
		if len(items) < 2 || len(items) > 10 {
			return out, fmt.Errorf("carousel requires 2–10 media items")
		}
		for _, it := range items {
			if it.Kind != "image" {
				return out, fmt.Errorf("carousel items must be kind image")
			}
			if err := mustInstagramFetchableMediaURL(it.URL); err != nil {
				return out, err
			}
			out.CarouselURLs = append(out.CarouselURLs, it.URL)
		}
	default:
		return out, fmt.Errorf("unsupported post_type %q", sp.PostType)
	}
	return out, nil
}

func mustHTTPSURL(s string) error {
	u, err := url.Parse(s)
	if err != nil || u.Scheme != "https" || u.Host == "" {
		return fmt.Errorf("media url must be a valid https URL")
	}
	return nil
}

func mustInstagramFetchableMediaURL(s string) error {
	if err := mustHTTPSURL(s); err != nil {
		return err
	}
	u, err := url.Parse(s)
	if err != nil {
		return err
	}
	host := strings.ToLower(u.Hostname())
	path := strings.ToLower(u.Path)
	q := strings.ToLower(u.RawQuery)
	switch {
	case strings.Contains(path, "/media/stream/"):
		return nil
	case strings.Contains(host, "drive.google.com"):
		if strings.Contains(path, "/file/d/") || strings.Contains(path, "/file/u/") ||
			strings.HasPrefix(path, "/open") {
			return fmt.Errorf("Google Drive share or preview links return a web page, not a JPEG file; Instagram cannot use them")
		}
	case host == "docs.google.com":
		return fmt.Errorf("Google Docs URLs are not direct media links")
	case strings.HasSuffix(host, ".dropbox.com") || host == "dropbox.com":
		if strings.Contains(path, "/s/") && !strings.Contains(q, "raw=1") && !strings.Contains(q, "dl=1") {
			return fmt.Errorf("Dropbox shared links often return HTML; append ?raw=1 or dl=1")
		}
	}
	return nil
}
