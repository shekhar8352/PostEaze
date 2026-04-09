package businessv1

import (
	"context"

	"github.com/google/uuid"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/services/meta_service"
)

// CreateFacebookChannel connects a Facebook Page using an OAuth code (single-shot exchange).
func CreateFacebookChannel(ctx context.Context, req modelsv1.CreateFacebookChannelRequest, userID string) (*modelsv1.CreateFacebookChannelResponse, error) {
	ownerUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	svc := meta_service.NewMetaService()
	channelID, name, err := svc.CreateFacebookPageChannel(ctx, req.Code, req.RedirectURI, req.PageID, req.ChannelName, ownerUserID, req.TeamID)
	if err != nil {
		return nil, err
	}
	return &modelsv1.CreateFacebookChannelResponse{
		ChannelID:   channelID,
		ChannelName: name,
	}, nil
}

// CreateFacebookChannelFromPageToken connects a Facebook Page using a page access token (e.g. from POST /meta/callback).
func CreateFacebookChannelFromPageToken(ctx context.Context, req modelsv1.CreateFacebookChannelRequest, userID string) (*modelsv1.CreateFacebookChannelResponse, error) {
	ownerUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	svc := meta_service.NewMetaService()
	channelID, name, err := svc.CreateFacebookPageChannelFromPageAccessToken(ctx, req.PageID, req.PageAccessToken, req.ChannelName, ownerUserID, req.TeamID)
	if err != nil {
		return nil, err
	}
	return &modelsv1.CreateFacebookChannelResponse{
		ChannelID:   channelID,
		ChannelName: name,
	}, nil
}
