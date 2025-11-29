package modelsv1

import "github.com/google/uuid"

type CreateInstagramChannelRequest struct {
	Code        string     `json:"code" binding:"required"`
	ChannelName string     `json:"channel_name" binding:"required"`
	OwnerUserID uuid.UUID  `json:"owner_user_id" binding:"required"`
	TeamID      *uuid.UUID `json:"team_id"`
}

type CreateInstagramChannelResponse struct {
	ChannelID   int64  `json:"channel_id"`
	ChannelName string `json:"channel_name"`
}
