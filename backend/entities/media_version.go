package entities

import "time"

type MediaVersion struct {
	ID              int64     `db:"id"`
	MediaAssetID    int64     `db:"media_asset_id"`
	VersionNumber   int       `db:"version_number"`
	Label           string    `db:"label"`
	StorageProvider string    `db:"storage_provider"`
	BlobURL         string    `db:"blob_url"`
	BlobPathKey     string    `db:"blob_path_key"`
	DriveFileID     *string   `db:"drive_file_id"`
	DriveRevisionID *string   `db:"drive_revision_id"`
	FileName        string    `db:"file_name"`
	ContentType     string    `db:"content_type"`
	FileSize        int64     `db:"file_size"`
	Metadata        []byte    `db:"metadata"` // JSONB
	Notes           string    `db:"notes"`
	CreatedAt       time.Time `db:"created_at"`
}
