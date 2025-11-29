package entities

import (
	"time"

	"github.com/google/uuid"
)

type Channel struct {
	ID                int64      `db:"id"`
	OwnerUserID       uuid.UUID  `db:"owner_user_id"`
	TeamID            *uuid.UUID `db:"team_id"`
	Provider          string     `db:"provider"`
	ProviderChannelID string     `db:"provider_channel_id"`
	DisplayName       *string    `db:"display_name"`
	Username          *string    `db:"username"`
	AvatarURL         *string    `db:"avatar_url"`
	IsActive          bool       `db:"is_active"`
	ErrorStatus       *string    `db:"error_status"`
	Metadata          []byte     `db:"metadata"` // JSONB stored as bytes
	ConnectedAt       time.Time  `db:"connected_at"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at"`
}

type ChannelToken struct {
	ID           int64      `db:"id"`
	ChannelID    int64      `db:"channel_id"`
	AccessToken  []byte     `db:"access_token"`  // Encrypted BYTEA
	RefreshToken []byte     `db:"refresh_token"` // Encrypted BYTEA
	TokenType    *string    `db:"token_type"`
	Scopes       *string    `db:"scopes"`
	ExpiresAt    *time.Time `db:"expires_at"`
	Revoked      bool       `db:"revoked"`
	IssuedAt     time.Time  `db:"issued_at"`
	LastUsedAt   *time.Time `db:"last_used_at"`
	CreatedAt    time.Time  `db:"created_at"`
}
