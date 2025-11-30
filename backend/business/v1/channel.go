package businessv1

import (
	"context"

	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/services/instagram_service"
)

func CreateInstagramChannel(ctx context.Context, req modelsv1.CreateInstagramChannelRequest) (*modelsv1.CreateInstagramChannelResponse, error) {
	service := instagram_service.NewInstagramService()

	channelResp, err := service.CreateChannel(
		ctx,
		req.Code,
		req.ChannelName,
		req.OwnerUserID,
		req.TeamID,
	)
	if err != nil {
		return nil, err
	}

	return &modelsv1.CreateInstagramChannelResponse{
		ChannelID:   channelResp.ChannelID,
		ChannelName: channelResp.ChannelName,
	}, nil
}
