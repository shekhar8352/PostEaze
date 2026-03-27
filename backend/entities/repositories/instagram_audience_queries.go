package repositories

import (
	"context"
	"database/sql"

	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/utils/database"
)

// GetLatestInstagramAudienceSnapshot returns the newest stored audience payload for a channel, or nil if none.
func GetLatestInstagramAudienceSnapshot(ctx context.Context, channelID int64) (*entities.InstagramAudienceSnapshot, error) {
	db := database.GetDB()
	const q = `
		SELECT id, channel_id, snapshot_date, raw, created_at
		FROM instagram_audience_snapshots
		WHERE channel_id = $1
		ORDER BY snapshot_date DESC, created_at DESC
		LIMIT 1
	`
	var s entities.InstagramAudienceSnapshot
	err := db.QueryRowContext(ctx, q, channelID).Scan(
		&s.ID, &s.ChannelID, &s.SnapshotDate, &s.Raw, &s.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}
