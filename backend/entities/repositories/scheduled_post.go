package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/utils/database"
)

// CreateScheduledPost inserts a new scheduled post row.
func CreateScheduledPost(ctx context.Context, sp *entities.ScheduledPost) error {
	db := database.GetDB()
	q := `
		INSERT INTO scheduled_posts (
			owner_user_id, channel_ids, platforms, scheduled_at, status,
			post_type, caption, media, provider_state
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`
	return db.QueryRowContext(ctx, q,
		sp.OwnerUserID,
		sp.ChannelIDs,
		sp.Platforms,
		sp.ScheduledAt,
		sp.Status,
		sp.PostType,
		sp.Caption,
		sp.Media,
		sp.ProviderState,
	).Scan(&sp.ID, &sp.CreatedAt, &sp.UpdatedAt)
}

// GetScheduledPostByID returns a row if it exists and belongs to ownerUserID.
func GetScheduledPostByID(ctx context.Context, id int64, ownerUserID uuid.UUID) (*entities.ScheduledPost, error) {
	db := database.GetDB()
	q := `
		SELECT id, owner_user_id, channel_ids, platforms, scheduled_at, status,
		       post_type, caption, media, provider_state, created_at, updated_at
		FROM scheduled_posts
		WHERE id = $1 AND owner_user_id = $2
	`
	var sp entities.ScheduledPost
	err := db.QueryRowContext(ctx, q, id, ownerUserID).Scan(
		&sp.ID,
		&sp.OwnerUserID,
		&sp.ChannelIDs,
		&sp.Platforms,
		&sp.ScheduledAt,
		&sp.Status,
		&sp.PostType,
		&sp.Caption,
		&sp.Media,
		&sp.ProviderState,
		&sp.CreatedAt,
		&sp.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &sp, nil
}

// GetScheduledPostByIDForJob loads a scheduled post by primary key (trusted worker / internal use only).
func GetScheduledPostByIDForJob(ctx context.Context, id int64) (*entities.ScheduledPost, error) {
	db := database.GetDB()
	q := `
		SELECT id, owner_user_id, channel_ids, platforms, scheduled_at, status,
		       post_type, caption, media, provider_state, created_at, updated_at
		FROM scheduled_posts
		WHERE id = $1
	`
	var sp entities.ScheduledPost
	err := db.QueryRowContext(ctx, q, id).Scan(
		&sp.ID,
		&sp.OwnerUserID,
		&sp.ChannelIDs,
		&sp.Platforms,
		&sp.ScheduledAt,
		&sp.Status,
		&sp.PostType,
		&sp.Caption,
		&sp.Media,
		&sp.ProviderState,
		&sp.CreatedAt,
		&sp.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &sp, nil
}

// ScheduledPostRangeFilters for calendar listing.
type ScheduledPostRangeFilters struct {
	OwnerUserID uuid.UUID
	From        time.Time
	To          time.Time
	ChannelID   *int64
}

// ListScheduledPostsInRange returns posts overlapping [from, to) for the owner.
func ListScheduledPostsInRange(ctx context.Context, f ScheduledPostRangeFilters) ([]entities.ScheduledPost, error) {
	db := database.GetDB()
	base := `
		SELECT id, owner_user_id, channel_ids, platforms, scheduled_at, status,
		       post_type, caption, media, provider_state, created_at, updated_at
		FROM scheduled_posts
		WHERE owner_user_id = $1
		  AND scheduled_at >= $2
		  AND scheduled_at < $3
		  AND status NOT IN ('cancelled')
	`
	var q string
	var args []interface{}
	if f.ChannelID != nil {
		q = base + ` AND $4 = ANY(channel_ids) ORDER BY scheduled_at ASC`
		args = []interface{}{f.OwnerUserID, f.From, f.To, *f.ChannelID}
	} else {
		q = base + ` ORDER BY scheduled_at ASC`
		args = []interface{}{f.OwnerUserID, f.From, f.To}
	}

	rows, err := db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanScheduledPosts(rows)
}

func scanScheduledPosts(rows *sql.Rows) ([]entities.ScheduledPost, error) {
	var out []entities.ScheduledPost
	for rows.Next() {
		var sp entities.ScheduledPost
		err := rows.Scan(
			&sp.ID,
			&sp.OwnerUserID,
			&sp.ChannelIDs,
			&sp.Platforms,
			&sp.ScheduledAt,
			&sp.Status,
			&sp.PostType,
			&sp.Caption,
			&sp.Media,
			&sp.ProviderState,
			&sp.CreatedAt,
			&sp.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		out = append(out, sp)
	}
	return out, rows.Err()
}

// UpdateScheduledPostStatusAndProviderState updates status and JSON state after Meta calls.
func UpdateScheduledPostStatusAndProviderState(ctx context.Context, id int64, ownerUserID uuid.UUID, status string, providerState []byte) error {
	db := database.GetDB()
	q := `
		UPDATE scheduled_posts
		SET status = $3, provider_state = $4, updated_at = NOW()
		WHERE id = $1 AND owner_user_id = $2
	`
	res, err := db.ExecContext(ctx, q, id, ownerUserID, status, providerState)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// CancelScheduledPost sets status to cancelled if still pending/submitting/scheduled/failed.
func CancelScheduledPost(ctx context.Context, id int64, ownerUserID uuid.UUID) error {
	db := database.GetDB()
	q := `
		UPDATE scheduled_posts
		SET status = 'cancelled', updated_at = NOW()
		WHERE id = $1 AND owner_user_id = $2
		  AND status IN ('pending', 'submitting', 'scheduled', 'failed')
	`
	res, err := db.ExecContext(ctx, q, id, ownerUserID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
