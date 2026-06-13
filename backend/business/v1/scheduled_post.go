package businessv1

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/entities/repositories"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/provider/instagram"
	"github.com/shekhar8352/PostEaze/mediapublish"
	"github.com/shekhar8352/PostEaze/services/publishing"
	"github.com/shekhar8352/PostEaze/tasks"
	"github.com/shekhar8352/PostEaze/utils/encryption"
)

var defaultInstagramPublisher = publishing.NewInstagramPublisher(instagram.NewInstagramProvider())

// CreateScheduledPost validates input, persists, calls Meta per channel, updates row.
func CreateScheduledPost(ctx context.Context, userIDStr string, req *modelsv1.CreateScheduledPostRequest) (*modelsv1.CreateScheduledPostResponse, int, error) {
	ownerID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, 400, fmt.Errorf("invalid user")
	}
	if err := validatePlatforms(req.Platforms); err != nil {
		return nil, 400, err
	}
	if req.PostType == "reel" || req.PostType == "story" {
		return nil, 501, fmt.Errorf("post type %q is not supported yet; use image, video, or carousel", req.PostType)
	}
	var schedUTC time.Time
	if req.PublishNow {
		schedUTC = time.Now().UTC()
	} else {
		if req.ScheduledAt == nil {
			return nil, 400, fmt.Errorf("scheduled_at is required when publish_now is false")
		}
		schedUTC = req.ScheduledAt.UTC()
		if err := instagram.ValidateScheduledPublishWindow(schedUTC, time.Now().UTC()); err != nil {
			return nil, 400, err
		}
	}
	if err := resolveScheduledMediaItems(ctx, ownerID, &req.Media, req.PostType); err != nil {
		return nil, 400, err
	}
	payload, err := buildSchedulePayload(req)
	if err != nil {
		return nil, 400, err
	}
	igChannels, ytChannels, err := validateAndPartitionChannels(ctx, userIDStr, req)
	if err != nil {
		return nil, 400, err
	}
	if len(igChannels) == 0 && len(ytChannels) == 0 {
		return nil, 400, fmt.Errorf("no valid channels for selected platforms")
	}
	if len(ytChannels) > 0 && req.PostType != "video" {
		return nil, 400, fmt.Errorf("youtube only supports video post type")
	}

	mediaJSON, err := json.Marshal(req.Media)
	if err != nil {
		return nil, 500, err
	}
	stateJSON := []byte(`{}`)

	initialStatus := string(entities.ScheduledStatusPending)
	if !req.PublishNow {
		initialStatus = string(entities.ScheduledStatusScheduled)
	}
	sp := &entities.ScheduledPost{
		OwnerUserID:   ownerID,
		ChannelIDs:    pq.Int64Array(append([]int64(nil), req.ChannelIDs...)),
		Platforms:     pq.StringArray(append([]string(nil), req.Platforms...)),
		ScheduledAt:   schedUTC,
		Status:        initialStatus,
		PostType:      req.PostType,
		Caption:       strPtrOrNil(req.Caption),
		Media:         mediaJSON,
		ProviderState: stateJSON,
	}
	if err := repositories.CreateScheduledPost(ctx, sp); err != nil {
		return nil, 500, err
	}

	// Optionally auto-link to a Studio Piece. Linking failures are logged via
	// the business layer's activity pipeline but never fail the schedule call;
	// the scheduled post has already been persisted at this point.
	if req.PieceID != nil && *req.PieceID > 0 {
		if _, err := LinkPieceScheduledPost(ctx, userIDStr, *req.PieceID, &modelsv1.LinkScheduledPostRequest{
			ScheduledPostID: sp.ID,
		}); err != nil {
			// Swallow but continue – the caller can retry linking from the UI.
			_ = err
		}
	}

	// Future posts: queue Asynq jobs at scheduled_at.
	if !req.PublishNow {
		stateMap := map[string]any{
			"publish_now": false,
			"queue":       "asynq",
			"run_at_utc":  schedUTC.Format(time.RFC3339),
			"instagram":   map[string]any{"channels": map[string]any{}},
			"youtube":     map[string]any{"channels": map[string]any{}},
		}
		results := make([]modelsv1.ChannelScheduleResult, 0, len(req.ChannelIDs))
		if igInner, ok := stateMap["instagram"].(map[string]any); ok {
			if inner, ok := igInner["channels"].(map[string]any); ok {
				for _, cid := range igChannels {
					results = append(results, modelsv1.ChannelScheduleResult{ChannelID: cid, Success: true})
					inner[fmt.Sprintf("%d", cid)] = map[string]any{"status": "queued"}
				}
			}
		}
		if ytInner, ok := stateMap["youtube"].(map[string]any); ok {
			if inner, ok := ytInner["channels"].(map[string]any); ok {
				for _, cid := range ytChannels {
					results = append(results, modelsv1.ChannelScheduleResult{ChannelID: cid, Success: true})
					inner[fmt.Sprintf("%d", cid)] = map[string]any{"status": "queued"}
				}
			}
		}
		stateBytes, _ := json.Marshal(stateMap)
		if err := repositories.UpdateScheduledPostStatusAndProviderState(ctx, sp.ID, ownerID, string(entities.ScheduledStatusScheduled), stateBytes); err != nil {
			return nil, 500, err
		}
		if len(igChannels) > 0 {
			if _, err := tasks.EnqueueInstagramScheduledPostPublish(sp.ID, schedUTC); err != nil {
				return nil, 500, fmt.Errorf("queue instagram job: %w", err)
			}
		}
		if len(ytChannels) > 0 {
			if _, err := tasks.EnqueueYouTubeScheduledPostPublish(sp.ID, schedUTC); err != nil {
				return nil, 500, fmt.Errorf("queue youtube job: %w", err)
			}
		}
		return &modelsv1.CreateScheduledPostResponse{
			ScheduledPostID: sp.ID,
			OverallStatus:   "scheduled",
			PublishNow:      false,
			Results:         results,
		}, 200, nil
	}

	_ = repositories.UpdateScheduledPostStatusAndProviderState(ctx, sp.ID, ownerID, string(entities.ScheduledStatusSubmitting), stateJSON)

	results := make([]modelsv1.ChannelScheduleResult, 0, len(req.ChannelIDs))
	igPub := defaultInstagramPublisher
	stateMap := map[string]any{
		"publish_now": req.PublishNow,
		"instagram":   map[string]any{"channels": map[string]any{}},
		"youtube":     map[string]any{"channels": map[string]any{}},
	}
	channelsMap, _ := stateMap["instagram"].(map[string]any)
	inner, _ := channelsMap["channels"].(map[string]any)
	ytChannelsMap, _ := stateMap["youtube"].(map[string]any)
	ytInner, _ := ytChannelsMap["channels"].(map[string]any)

	igSuccessN := 0
	for _, cid := range igChannels {
		res := modelsv1.ChannelScheduleResult{ChannelID: cid}
		tok, err := repositories.GetLatestTokenByChannelID(ctx, cid)
		if err != nil {
			res.ErrorMessage = "token: " + err.Error()
			results = append(results, res)
			inner[fmt.Sprintf("%d", cid)] = map[string]any{"error": res.ErrorMessage}
			continue
		}
		plain, err := encryption.Decrypt(tok.AccessToken)
		if err != nil {
			res.ErrorMessage = "decrypt token: " + err.Error()
			results = append(results, res)
			inner[fmt.Sprintf("%d", cid)] = map[string]any{"error": res.ErrorMessage}
			continue
		}
		ch, err := repositories.GetChannelByID(ctx, cid)
		if err != nil {
			res.ErrorMessage = "channel: " + err.Error()
			results = append(results, res)
			inner[fmt.Sprintf("%d", cid)] = map[string]any{"error": res.ErrorMessage}
			continue
		}
		igUserID, err := instagram.UserIDFromChannelMetadata(ch.Metadata)
		if err != nil {
			res.ErrorMessage = err.Error()
			results = append(results, res)
			inner[fmt.Sprintf("%d", cid)] = map[string]any{"error": res.ErrorMessage}
			continue
		}

		pubPayload := publishing.SchedulePayload{
			PostType:     payload.PostType,
			Caption:      payload.Caption,
			PublishNow:   true,
			ScheduledAt:  schedUTC,
			ImageURL:     payload.ImageURL,
			VideoURL:     payload.VideoURL,
			CarouselURLs: payload.CarouselURLs,
		}
		creationID, publishedID, schedErr := igPub.Schedule(ctx, plain, igUserID, pubPayload)
		if schedErr != nil {
			res.ErrorMessage = schedErr.Error()
			inner[fmt.Sprintf("%d", cid)] = map[string]any{"error": res.ErrorMessage}
		} else {
			res.Success = true
			res.CreationID = creationID
			res.PublishedID = publishedID
			igSuccessN++
			inner[fmt.Sprintf("%d", cid)] = map[string]any{
				"creation_id":          creationID,
				"published_media_id": publishedID,
			}
		}
		results = append(results, res)
	}

	ytQueued := 0
	for _, cid := range ytChannels {
		res := modelsv1.ChannelScheduleResult{ChannelID: cid}
		ytInner[fmt.Sprintf("%d", cid)] = map[string]any{"status": "queued"}
		res.Success = true
		ytQueued++
		results = append(results, res)
	}
	if len(ytChannels) > 0 {
		if _, err := tasks.EnqueueYouTubeScheduledPostPublish(sp.ID, time.Now().UTC()); err != nil {
			for _, cid := range ytChannels {
				ytInner[fmt.Sprintf("%d", cid)] = map[string]any{"error": err.Error()}
			}
		}
	}

	stateBytes, _ := json.Marshal(stateMap)
	overall := "failed"
	finalStatus := string(entities.ScheduledStatusFailed)
	igDone := len(igChannels) == 0 || igSuccessN == len(igChannels)
	ytDone := len(ytChannels) == 0 || ytQueued == len(ytChannels)
	if igDone && ytDone && len(ytChannels) == 0 {
		if igSuccessN == len(igChannels) {
			overall = "published"
			finalStatus = string(entities.ScheduledStatusPublished)
		} else if igSuccessN > 0 {
			overall = "partial_failure"
		}
	} else if igDone && len(ytChannels) > 0 {
		overall = "submitting"
		finalStatus = string(entities.ScheduledStatusSubmitting)
	} else if igSuccessN > 0 || ytQueued > 0 {
		overall = "partial_failure"
	}
	_ = repositories.UpdateScheduledPostStatusAndProviderState(ctx, sp.ID, ownerID, finalStatus, stateBytes)

	if finalStatus == string(entities.ScheduledStatusPublished) {
		mediapublish.MarkLinkedMediaAssetsPublished(ctx, ownerID, mediaJSON)
	}

	// Propagate lifecycle to any linked Studio Piece (move to published / log failure).
	_ = OnScheduledPostPublished(ctx, sp.ID, finalStatus == string(entities.ScheduledStatusPublished))

	resp := &modelsv1.CreateScheduledPostResponse{
		ScheduledPostID: sp.ID,
		OverallStatus:   overall,
		PublishNow:      req.PublishNow,
		Results:         results,
	}
	return resp, 200, nil
}

