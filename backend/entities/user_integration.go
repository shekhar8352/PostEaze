package entities

import (
	"time"

	"github.com/google/uuid"
)

const UserIntegrationProviderGoogleDrive = "google_drive"

type UserIntegration struct {
	ID                    int64      `db:"id"`
	UserID                uuid.UUID  `db:"user_id"`
	Provider              string     `db:"provider"`
	ProviderAccountEmail  *string    `db:"provider_account_email"`
	AccessToken           []byte     `db:"access_token"`
	RefreshToken          []byte     `db:"refresh_token"`
	Scopes                *string    `db:"scopes"`
	ExpiresAt             *time.Time `db:"expires_at"`
	Revoked               bool       `db:"revoked"`
	CreatedAt             time.Time  `db:"created_at"`
	UpdatedAt             time.Time  `db:"updated_at"`
}
