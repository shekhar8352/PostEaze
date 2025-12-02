package businessv1

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/shekhar8352/PostEaze/entities/repositories"
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

func GetChannels(ctx context.Context, userID string, provider string) (*modelsv1.GetChannelsResponse, error) {
	channels, err := repositories.GetChannelsByUserID(ctx, userID, provider)
	if err != nil {
		return nil, err
	}

	channelInfos := make([]modelsv1.ChannelInfo, 0, len(channels))
	for _, ch := range channels {
		// Unmarshal metadata
		var metadata map[string]interface{}
		if len(ch.Metadata) > 0 {
			if err := json.Unmarshal(ch.Metadata, &metadata); err != nil {
				// If unmarshal fails, use empty map
				metadata = make(map[string]interface{})
			}
		} else {
			metadata = make(map[string]interface{})
		}

		channelName := ""
		if ch.DisplayName != nil {
			channelName = *ch.DisplayName
		}

		channelInfos = append(channelInfos, modelsv1.ChannelInfo{
			ChannelID:         ch.ID,
			ChannelName:       channelName,
			Provider:          ch.Provider,
			ProviderChannelID: ch.ProviderChannelID,
			IsActive:          ch.IsActive,
			Metadata:          metadata,
			CreatedAt:         ch.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return &modelsv1.GetChannelsResponse{
		Channels: channelInfos,
		Total:    len(channelInfos),
	}, nil
}