type builtPayload struct {
	PostType     string
	Caption      string
	ImageURL     string
	VideoURL     string
	CarouselURLs []string
}

func buildSchedulePayload(req *modelsv1.CreateScheduledPostRequest) (builtPayload, error) {
	var out builtPayload
	out.PostType = req.PostType
	out.Caption = req.Caption
	items := req.Media.Items
	switch req.PostType {
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
		return out, fmt.Errorf("unsupported post_type %q", req.PostType)
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

// mustInstagramFetchableMediaURL rejects URLs that typically return HTML (viewer pages)
// instead of raw media. Meta cURLs the URL and expects image/video bytes; see
// https://developers.facebook.com/docs/instagram-platform/content-publishing/
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
		// Signed proxy URLs serve raw media bytes for Meta ingestion.
		return nil
	case strings.Contains(host, "drive.google.com"):
		// Viewer/share pages are HTML; Meta cURLs the URL and expects raw JPEG/video bytes.
		if strings.Contains(path, "/file/d/") || strings.Contains(path, "/file/u/") ||
			strings.HasPrefix(path, "/open") {
			return fmt.Errorf("Google Drive share or preview links return a web page, not a JPEG file; Instagram cannot use them. Host the image on a CDN or static HTTPS URL whose response is image/jpeg (see Meta content publishing docs)")
		}
	case host == "docs.google.com":
		return fmt.Errorf("Google Docs URLs are not direct media links; use a public HTTPS URL that returns raw image or video bytes")
	case strings.HasSuffix(host, ".dropbox.com") || host == "dropbox.com":
		if strings.Contains(path, "/s/") && !strings.Contains(q, "raw=1") && !strings.Contains(q, "dl=1") {
			return fmt.Errorf("Dropbox shared folder links often return HTML; append ?raw=1 (or use dl=1) so the URL returns the file bytes")
		}
	}
	return nil
}

