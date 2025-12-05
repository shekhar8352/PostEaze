package businessv1

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/shekhar8352/PostEaze/entities/repositories"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/services/instagram_service"
	"github.com/shekhar8352/PostEaze/utils/encryption"
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

func GetPageDetails(ctx context.Context, channelID int64, userID string) (*modelsv1.GetPageDetailsResponse, error) {
	// 1. Get channel to verify ownership
	channel, err := repositories.GetChannelByID(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("channel not found: %w", err)
	}

	// 2. Verify the channel belongs to the user
	if channel.OwnerUserID.String() != userID {
		return nil, fmt.Errorf("unauthorized: channel does not belong to user")
	}

	// 3. Get the latest access token
	token, err := repositories.GetLatestTokenByChannelID(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// 4. Decrypt the access token
	decryptedToken, err := encryption.Decrypt(token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt access token: %w", err)
	}

	// 5. Call Instagram API to get page details
	service := instagram_service.NewInstagramService()
	pageDetails, err := service.GetPageDetails(ctx, decryptedToken)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page details from Instagram: %w", err)
	}

	return pageDetails, nil
}

// SubscribeToWebhooks subscribes an existing channel to Meta webhooks
func SubscribeToWebhooks(ctx context.Context, channelID int64, userID string, fields []string) (*modelsv1.SubscribeWebhooksResponse, error) {
	// 1. Get channel to verify ownership
	channel, err := repositories.GetChannelByID(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("channel not found: %w", err)
	}

	// 2. Verify the channel belongs to the user
	if channel.OwnerUserID.String() != userID {
		return nil, fmt.Errorf("unauthorized: channel does not belong to user")
	}

	// 3. Get the latest access token
	token, err := repositories.GetLatestTokenByChannelID(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// 4. Decrypt the access token
	decryptedToken, err := encryption.Decrypt(token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt access token: %w", err)
	}

	// 5. Get Instagram user ID from GetPageDetails
	service := instagram_service.NewInstagramService()
	pageDetails, err := service.GetPageDetails(ctx, decryptedToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get page details: %w", err)
	}

	// 6. Set default fields if not provided
	if len(fields) == 0 {
		fields = []string{"comments", "mentions", "story_insights"}
	}

	// 7. Call Instagram service to subscribe to webhooks using the fetched Instagram user ID
	err = service.SubscribeToWebhooks(ctx, decryptedToken, pageDetails.ID, fields)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to webhooks: %w", err)
	}

	return &modelsv1.SubscribeWebhooksResponse{
		Success: true,
		Fields:  fields,
	}, nil
}
