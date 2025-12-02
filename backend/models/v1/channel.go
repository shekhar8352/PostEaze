package modelsv1

import "github.com/google/uuid"

type CreateInstagramChannelRequest struct {
	Code        string                 `json:"code" binding:"required"`
	ChannelName string                 `json:"channel_name" binding:"required"`
	TeamID      *uuid.UUID             `json:"team_id"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type CreateInstagramChannelResponse struct {
	ChannelID   int64  `json:"channel_id"`
	ChannelName string `json:"channel_name"`
}

type GetChannelsRequest struct {
	Provider string `form:"provider"` // Optional filter by provider (instagram, facebook, etc.)
}

type ChannelInfo struct {
	ChannelID         int64                  `json:"channel_id"`
	ChannelName       string                 `json:"channel_name"`
	Provider          string                 `json:"provider"`
	ProviderChannelID string                 `json:"provider_channel_id"`
	IsActive          bool                   `json:"is_active"`
	Metadata          map[string]interface{} `json:"metadata"`
	CreatedAt         string                 `json:"created_at"`
}

type GetChannelsResponse struct {
	Channels []ChannelInfo `json:"channels"`
	Total    int           `json:"total"`
}

type GetPageDetailsRequest struct {
	ChannelID int64 `form:"channel_id" binding:"required"`
}

type GetPageDetailsResponse struct {
	ID                string `json:"id"`
	Username          string `json:"username"`
	Name              string `json:"name"`
	Biography         string `json:"biography"`
	FollowersCount    int    `json:"followers_count"`
	FollowsCount      int    `json:"follows_count"`
	MediaCount        int    `json:"media_count"`
	ProfilePictureURL string `json:"profile_picture_url"`
	Website           string `json:"website"`
}
