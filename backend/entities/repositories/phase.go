package repositories

import (
	"context"
	"database/sql"

	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/utils/database"
)

// CreatePhase inserts a phase; ID and timestamps are populated on the passed struct.
func CreatePhase(ctx context.Context, p *entities.Phase) error {
	db := database.GetDB()
	q := `
		INSERT INTO phases
			(studio_id, name, slug, order_index, kind, wip_limit, is_default, color)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`
	return db.QueryRowContext(ctx, q,
		p.StudioID, p.Name, p.Slug, p.OrderIndex, p.Kind,
		p.WIPLimit, p.IsDefault, p.Color,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

// GetPhaseByID fetches a phase by id.
func GetPhaseByID(ctx context.Context, id int64) (*entities.Phase, error) {
	db := database.GetDB()
	q := `
		SELECT id, studio_id, name, slug, order_index, kind,
		       wip_limit, is_default, color, created_at, updated_at
		FROM phases WHERE id = $1
	`
	var p entities.Phase
	err := db.QueryRowContext(ctx, q, id).Scan(
		&p.ID, &p.StudioID, &p.Name, &p.Slug, &p.OrderIndex, &p.Kind,
		&p.WIPLimit, &p.IsDefault, &p.Color, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// ListPhasesByStudio returns all phases for a studio ordered by order_index.
func ListPhasesByStudio(ctx context.Context, studioID int64) ([]entities.Phase, error) {
	db := database.GetDB()
	q := `
		SELECT id, studio_id, name, slug, order_index, kind,
		       wip_limit, is_default, color, created_at, updated_at
		FROM phases
		WHERE studio_id = $1
		ORDER BY order_index ASC, id ASC
	`
	rows, err := db.QueryContext(ctx, q, studioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []entities.Phase
	for rows.Next() {
		var p entities.Phase
		if err := rows.Scan(
			&p.ID, &p.StudioID, &p.Name, &p.Slug, &p.OrderIndex, &p.Kind,
			&p.WIPLimit, &p.IsDefault, &p.Color, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// GetPhaseByStudioKind returns the first phase in a studio matching a kind. Used
// for lifecycle auto-transitions (e.g., find the "scheduled" phase of a studio).
func GetPhaseByStudioKind(ctx context.Context, studioID int64, kind string) (*entities.Phase, error) {
	db := database.GetDB()
	q := `
		SELECT id, studio_id, name, slug, order_index, kind,
		       wip_limit, is_default, color, created_at, updated_at
		FROM phases
		WHERE studio_id = $1 AND kind = $2
		ORDER BY order_index ASC, id ASC
		LIMIT 1
	`
	var p entities.Phase
	err := db.QueryRowContext(ctx, q, studioID, kind).Scan(
		&p.ID, &p.StudioID, &p.Name, &p.Slug, &p.OrderIndex, &p.Kind,
		&p.WIPLimit, &p.IsDefault, &p.Color, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// UpdatePhase updates the mutable metadata on a phase (not order).
func UpdatePhase(ctx context.Context, id int64, name, slug, kind, color string, wipLimit *int) error {
	db := database.GetDB()
	q := `
		UPDATE phases
		SET name = $2, slug = $3, kind = $4, color = $5, wip_limit = $6, updated_at = NOW()
		WHERE id = $1
	`
	res, err := db.ExecContext(ctx, q, id, name, slug, kind, color, wipLimit)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// ReorderPhases sets order_index for a batch of phase IDs (0-based in the order
// provided). It runs in a transaction for atomicity.
func ReorderPhases(ctx context.Context, studioID int64, orderedPhaseIDs []int64) error {
	db := database.GetDB()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	for idx, pid := range orderedPhaseIDs {
		if _, err := tx.ExecContext(ctx,
			`UPDATE phases SET order_index = $1, updated_at = NOW() WHERE id = $2 AND studio_id = $3`,
			idx, pid, studioID,
		); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// DeletePhase removes a phase. Caller must ensure no pieces reference it
// (ON DELETE RESTRICT on pieces.phase_id).
func DeletePhase(ctx context.Context, id int64) error {
	db := database.GetDB()
	res, err := db.ExecContext(ctx, `DELETE FROM phases WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// CountPiecesInPhase returns the number of active pieces currently in a phase
// (useful for WIP-limit enforcement and for safe delete checks).
func CountPiecesInPhase(ctx context.Context, phaseID int64) (int, error) {
	db := database.GetDB()
	var n int
	err := db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM pieces WHERE phase_id = $1 AND status = 'active'`,
		phaseID,
	).Scan(&n)
	return n, err
}