func validatePlatforms(platforms []string) error {
	for _, p := range platforms {
		if p != publishing.PlatformInstagram && p != publishing.PlatformYouTube {
			return fmt.Errorf("unsupported platform %q", p)
		}
	}
	return nil
}

func validateAndPartitionChannels(ctx context.Context, userIDStr string, req *modelsv1.CreateScheduledPostRequest) (ig []int64, yt []int64, err error) {
	platformSet := make(map[string]bool)
	for _, p := range req.Platforms {
		platformSet[p] = true
	}
	for _, cid := range req.ChannelIDs {
		ok, ch, e := repositories.UserCanAccessChannel(ctx, cid, userIDStr)
		if e != nil || ch == nil {
			return nil, nil, fmt.Errorf("channel %d not found or inaccessible", cid)
		}
		if !ok {
			return nil, nil, fmt.Errorf("no access to channel %d", cid)
		}
		switch ch.Provider {
		case "instagram":
			if !platformSet[publishing.PlatformInstagram] {
				return nil, nil, fmt.Errorf("channel %d is instagram but instagram not in platforms", cid)
			}
			ig = append(ig, cid)
		case "youtube":
			if !platformSet[publishing.PlatformYouTube] {
				return nil, nil, fmt.Errorf("channel %d is youtube but youtube not in platforms", cid)
			}
			yt = append(yt, cid)
		default:
			return nil, nil, fmt.Errorf("channel %d provider %q is not supported for publishing", cid, ch.Provider)
		}
	}
	return ig, yt, nil
}

