package entities

import (
	"time"

	"github.com/google/uuid"
)

type MediaAssetStatus string

const (
	MediaAssetStatusDraft     MediaAssetStatus = "draft"
	MediaAssetStatusReady     MediaAssetStatus = "ready"
	MediaAssetStatusPublished MediaAssetStatus = "published"
)

type MediaAsset struct {
	ID               int64     `db:"id"`
	OwnerUserID      uuid.UUID `db:"owner_user_id"`
	TeamID           *uuid.UUID `db:"team_id"`
	Title            string    `db:"title"`
	AssetType        string    `db:"asset_type"` // "photo" | "video"
	Status           string    `db:"status"`
	CurrentVersionID *int64    `db:"current_version_id"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}
