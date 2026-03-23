package repositories

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/utils/database"
)

// UserCanAccessChannel reports whether the user may use the channel (owner or active team member).
func UserCanAccessChannel(ctx context.Context, channelID int64, userIDStr string) (bool, *entities.Channel, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return false, nil, err
	}

	ch, err := GetChannelByID(ctx, channelID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil, sql.ErrNoRows
		}
		return false, nil, err
	}

	if ch.OwnerUserID == userID {
		return true, ch, nil
	}

	if ch.TeamID == nil {
		return false, ch, nil
	}

	db := database.GetDB()
	var exists bool
	err = db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM team_members
			WHERE team_id = $1 AND user_id = $2 AND status = 'active'
		)
	`, ch.TeamID, userID).Scan(&exists)
	if err != nil {
		return false, ch, err
	}

	return exists, ch, nil
}