func resolveScheduledMediaItems(ctx context.Context, ownerID uuid.UUID, media *modelsv1.ScheduledMediaPayload, postType string) error {
	for i := range media.Items {
		item := &media.Items[i]
		if item.MediaVersionID != nil {
			v, oid, err := repositories.GetMediaVersionWithOwner(ctx, *item.MediaVersionID)
			if err != nil {
				return err
			}
			if v == nil || oid != ownerID {
				return fmt.Errorf("media version %d not found", *item.MediaVersionID)
			}
			url, err := ResolveVersionMediaURL(v)
			if err != nil {
				return err
			}
			item.URL = url
			if item.Kind == "" {
				if strings.HasPrefix(v.ContentType, "video/") {
					item.Kind = "video"
				} else {
					item.Kind = "image"
				}
			}
		}
		if item.URL == "" {
			return fmt.Errorf("media item requires url or media_version_id")
		}
	}
	return nil
}

func strPtrOrNil(s string) *string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return &s
}

// ListScheduledPostsForCalendar returns posts in range for owner.
func ListScheduledPostsForCalendar(ctx context.Context, userIDStr string, q modelsv1.ListScheduledPostsQuery) (*modelsv1.ListScheduledPostsResponse, int, error) {
	ownerID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, 400, fmt.Errorf("invalid user")
	}
	fromT, toT, err := parseCalendarRange(q.From, q.To)
	if err != nil {
		return nil, 400, err
	}
	rows, err := repositories.ListScheduledPostsInRange(ctx, repositories.ScheduledPostRangeFilters{
		OwnerUserID: ownerID,
		From:        fromT,
		To:          toT,
		ChannelID:   q.ChannelID,
	})
	if err != nil {
		return nil, 500, err
	}
	out := make([]modelsv1.ScheduledPostListItem, 0, len(rows))
	for _, sp := range rows {
		out = append(out, mapScheduledPostToListItem(sp))
	}
	return &modelsv1.ListScheduledPostsResponse{Posts: out}, 200, nil
}

