package businessv1

import (
	"context"

	"github.com/google/uuid"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/services/instagram_service"
)

func CreateInstagramChannel(ctx context.Context, req modelsv1.CreateInstagramChannelRequest, userID string) (*modelsv1.CreateInstagramChannelResponse, error) {
	service := instagram_service.NewInstagramService()

	// Parse user ID string to UUID
	ownerUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	channelResp, err := service.CreateChannel(
		ctx,
		req.Code,
		req.ChannelName,
		ownerUserID,
		req.TeamID,
		req.Metadata,
	)
	if err != nil {
		return nil, err
	}

	return &modelsv1.CreateInstagramChannelResponse{
		ChannelID:   channelResp.ChannelID,
		ChannelName: channelResp.ChannelName,
	}, nil
}
