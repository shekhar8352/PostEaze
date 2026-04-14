package instagram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
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

// ContainerStatusResponse is returned by GET /{creation-id}?fields=status_code,status.
type ContainerStatusResponse struct {
	StatusCode string `json:"status_code"`
	Status     string `json:"status"`
}

func graphAPIBase() string {
	if b := os.Getenv("INSTAGRAM_GRAPH_BASE_URL"); b != "" {
		return strings.TrimRight(b, "/")
	}
	// Default aligns with Meta’s Content Publishing examples (Instagram Login host).
	return "https://graph.instagram.com/v25.0"
}

// graphPOSTJSON sends POST requests the way Meta documents them for content publishing:
// Content-Type: application/json and Authorization: Bearer <token>.
// See https://developers.facebook.com/docs/instagram-platform/content-publishing/
func (p *InstagramProviderImpl) graphPOSTJSON(ctx context.Context, accessToken, endpoint string, body map[string]any) ([]byte, error) {
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		if ge, ok := GraphAPIErrorFromBody(respBody); ok {
			return nil, ge
		}
		return nil, fmt.Errorf("instagram graph POST %s: %s: %s", endpoint, resp.Status, string(respBody))
	}
	return respBody, nil
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

	body := map[string]any{
		"media_type": "CAROUSEL",
		"children":   strings.Join(childIDs, ","),
	}
	if caption != "" {
		body["caption"] = caption
	}
	if scheduledAt != nil {
		body["scheduled_publish_time"] = scheduledAt.UTC().Unix()
	}

	endpoint := fmt.Sprintf("%s/%s/media", graphAPIBase(), igUserID)
	return p.parseCreateMediaResponse(p.graphPOSTJSON(ctx, accessToken, endpoint, body))
}

func (p *InstagramProviderImpl) createCarouselChild(ctx context.Context, accessToken, igUserID, imageURL string) (string, error) {
	body := map[string]any{
		"image_url":         imageURL,
		"is_carousel_item":  true,
	}
	endpoint := fmt.Sprintf("%s/%s/media", graphAPIBase(), igUserID)
	return p.parseCreateMediaResponse(p.graphPOSTJSON(ctx, accessToken, endpoint, body))
}

func (p *InstagramProviderImpl) createImageOrVideoContainer(ctx context.Context, accessToken, igUserID, mediaType, imageURL, videoURL, caption string, scheduledAt *time.Time) (string, error) {
	body := make(map[string]any)
	switch mediaType {
	case "IMAGE":
		if imageURL == "" {
			return "", fmt.Errorf("image_url required")
		}
		body["image_url"] = imageURL
	case "VIDEO":
		if videoURL == "" {
			return "", fmt.Errorf("video_url required")
		}
		body["media_type"] = "VIDEO"
		body["video_url"] = videoURL
	default:
		return "", fmt.Errorf("unsupported mediaType %q", mediaType)
	}
	if caption != "" {
		body["caption"] = caption
	}
	if scheduledAt != nil {
		body["scheduled_publish_time"] = scheduledAt.UTC().Unix()
	}
	endpoint := fmt.Sprintf("%s/%s/media", graphAPIBase(), igUserID)
	return p.parseCreateMediaResponse(p.graphPOSTJSON(ctx, accessToken, endpoint, body))
}

func (p *InstagramProviderImpl) parseCreateMediaResponse(respBody []byte, err error) (string, error) {
	if err != nil {
		return "", err
	}
	var out CreateMediaContainerResponse
	if err := json.Unmarshal(respBody, &out); err != nil {
		return "", fmt.Errorf("parse create media: %w", err)
	}
	if out.ID == "" {
		return "", fmt.Errorf("empty creation id in response: %s", string(respBody))
	}
	return out.ID, nil
}

// PublishMedia submits media_publish for a container id (required after create, including scheduled).
func (p *InstagramProviderImpl) PublishMedia(ctx context.Context, accessToken, igUserID, creationID string) (publishedMediaID string, err error) {
	endpoint := fmt.Sprintf("%s/%s/media_publish", graphAPIBase(), igUserID)
	body := map[string]any{"creation_id": creationID}
	respBody, err := p.graphPOSTJSON(ctx, accessToken, endpoint, body)
	if err != nil {
		return "", err
	}
	var out PublishMediaResponse
	if err := json.Unmarshal(respBody, &out); err != nil {
		return "", fmt.Errorf("parse media_publish: %w", err)
	}
	return out.ID, nil
}

// GetMediaContainerStatus fetches processing status for a media container.
func (p *InstagramProviderImpl) GetMediaContainerStatus(ctx context.Context, accessToken, creationID string) (ContainerStatusResponse, error) {
	v := url.Values{}
	v.Set("fields", "status_code,status")
	v.Set("access_token", accessToken)
	endpoint := fmt.Sprintf("%s/%s?%s", graphAPIBase(), creationID, v.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return ContainerStatusResponse{}, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ContainerStatusResponse{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ContainerStatusResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		if ge, ok := GraphAPIErrorFromBody(body); ok {
			return ContainerStatusResponse{}, ge
		}
		return ContainerStatusResponse{}, fmt.Errorf("get media container status: %s: %s", resp.Status, string(body))
	}
	var out ContainerStatusResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return ContainerStatusResponse{}, fmt.Errorf("parse media container status: %w", err)
	}
	return out, nil
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
