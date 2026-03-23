package repositories

import (
	"context"

	"github.com/shekhar8352/PostEaze/utils/database"
)

// PostBelongsToChannel returns true if the post exists and channel_ids contains channelID.
func PostBelongsToChannel(ctx context.Context, postID, channelID int64) (bool, error) {
	db := database.GetDB()
	var ok bool
	err := db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM posts
			WHERE id = $1 AND channel_ids @> jsonb_build_array($2::bigint)
		)
	`, postID, channelID).Scan(&ok)
	return ok, err
}
