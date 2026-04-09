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

// CreateFacebookChannelRequest connects a Facebook Page.
// Use either (1) code + redirect_uri + page_id after a single exchange, or
// (2) page_id + page_access_token from POST /meta/callback (same OAuth session; code is consumed by callback).
type CreateFacebookChannelRequest struct {
	Code            string     `json:"code"`
	RedirectURI     string     `json:"redirect_uri"`
	PageID          string     `json:"page_id"`
	PageAccessToken string     `json:"page_access_token"`
	ChannelName     string     `json:"channel_name"`
	TeamID          *uuid.UUID `json:"team_id"`
}

// CreateFacebookChannelResponse is returned after connecting a Facebook Page.
type CreateFacebookChannelResponse struct {
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

// SubscribeWebhooksRequest for subscribing to Meta webhooks
type SubscribeWebhooksRequest struct {
	ChannelID int64    `json:"channel_id" binding:"required"`
	Fields    []string `json:"fields"` // Optional, defaults to ["comments", "mentions", "story_insights"]
}

// SubscribeWebhooksResponse for webhook subscription result
type SubscribeWebhooksResponse struct {
	Success bool     `json:"success"`
	Fields  []string `json:"subscribed_fields"`
}
