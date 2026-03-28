package instagram

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// ContentPublishPostType matches scheduling payloads (feed types).
type ContentPublishPostType string

const (
	ContentTypeImage    ContentPublishPostType = "image"
	ContentTypeVideo    ContentPublishPostType = "video"
	ContentTypeCarousel ContentPublishPostType = "carousel"
)

// ContentPublishRequest is input for creating a scheduled or immediate IG media container.
type ContentPublishRequest struct {
	PostType          ContentPublishPostType
	Caption           string
	ScheduledAt       *time.Time // if set, sends scheduled_publish_time (Unix UTC)
	ImageURL          string     // single image or ignored for carousel/video
	VideoURL          string
	CarouselImageURLs []string   // 2–10 images for carousel
}

// CreateMediaContainerResponse is returned by POST /{ig-user-id}/media.
type CreateMediaContainerResponse struct {
	ID string `json:"id"`
}

// PublishMediaResponse is returned by POST /{ig-user-id}/media_publish.
type PublishMediaResponse struct {
	ID string `json:"id"`
}

func graphAPIBase() string {
	if b := os.Getenv("INSTAGRAM_GRAPH_BASE_URL"); b != "" {
		return strings.TrimRight(b, "/")
	}
	return "https://graph.instagram.com/v24.0"
}

// CreateMediaContainer creates one IG media container (image, video, or carousel parent/child).
// For carousel, call CreateCarouselContainers instead.
func (p *InstagramProviderImpl) CreateMediaContainer(ctx context.Context, accessToken, igUserID string, req *ContentPublishRequest) (creationID string, err error) {
	if req == nil {
		return "", fmt.Errorf("nil request")
	}
	switch req.PostType {
	case ContentTypeImage:
		return p.createImageOrVideoContainer(ctx, accessToken, igUserID, "IMAGE", req.ImageURL, "", req.Caption, req.ScheduledAt)
	case ContentTypeVideo:
		return p.createImageOrVideoContainer(ctx, accessToken, igUserID, "VIDEO", "", req.VideoURL, req.Caption, req.ScheduledAt)
	case ContentTypeCarousel:
		return "", fmt.Errorf("use CreateCarouselContainers for carousel")
	default:
		return "", fmt.Errorf("unsupported post type: %s", req.PostType)
	}
}

// CreateCarouselContainers creates child image containers then a CAROUSEL parent; returns parent creation id.
func (p *InstagramProviderImpl) CreateCarouselContainers(ctx context.Context, accessToken, igUserID string, imageURLs []string, caption string, scheduledAt *time.Time) (parentCreationID string, err error) {
	if len(imageURLs) < 2 || len(imageURLs) > 10 {
		return "", fmt.Errorf("carousel requires 2–10 images, got %d", len(imageURLs))
	}
	childIDs := make([]string, 0, len(imageURLs))
	for _, u := range imageURLs {
		cid, err := p.createCarouselChild(ctx, accessToken, igUserID, u)
		if err != nil {
			return "", err
		}
		childIDs = append(childIDs, cid)
	}

	v := url.Values{}
	v.Set("access_token", accessToken)
	v.Set("media_type", "CAROUSEL")
	v.Set("children", strings.Join(childIDs, ","))
	if caption != "" {
		v.Set("caption", caption)
	}
	if scheduledAt != nil {
		v.Set("scheduled_publish_time", strconv.FormatInt(scheduledAt.UTC().Unix(), 10))
	}

	endpoint := fmt.Sprintf("%s/%s/media", graphAPIBase(), igUserID)
	return p.postMediaContainer(ctx, endpoint, v)
}

func (p *InstagramProviderImpl) createCarouselChild(ctx context.Context, accessToken, igUserID, imageURL string) (string, error) {
	v := url.Values{}
	v.Set("access_token", accessToken)
	v.Set("image_url", imageURL)
	v.Set("is_carousel_item", "true")
	endpoint := fmt.Sprintf("%s/%s/media", graphAPIBase(), igUserID)
	return p.postMediaContainer(ctx, endpoint, v)
}

func (p *InstagramProviderImpl) createImageOrVideoContainer(ctx context.Context, accessToken, igUserID, mediaType, imageURL, videoURL, caption string, scheduledAt *time.Time) (string, error) {
	v := url.Values{}
	v.Set("access_token", accessToken)
	if mediaType == "IMAGE" {
		if imageURL == "" {
			return "", fmt.Errorf("image_url required")
		}
		v.Set("image_url", imageURL)
	} else if mediaType == "VIDEO" {
		if videoURL == "" {
			return "", fmt.Errorf("video_url required")
		}
		v.Set("media_type", "VIDEO")
		v.Set("video_url", videoURL)
	}
	if caption != "" {
		v.Set("caption", caption)
	}
	if scheduledAt != nil {
		v.Set("scheduled_publish_time", strconv.FormatInt(scheduledAt.UTC().Unix(), 10))
	}
	endpoint := fmt.Sprintf("%s/%s/media", graphAPIBase(), igUserID)
	return p.postMediaContainer(ctx, endpoint, v)
}

func (p *InstagramProviderImpl) postMediaContainer(ctx context.Context, endpoint string, form url.Values) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		if ge, ok := GraphAPIErrorFromBody(body); ok {
			return "", ge
		}
		return "", fmt.Errorf("create media container: %s: %s", resp.Status, string(body))
	}
	var out CreateMediaContainerResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return "", fmt.Errorf("parse create media: %w", err)
	}
	if out.ID == "" {
		return "", fmt.Errorf("empty creation id in response: %s", string(body))
	}
	return out.ID, nil
}

// PublishMedia submits media_publish for a container id (required after create, including scheduled).
func (p *InstagramProviderImpl) PublishMedia(ctx context.Context, accessToken, igUserID, creationID string) (publishedMediaID string, err error) {
	v := url.Values{}
	v.Set("access_token", accessToken)
	v.Set("creation_id", creationID)
	endpoint := fmt.Sprintf("%s/%s/media_publish", graphAPIBase(), igUserID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(v.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		if ge, ok := GraphAPIErrorFromBody(body); ok {
			return "", ge
		}
		return "", fmt.Errorf("media_publish: %s: %s", resp.Status, string(body))
	}
	var out PublishMediaResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return "", fmt.Errorf("parse media_publish: %w", err)
	}
	return out.ID, nil
}

// ValidateScheduledPublishWindow checks Meta’s documented window (10 min – 75 days). Adjust if Meta changes.
func ValidateScheduledPublishWindow(scheduledUTC time.Time, now time.Time) error {
	min := now.Add(10 * time.Minute)
	max := now.Add(75 * 24 * time.Hour)
	if scheduledUTC.Before(min) || scheduledUTC.After(max) {
		return fmt.Errorf("scheduled_at must be between 10 minutes and 75 days from now (UTC)")
	}
	return nil
}
