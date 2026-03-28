package publishing

import (
	"context"
	"fmt"

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
	sched := payload.ScheduledAt
	switch payload.PostType {
	case "image":
		req := &instagram.ContentPublishRequest{
			PostType:    instagram.ContentTypeImage,
			Caption:     payload.Caption,
			ScheduledAt: &sched,
			ImageURL:    payload.ImageURL,
		}
		creationID, err = p.API.CreateMediaContainer(ctx, accessToken, igUserID, req)
	case "video":
		req := &instagram.ContentPublishRequest{
			PostType:    instagram.ContentTypeVideo,
			Caption:     payload.Caption,
			ScheduledAt: &sched,
			VideoURL:    payload.VideoURL,
		}
		creationID, err = p.API.CreateMediaContainer(ctx, accessToken, igUserID, req)
	case "carousel":
		creationID, err = p.API.CreateCarouselContainers(ctx, accessToken, igUserID, payload.CarouselURLs, payload.Caption, &sched)
	default:
		return "", "", fmt.Errorf("unsupported post type for instagram: %s", payload.PostType)
	}
	if err != nil {
		return "", "", err
	}
	publishedMediaID, err = p.API.PublishMedia(ctx, accessToken, igUserID, creationID)
	if err != nil {
		return creationID, "", err
	}
	return creationID, publishedMediaID, nil
}
