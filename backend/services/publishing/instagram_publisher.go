package publishing

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/shekhar8352/PostEaze/provider/instagram"
)

// InstagramPublisher implements Publisher using Instagram Graph content publishing.
type InstagramPublisher struct {
	API instagram.InstagramProvider
}

func NewInstagramPublisher(api instagram.InstagramProvider) *InstagramPublisher {
	return &InstagramPublisher{API: api}
}

func (p *InstagramPublisher) Platform() string { return PlatformInstagram }

func (p *InstagramPublisher) Schedule(ctx context.Context, accessToken, igUserID string, payload SchedulePayload) (creationID, publishedMediaID string, err error) {
	if p.API == nil {
		return "", "", fmt.Errorf("instagram API not configured")
	}
	var schedPtr *time.Time
	if !payload.PublishNow {
		t := payload.ScheduledAt
		schedPtr = &t
	}
	switch payload.PostType {
	case "image":
		req := &instagram.ContentPublishRequest{
			PostType:    instagram.ContentTypeImage,
			Caption:     payload.Caption,
			ScheduledAt: schedPtr,
			ImageURL:    payload.ImageURL,
		}
		creationID, err = p.API.CreateMediaContainer(ctx, accessToken, igUserID, req)
	case "video":
		req := &instagram.ContentPublishRequest{
			PostType:    instagram.ContentTypeVideo,
			Caption:     payload.Caption,
			ScheduledAt: schedPtr,
			VideoURL:    payload.VideoURL,
		}
		creationID, err = p.API.CreateMediaContainer(ctx, accessToken, igUserID, req)
	case "carousel":
		creationID, err = p.API.CreateCarouselContainers(ctx, accessToken, igUserID, payload.CarouselURLs, payload.Caption, schedPtr)
	default:
		return "", "", fmt.Errorf("unsupported post type for instagram: %s", payload.PostType)
	}
	if err != nil {
		return "", "", err
	}

	// Meta must finish ingesting remote media (status FINISHED) before media_publish.
	// We previously skipped this for scheduled posts, which led to immediate media_publish
	// while the container was still IN_PROGRESS and opaque errors (e.g. code=1).
	if err := p.waitForContainerReady(ctx, accessToken, creationID); err != nil {
		return creationID, "", err
	}
	publishedMediaID, err = p.API.PublishMedia(ctx, accessToken, igUserID, creationID)
	if err != nil && payload.PublishNow && isContainerNotReadyError(err) {
		if waitErr := p.waitForContainerReady(ctx, accessToken, creationID); waitErr != nil {
			return creationID, "", waitErr
		}
		publishedMediaID, err = p.API.PublishMedia(ctx, accessToken, igUserID, creationID)
	}
	if err != nil {
		return creationID, "", err
	}
	return creationID, publishedMediaID, nil
}

func (p *InstagramPublisher) waitForContainerReady(ctx context.Context, accessToken, creationID string) error {
	const (
		maxWait      = 90 * time.Second
		pollInterval = 3 * time.Second
	)
	deadline := time.Now().Add(maxWait)
	var lastStatus string
	for time.Now().Before(deadline) {
		status, err := p.API.GetMediaContainerStatus(ctx, accessToken, creationID)
		if err != nil {
			return err
		}
		code := strings.ToUpper(strings.TrimSpace(status.StatusCode))
		lastStatus = code
		switch code {
		case "FINISHED", "PUBLISHED":
			return nil
		case "ERROR", "EXPIRED":
			return fmt.Errorf("instagram container is not publishable: status_code=%s status=%s", status.StatusCode, status.Status)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(pollInterval):
		}
	}
	return fmt.Errorf("timed out waiting for instagram container readiness (last_status=%s)", lastStatus)
}

func isContainerNotReadyError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "code=9007") ||
		strings.Contains(msg, "error_subcode=2207027") ||
		strings.Contains(msg, "media id is not available")
}
