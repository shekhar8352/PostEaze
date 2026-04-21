package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/utils/database"
)

// CreatePiece inserts a new piece and populates ID/timestamps on the struct.
func CreatePiece(ctx context.Context, p *entities.Piece) error {
	db := database.GetDB()
	q := `
		INSERT INTO pieces
			(studio_id, phase_id, title, description, content_type, status,
			 assignee_user_id, due_at, position, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`
	return db.QueryRowContext(ctx, q,
		p.StudioID, p.PhaseID, p.Title, p.Description, p.ContentType, p.Status,
		p.AssigneeUserID, p.DueAt, p.Position, p.CreatedBy,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

// GetPieceByID fetches a piece by id.
func GetPieceByID(ctx context.Context, id int64) (*entities.Piece, error) {
	db := database.GetDB()
	q := `
		SELECT id, studio_id, phase_id, title, description, content_type, status,
		       assignee_user_id, due_at, position, created_by,
		       created_at, updated_at, archived_at
		FROM pieces
		WHERE id = $1
	`
	var p entities.Piece
	err := db.QueryRowContext(ctx, q, id).Scan(
		&p.ID, &p.StudioID, &p.PhaseID, &p.Title, &p.Description, &p.ContentType, &p.Status,
		&p.AssigneeUserID, &p.DueAt, &p.Position, &p.CreatedBy,
		&p.CreatedAt, &p.UpdatedAt, &p.ArchivedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// ListPiecesFilters narrows a list query.
type ListPiecesFilters struct {
	StudioID       int64
	PhaseID        *int64     // optional
	Status         string     // "", "active", "archived"
	AssigneeUserID *uuid.UUID // optional
}

// ListPiecesByStudio returns pieces for the board view, ordered by position.
func ListPiecesByStudio(ctx context.Context, f ListPiecesFilters) ([]entities.Piece, error) {
	db := database.GetDB()
	base := `
		SELECT id, studio_id, phase_id, title, description, content_type, status,
		       assignee_user_id, due_at, position, created_by,
		       created_at, updated_at, archived_at
		FROM pieces
		WHERE studio_id = $1
	`
	args := []interface{}{f.StudioID}
	extra := ""

	if f.PhaseID != nil {
		extra += fmt.Sprintf(" AND phase_id = $%d", len(args)+1)
		args = append(args, *f.PhaseID)
	}
	if f.Status != "" {
		extra += fmt.Sprintf(" AND status = $%d", len(args)+1)
		args = append(args, f.Status)
	} else {
		extra += " AND status = 'active'"
	}
	if f.AssigneeUserID != nil {
		extra += fmt.Sprintf(" AND assignee_user_id = $%d", len(args)+1)
		args = append(args, *f.AssigneeUserID)
	}

	order := " ORDER BY phase_id ASC, position ASC, id ASC"

	rows, err := db.QueryContext(ctx, base+extra+order, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []entities.Piece
	for rows.Next() {
		var p entities.Piece
		if err := rows.Scan(
			&p.ID, &p.StudioID, &p.PhaseID, &p.Title, &p.Description, &p.ContentType, &p.Status,
			&p.AssigneeUserID, &p.DueAt, &p.Position, &p.CreatedBy,
			&p.CreatedAt, &p.UpdatedAt, &p.ArchivedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// GetNeighborPositionsInPhase returns the positions flanking an intended insertion
// slot inside the given phase, ordered by position. Used by the fractional-index
// move/reorder logic.
//
// If beforePieceID is non-nil, we return (prev, that piece's position).
// If afterPieceID is non-nil, we return (that piece's position, next).
// If both are nil, we return (last_position_in_phase, "").
func GetPhasePositions(ctx context.Context, phaseID int64) ([]string, error) {
	db := database.GetDB()
	rows, err := db.QueryContext(ctx,
		`SELECT position FROM pieces WHERE phase_id = $1 AND status = 'active' ORDER BY position ASC`,
		phaseID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// UpdatePieceFields updates the basic content fields of a piece.
func UpdatePieceFields(
	ctx context.Context,
	id int64,
	title, description, contentType string,
	assigneeUserID *uuid.UUID,
	dueAt *sql.NullTime,
) error {
	db := database.GetDB()
	q := `
		UPDATE pieces
		SET title = $2, description = $3, content_type = $4,
		    assignee_user_id = $5, due_at = $6, updated_at = NOW()
		WHERE id = $1
	`
	var due interface{}
	if dueAt != nil && dueAt.Valid {
		due = dueAt.Time
	} else {
		due = nil
	}
	res, err := db.ExecContext(ctx, q, id, title, description, contentType, assigneeUserID, due)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// MovePiece updates the phase_id and position of a piece atomically.
func MovePiece(ctx context.Context, id, phaseID int64, position string) error {
	db := database.GetDB()
	q := `
		UPDATE pieces
		SET phase_id = $2, position = $3, updated_at = NOW()
		WHERE id = $1
	`
	res, err := db.ExecContext(ctx, q, id, phaseID, position)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// SetPieceStatus sets the lifecycle status of a piece ("active" | "archived").
// When archiving, archived_at is set to NOW(); when restoring it is cleared.
func SetPieceStatus(ctx context.Context, id int64, status string) error {
	db := database.GetDB()
	var q string
	if status == string(entities.PieceStatusArchived) {
		q = `UPDATE pieces SET status = $2, archived_at = NOW(), updated_at = NOW() WHERE id = $1`
	} else {
		q = `UPDATE pieces SET status = $2, archived_at = NULL, updated_at = NOW() WHERE id = $1`
	}
	res, err := db.ExecContext(ctx, q, id, status)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// DeletePiece removes a piece (cascades to joins/activities/comments).
func DeletePiece(ctx context.Context, id int64) error {
	db := database.GetDB()
	res, err := db.ExecContext(ctx, `DELETE FROM pieces WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// ---- piece_assets ----

// LinkPieceAsset creates a (piece, asset, role) link. Safe to call when the row
// may already exist; duplicates return nil.
func LinkPieceAsset(ctx context.Context, pieceID, mediaAssetID int64, role string) error {
	db := database.GetDB()
	q := `
		INSERT INTO piece_assets (piece_id, media_asset_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING
	`
	_, err := db.ExecContext(ctx, q, pieceID, mediaAssetID, role)
	return err
}

// UnlinkPieceAsset removes a specific (piece, asset, role) link.
func UnlinkPieceAsset(ctx context.Context, pieceID, mediaAssetID int64, role string) error {
	db := database.GetDB()
	_, err := db.ExecContext(ctx,
		`DELETE FROM piece_assets WHERE piece_id = $1 AND media_asset_id = $2 AND role = $3`,
		pieceID, mediaAssetID, role,
	)
	return err
}

// ListPieceAssets returns all asset links for a piece.
func ListPieceAssets(ctx context.Context, pieceID int64) ([]entities.PieceAsset, error) {
	db := database.GetDB()
	rows, err := db.QueryContext(ctx,
		`SELECT piece_id, media_asset_id, role, added_at
		 FROM piece_assets WHERE piece_id = $1 ORDER BY added_at ASC`,
		pieceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []entities.PieceAsset
	for rows.Next() {
		var pa entities.PieceAsset
		if err := rows.Scan(&pa.PieceID, &pa.MediaAssetID, &pa.Role, &pa.AddedAt); err != nil {
			return nil, err
		}
		out = append(out, pa)
	}
	return out, rows.Err()
}

// ---- piece_scheduled_posts ----

// LinkPieceScheduledPost ties a scheduled_post to a piece (deduped on PK).
func LinkPieceScheduledPost(ctx context.Context, pieceID, scheduledPostID int64, role string) error {
	db := database.GetDB()
	q := `
		INSERT INTO piece_scheduled_posts (piece_id, scheduled_post_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING
	`
	_, err := db.ExecContext(ctx, q, pieceID, scheduledPostID, role)
	return err
}

// UnlinkPieceScheduledPost removes the (piece, scheduled_post) link.
func UnlinkPieceScheduledPost(ctx context.Context, pieceID, scheduledPostID int64) error {
	db := database.GetDB()
	_, err := db.ExecContext(ctx,
		`DELETE FROM piece_scheduled_posts WHERE piece_id = $1 AND scheduled_post_id = $2`,
		pieceID, scheduledPostID,
	)
	return err
}

// ListPieceScheduledPosts returns all scheduled_post links for a piece.
func ListPieceScheduledPosts(ctx context.Context, pieceID int64) ([]entities.PieceScheduledPost, error) {
	db := database.GetDB()
	rows, err := db.QueryContext(ctx,
		`SELECT piece_id, scheduled_post_id, role, created_at
		 FROM piece_scheduled_posts WHERE piece_id = $1 ORDER BY created_at ASC`,
		pieceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []entities.PieceScheduledPost
	for rows.Next() {
		var ps entities.PieceScheduledPost
		if err := rows.Scan(&ps.PieceID, &ps.ScheduledPostID, &ps.Role, &ps.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, ps)
	}
	return out, rows.Err()
}

// GetPieceIDByScheduledPost returns the piece linked to a scheduled_post as
// "primary" (or the first link found). Used by the publishing lifecycle hook
// to resolve which Piece to transition when a scheduled_post changes state.
func GetPieceIDByScheduledPost(ctx context.Context, scheduledPostID int64) (int64, error) {
	db := database.GetDB()
	var pieceID int64
	err := db.QueryRowContext(ctx,
		`SELECT piece_id FROM piece_scheduled_posts
		 WHERE scheduled_post_id = $1
		 ORDER BY (role = 'primary') DESC, created_at ASC
		 LIMIT 1`,
		scheduledPostID,
	).Scan(&pieceID)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return pieceID, err
}

// ---- piece_activities ----

// CreatePieceActivity inserts an audit log row.
func CreatePieceActivity(ctx context.Context, a *entities.PieceActivity) error {
	db := database.GetDB()
	q := `
		INSERT INTO piece_activities (piece_id, actor_user_id, kind, payload)
		VALUES ($1, $2, $3, COALESCE($4, '{}'::jsonb))
		RETURNING id, created_at
	`
	return db.QueryRowContext(ctx, q, a.PieceID, a.ActorUserID, a.Kind, a.Payload).
		Scan(&a.ID, &a.CreatedAt)
}

// ListPieceActivities returns the most recent activities for a piece (newest first).
func ListPieceActivities(ctx context.Context, pieceID int64, limit int) ([]entities.PieceActivity, error) {
	if limit <= 0 {
		limit = 50
	}
	db := database.GetDB()
	rows, err := db.QueryContext(ctx,
		`SELECT id, piece_id, actor_user_id, kind, payload, created_at
		 FROM piece_activities WHERE piece_id = $1
		 ORDER BY created_at DESC LIMIT $2`,
		pieceID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []entities.PieceActivity
	for rows.Next() {
		var a entities.PieceActivity
		if err := rows.Scan(&a.ID, &a.PieceID, &a.ActorUserID, &a.Kind, &a.Payload, &a.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// ---- piece_comments ----

// CreatePieceComment inserts a comment and populates ID/timestamps.
func CreatePieceComment(ctx context.Context, c *entities.PieceComment) error {
	db := database.GetDB()
	q := `
		INSERT INTO piece_comments (piece_id, user_id, body, parent_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	return db.QueryRowContext(ctx, q, c.PieceID, c.UserID, c.Body, c.ParentID).
		Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

// ListPieceComments returns the comments for a piece ordered chronologically.
func ListPieceComments(ctx context.Context, pieceID int64) ([]entities.PieceComment, error) {
	db := database.GetDB()
	rows, err := db.QueryContext(ctx,
		`SELECT id, piece_id, user_id, body, parent_id, created_at, updated_at
		 FROM piece_comments WHERE piece_id = $1
		 ORDER BY created_at ASC, id ASC`,
		pieceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []entities.PieceComment
	for rows.Next() {
		var c entities.PieceComment
		if err := rows.Scan(&c.ID, &c.PieceID, &c.UserID, &c.Body, &c.ParentID, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// DeletePieceComment deletes a comment by id (cascades to replies via FK).
func DeletePieceComment(ctx context.Context, id int64) error {
	db := database.GetDB()
	res, err := db.ExecContext(ctx, `DELETE FROM piece_comments WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
