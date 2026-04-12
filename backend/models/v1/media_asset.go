package modelsv1

import "time"

// --- Requests ---

type CreateMediaAssetRequest struct {
	Title     string `json:"title" binding:"required"`
	AssetType string `json:"asset_type" binding:"required,oneof=photo video"`
	Label     string `json:"label"`
}

type UpdateMediaAssetRequest struct {
	Title  string `json:"title"`
	Status string `json:"status" binding:"omitempty,oneof=draft ready"`
}

type SetCurrentVersionRequest struct {
	VersionID int64 `json:"version_id" binding:"required"`
}

type AddVersionRequest struct {
	Label string `json:"label"`
	Notes string `json:"notes"`
}

type PublishMediaAssetRequest struct {
	ChannelIDs []int64 `json:"channel_ids" binding:"required,min=1"`
	Caption    string  `json:"caption"`
}

// --- Responses ---

type MediaAssetResponse struct {
	ID               int64                  `json:"id"`
	Title            string                 `json:"title"`
	AssetType        string                 `json:"asset_type"`
	Status           string                 `json:"status"`
	CurrentVersionID *int64                 `json:"current_version_id"`
	Versions         []MediaVersionResponse `json:"versions,omitempty"`
	CreatedAt        string                 `json:"created_at"`
	UpdatedAt        string                 `json:"updated_at"`
}

type MediaVersionResponse struct {
	ID            int64  `json:"id"`
	VersionNumber int    `json:"version_number"`
	Label         string `json:"label"`
	BlobURL       string `json:"blob_url"`
	FileName      string `json:"file_name"`
	ContentType   string `json:"content_type"`
	FileSize      int64  `json:"file_size"`
	Metadata      any    `json:"metadata"`
	Notes         string `json:"notes"`
	CreatedAt     string `json:"created_at"`
}

type MediaAssetListResponse struct {
	Assets []MediaAssetResponse `json:"assets"`
	Total  int                  `json:"total"`
	Limit  int                  `json:"limit"`
	Offset int                  `json:"offset"`
}

type UploadResponse struct {
	URL         string `json:"url"`
	Pathname    string `json:"pathname"`
	ContentType string `json:"content_type"`
	FileSize    int64  `json:"file_size"`
}

type ListMediaAssetsQuery struct {
	Status string `form:"status"`
	Limit  int    `form:"limit,default=20"`
	Offset int    `form:"offset,default=0"`
}

// Helper to format time consistently
func FormatTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}
