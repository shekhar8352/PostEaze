package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/utils/database"
)

func CreateMediaAsset(ctx context.Context, a *entities.MediaAsset) error {
	db := database.GetDB()
	q := `
		INSERT INTO media_assets (owner_user_id, team_id, title, asset_type, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	return db.QueryRowContext(ctx, q,
		a.OwnerUserID, a.TeamID, a.Title, a.AssetType, a.Status,
	).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
}

func GetMediaAssetByID(ctx context.Context, id int64, ownerID uuid.UUID) (*entities.MediaAsset, error) {
	db := database.GetDB()
	q := `
		SELECT id, owner_user_id, team_id, title, asset_type, status,
		       current_version_id, created_at, updated_at
		FROM media_assets
		WHERE id = $1 AND owner_user_id = $2
	`
	var a entities.MediaAsset
	err := db.QueryRowContext(ctx, q, id, ownerID).Scan(
		&a.ID, &a.OwnerUserID, &a.TeamID, &a.Title, &a.AssetType,
		&a.Status, &a.CurrentVersionID, &a.CreatedAt, &a.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

type ListMediaAssetsFilters struct {
	OwnerUserID uuid.UUID
	Status      string
	Limit       int
	Offset      int
}

func ListMediaAssets(ctx context.Context, f ListMediaAssetsFilters) ([]entities.MediaAsset, int, error) {
	db := database.GetDB()

	countBase := `SELECT COUNT(*) FROM media_assets WHERE owner_user_id = $1`
	queryBase := `
		SELECT id, owner_user_id, team_id, title, asset_type, status,
		       current_version_id, created_at, updated_at
		FROM media_assets
		WHERE owner_user_id = $1
	`

	args := []interface{}{f.OwnerUserID}
	filterClause := ""
	if f.Status != "" {
		filterClause = fmt.Sprintf(" AND status = $%d", len(args)+1)
		args = append(args, f.Status)
	}

	var total int
	if err := db.QueryRowContext(ctx, countBase+filterClause, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	orderClause := " ORDER BY updated_at DESC"
	limitClause := fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, f.Limit, f.Offset)

	rows, err := db.QueryContext(ctx, queryBase+filterClause+orderClause+limitClause, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []entities.MediaAsset
	for rows.Next() {
		var a entities.MediaAsset
		if err := rows.Scan(
			&a.ID, &a.OwnerUserID, &a.TeamID, &a.Title, &a.AssetType,
			&a.Status, &a.CurrentVersionID, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		out = append(out, a)
	}
	return out, total, rows.Err()
}

func UpdateMediaAsset(ctx context.Context, id int64, ownerID uuid.UUID, title string, status string) error {
	db := database.GetDB()
	q := `
		UPDATE media_assets
		SET title = $3, status = $4, updated_at = NOW()
		WHERE id = $1 AND owner_user_id = $2
	`
	res, err := db.ExecContext(ctx, q, id, ownerID, title, status)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func SetCurrentVersion(ctx context.Context, assetID int64, ownerID uuid.UUID, versionID int64) error {
	db := database.GetDB()
	q := `
		UPDATE media_assets
		SET current_version_id = $3, updated_at = NOW()
		WHERE id = $1 AND owner_user_id = $2
	`
	res, err := db.ExecContext(ctx, q, assetID, ownerID, versionID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func UpdateMediaAssetStatus(ctx context.Context, id int64, ownerID uuid.UUID, status string) error {
	db := database.GetDB()
	q := `
		UPDATE media_assets SET status = $3, updated_at = NOW()
		WHERE id = $1 AND owner_user_id = $2
	`
	res, err := db.ExecContext(ctx, q, id, ownerID, status)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func DeleteMediaAsset(ctx context.Context, id int64, ownerID uuid.UUID) error {
	db := database.GetDB()
	q := `DELETE FROM media_assets WHERE id = $1 AND owner_user_id = $2`
	res, err := db.ExecContext(ctx, q, id, ownerID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