func parseCalendarRange(fromStr, toStr string) (time.Time, time.Time, error) {
	fromT, err := parseFlexibleDate(fromStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("from: %w", err)
	}
	toT, err := parseFlexibleDate(toStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("to: %w", err)
	}
	// If date-only (midnight), treat `to` as exclusive end of that calendar day in UTC
	if len(toStr) == 10 {
		toT = toT.Add(24 * time.Hour)
	}
	if !toT.After(fromT) {
		return time.Time{}, time.Time{}, fmt.Errorf("to must be after from")
	}
	return fromT.UTC(), toT.UTC(), nil
}

func parseFlexibleDate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	return time.ParseInLocation("2006-01-02", s, time.UTC)
}

func mapScheduledPostToListItem(sp entities.ScheduledPost) modelsv1.ScheduledPostListItem {
	var media any
	_ = json.Unmarshal(sp.Media, &media)
	var prov any
	_ = json.Unmarshal(sp.ProviderState, &prov)
	caption := ""
	if sp.Caption != nil {
		caption = *sp.Caption
	}
	return modelsv1.ScheduledPostListItem{
		ID:             sp.ID,
		ChannelIDs:     []int64(sp.ChannelIDs),
		Platforms:      []string(sp.Platforms),
		ScheduledAt:    sp.ScheduledAt.UTC().Format(time.RFC3339),
		Status:         sp.Status,
		PostType:       sp.PostType,
		Caption:        caption,
		Media:          media,
		ProviderState:  prov,
		CreatedAt:      sp.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:      sp.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

// GetScheduledPost returns one post for owner.
func GetScheduledPost(ctx context.Context, userIDStr string, id int64) (*modelsv1.ScheduledPostListItem, int, error) {
	ownerID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, 400, fmt.Errorf("invalid user")
	}
	sp, err := repositories.GetScheduledPostByID(ctx, id, ownerID)
	if err != nil {
		return nil, 500, err
	}
	if sp == nil {
		return nil, 404, fmt.Errorf("not found")
	}
	item := mapScheduledPostToListItem(*sp)
	return &item, 200, nil
}

// CancelScheduledPost marks cancelled when allowed.
func CancelScheduledPost(ctx context.Context, userIDStr string, id int64) (int, error) {
	ownerID, err := uuid.Parse(userIDStr)
	if err != nil {
		return 400, fmt.Errorf("invalid user")
	}
	err = repositories.CancelScheduledPost(ctx, id, ownerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 404, fmt.Errorf("not found or cannot cancel")
		}
		return 500, err
	}
	return 200, nil
}
