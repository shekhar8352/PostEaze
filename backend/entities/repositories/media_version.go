package repositories

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/utils/database"
)

func CreateMediaVersion(ctx context.Context, v *entities.MediaVersion) error {
	db := database.GetDB()
	q := `
		INSERT INTO media_versions (
			media_asset_id, version_number, label, blob_url, blob_path_key,
			file_name, content_type, file_size, metadata, notes
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at
	`
	return db.QueryRowContext(ctx, q,
		v.MediaAssetID, v.VersionNumber, v.Label, v.BlobURL, v.BlobPathKey,
		v.FileName, v.ContentType, v.FileSize, v.Metadata, v.Notes,
	).Scan(&v.ID, &v.CreatedAt)
}

// ListCurrentVersionsForAssets returns the current version row for each listed asset ID (same owner).
// Used by list APIs so clients can render thumbnails without loading full version history per asset.
func ListCurrentVersionsForAssets(ctx context.Context, ownerID uuid.UUID, assetIDs []int64) (map[int64]entities.MediaVersion, error) {
	out := make(map[int64]entities.MediaVersion)
	if len(assetIDs) == 0 {
		return out, nil
	}
	db := database.GetDB()
	q := `
		SELECT v.id, v.media_asset_id, v.version_number, v.label, v.blob_url, v.blob_path_key,
		       v.file_name, v.content_type, v.file_size, v.metadata, v.notes, v.created_at
		FROM media_versions v
		INNER JOIN media_assets a ON a.id = v.media_asset_id AND a.current_version_id = v.id
		WHERE a.owner_user_id = $1 AND a.id = ANY($2)
	`
	rows, err := db.QueryContext(ctx, q, ownerID, pq.Array(assetIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var v entities.MediaVersion
		if err := rows.Scan(
			&v.ID, &v.MediaAssetID, &v.VersionNumber, &v.Label, &v.BlobURL,
			&v.BlobPathKey, &v.FileName, &v.ContentType, &v.FileSize,
			&v.Metadata, &v.Notes, &v.CreatedAt,
		); err != nil {
			return nil, err
		}
		out[v.MediaAssetID] = v
	}
	return out, rows.Err()
}

func ListVersionsByAssetID(ctx context.Context, assetID int64) ([]entities.MediaVersion, error) {
	db := database.GetDB()
	q := `
		SELECT id, media_asset_id, version_number, label, blob_url, blob_path_key,
		       file_name, content_type, file_size, metadata, notes, created_at
		FROM media_versions
		WHERE media_asset_id = $1
		ORDER BY version_number ASC
	`
	rows, err := db.QueryContext(ctx, q, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []entities.MediaVersion
	for rows.Next() {
		var v entities.MediaVersion
		if err := rows.Scan(
			&v.ID, &v.MediaAssetID, &v.VersionNumber, &v.Label, &v.BlobURL,
			&v.BlobPathKey, &v.FileName, &v.ContentType, &v.FileSize,
			&v.Metadata, &v.Notes, &v.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

func GetMediaVersionByID(ctx context.Context, versionID int64, assetID int64) (*entities.MediaVersion, error) {
	db := database.GetDB()
	q := `
		SELECT id, media_asset_id, version_number, label, blob_url, blob_path_key,
		       file_name, content_type, file_size, metadata, notes, created_at
		FROM media_versions
		WHERE id = $1 AND media_asset_id = $2
	`
	var v entities.MediaVersion
	err := db.QueryRowContext(ctx, q, versionID, assetID).Scan(
		&v.ID, &v.MediaAssetID, &v.VersionNumber, &v.Label, &v.BlobURL,
		&v.BlobPathKey, &v.FileName, &v.ContentType, &v.FileSize,
		&v.Metadata, &v.Notes, &v.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func NextVersionNumber(ctx context.Context, assetID int64) (int, error) {
	db := database.GetDB()
	var n sql.NullInt64
	err := db.QueryRowContext(ctx,
		`SELECT MAX(version_number) FROM media_versions WHERE media_asset_id = $1`, assetID,
	).Scan(&n)
	if err != nil {
		return 1, err
	}
	if !n.Valid {
		return 1, nil
	}
	return int(n.Int64) + 1, nil
}

func DeleteMediaVersion(ctx context.Context, versionID int64, assetID int64) error {
	db := database.GetDB()
	q := `DELETE FROM media_versions WHERE id = $1 AND media_asset_id = $2`
	res, err := db.ExecContext(ctx, q, versionID, assetID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// GetAllBlobURLsForAsset returns blob URLs for every version of an asset (for bulk deletion).
func GetAllBlobURLsForAsset(ctx context.Context, assetID int64) ([]string, error) {
	db := database.GetDB()
	rows, err := db.QueryContext(ctx,
		`SELECT blob_url FROM media_versions WHERE media_asset_id = $1`, assetID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []string
	for rows.Next() {
		var u string
		if err := rows.Scan(&u); err != nil {
			return nil, err
		}
		urls = append(urls, u)
	}
	return urls, rows.Err()
}
