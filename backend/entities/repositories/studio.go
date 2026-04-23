package repositories

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/utils/database"
)

// CreateStudio inserts a new studio for a team. The generated ID and timestamps
// are populated on the passed struct.
func CreateStudio(ctx context.Context, s *entities.Studio) error {
	db := database.GetDB()
	q := `
		INSERT INTO studios (team_id, name, piece_label)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	return db.QueryRowContext(ctx, q, s.TeamID, s.Name, s.PieceLabel).
		Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

// GetStudioByID fetches a studio by its primary key.
func GetStudioByID(ctx context.Context, id int64) (*entities.Studio, error) {
	db := database.GetDB()
	q := `
		SELECT id, team_id, name, piece_label, created_at, updated_at
		FROM studios WHERE id = $1
	`
	var s entities.Studio
	err := db.QueryRowContext(ctx, q, id).Scan(
		&s.ID, &s.TeamID, &s.Name, &s.PieceLabel, &s.CreatedAt, &s.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// GetStudioByTeamID returns the (MVP: single) studio for a team, or nil.
func GetStudioByTeamID(ctx context.Context, teamID uuid.UUID) (*entities.Studio, error) {
	db := database.GetDB()
	q := `
		SELECT id, team_id, name, piece_label, created_at, updated_at
		FROM studios WHERE team_id = $1
		ORDER BY id ASC
		LIMIT 1
	`
	var s entities.Studio
	err := db.QueryRowContext(ctx, q, teamID).Scan(
		&s.ID, &s.TeamID, &s.Name, &s.PieceLabel, &s.CreatedAt, &s.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// ListStudiosByTeam returns all studios owned by a team.
func ListStudiosByTeam(ctx context.Context, teamID uuid.UUID) ([]entities.Studio, error) {
	db := database.GetDB()
	q := `
		SELECT id, team_id, name, piece_label, created_at, updated_at
		FROM studios WHERE team_id = $1
		ORDER BY id ASC
	`
	rows, err := db.QueryContext(ctx, q, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []entities.Studio
	for rows.Next() {
		var s entities.Studio
		if err := rows.Scan(
			&s.ID, &s.TeamID, &s.Name, &s.PieceLabel, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// UpdateStudio updates a studio's mutable fields.
func UpdateStudio(ctx context.Context, id int64, name, pieceLabel string) error {
	db := database.GetDB()
	q := `
		UPDATE studios
		SET name = $2, piece_label = $3, updated_at = NOW()
		WHERE id = $1
	`
	res, err := db.ExecContext(ctx, q, id, name, pieceLabel)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// DeleteStudio removes a studio (cascades to phases/pieces/joins).
func DeleteStudio(ctx context.Context, id int64) error {
	db := database.GetDB()
	res, err := db.ExecContext(ctx, `DELETE FROM studios WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
