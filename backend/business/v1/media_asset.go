package businessv1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/google/uuid"
	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/entities/repositories"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/utils/blobstore"
)

const maxUploadSize = 50 << 20 // 50 MB

func UploadMedia(ctx context.Context, userID string, filename string, contentType string, fileSize int64, body io.Reader) (*modelsv1.UploadResponse, int, error) {
	if fileSize > maxUploadSize {
		return nil, 400, fmt.Errorf("file size exceeds 50MB limit")
	}
	ct := strings.Split(contentType, ";")[0]
	ct = strings.TrimSpace(ct)
	if !blobstore.IsAllowedContentType(ct) {
		return nil, 400, fmt.Errorf("unsupported content type %q; allowed: jpeg, png, webp, mp4, mov, webm", ct)
	}

	store := blobstore.Get()
	if store == nil {
		return nil, 500, fmt.Errorf("blob storage not configured")
	}

	ext := extensionFromContentType(ct)
	pathname := fmt.Sprintf("%s/%s%s", userID, uuid.New().String(), ext)

	result, err := store.Upload(ctx, pathname, ct, body)
	if err != nil {
		return nil, 500, fmt.Errorf("upload failed: %w", err)
	}

	return &modelsv1.UploadResponse{
		URL:         result.URL,
		Pathname:    result.Pathname,
		ContentType: ct,
		FileSize:    fileSize,
	}, 200, nil
}

func CreateMediaAsset(ctx context.Context, userIDStr string, req *modelsv1.CreateMediaAssetRequest, upload *modelsv1.UploadResponse) (*modelsv1.MediaAssetResponse, int, error) {
	ownerID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, 400, fmt.Errorf("invalid user")
	}

	asset := &entities.MediaAsset{
		OwnerUserID: ownerID,
		Title:       req.Title,
		AssetType:   req.AssetType,
		Status:      string(entities.MediaAssetStatusDraft),
	}
	if err := repositories.CreateMediaAsset(ctx, asset); err != nil {
		return nil, 500, fmt.Errorf("create asset: %w", err)
	}

	label := req.Label
	if label == "" {
		label = "raw"
	}

	version := &entities.MediaVersion{
		MediaAssetID:  asset.ID,
		VersionNumber: 1,
		Label:         label,
		BlobURL:       upload.URL,
		BlobPathKey:   upload.Pathname,
		FileName:      req.Title,
		ContentType:   upload.ContentType,
		FileSize:      upload.FileSize,
		Metadata:      []byte(`{}`),
		Notes:         "",
	}
	if err := repositories.CreateMediaVersion(ctx, version); err != nil {
		return nil, 500, fmt.Errorf("create version: %w", err)
	}

	if err := repositories.SetCurrentVersion(ctx, asset.ID, ownerID, version.ID); err != nil {
		return nil, 500, fmt.Errorf("set current version: %w", err)
	}
	asset.CurrentVersionID = &version.ID

	return mapAssetToResponse(asset, []entities.MediaVersion{*version}), 200, nil
}

func GetMediaAsset(ctx context.Context, userIDStr string, assetID int64) (*modelsv1.MediaAssetResponse, int, error) {
	ownerID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, 400, fmt.Errorf("invalid user")
	}
	asset, err := repositories.GetMediaAssetByID(ctx, assetID, ownerID)
	if err != nil {
		return nil, 500, err
	}
	if asset == nil {
		return nil, 404, fmt.Errorf("asset not found")
	}
	versions, err := repositories.ListVersionsByAssetID(ctx, assetID)
	if err != nil {
		return nil, 500, err
	}
	return mapAssetToResponse(asset, versions), 200, nil
}

func ListMediaAssets(ctx context.Context, userIDStr string, q modelsv1.ListMediaAssetsQuery) (*modelsv1.MediaAssetListResponse, int, error) {
	ownerID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, 400, fmt.Errorf("invalid user")
	}
	if q.Limit <= 0 || q.Limit > 100 {
		q.Limit = 20
	}
	if q.Offset < 0 {
		q.Offset = 0
	}

	assets, total, err := repositories.ListMediaAssets(ctx, repositories.ListMediaAssetsFilters{
		OwnerUserID: ownerID,
		Status:      q.Status,
		Limit:       q.Limit,
		Offset:      q.Offset,
	})
	if err != nil {
		return nil, 500, err
	}

	assetIDs := make([]int64, 0, len(assets))
	for i := range assets {
		assetIDs = append(assetIDs, assets[i].ID)
	}
	currentByAssetID, err := repositories.ListCurrentVersionsForAssets(ctx, ownerID, assetIDs)
	if err != nil {
		return nil, 500, err
	}

	items := make([]modelsv1.MediaAssetResponse, 0, len(assets))
	for i := range assets {
		var versions []entities.MediaVersion
		if v, ok := currentByAssetID[assets[i].ID]; ok {
			versions = []entities.MediaVersion{v}
		}
		items = append(items, *mapAssetToResponse(&assets[i], versions))
	}

	return &modelsv1.MediaAssetListResponse{
		Assets: items,
		Total:  total,
		Limit:  q.Limit,
		Offset: q.Offset,
	}, 200, nil
}

func UpdateMediaAsset(ctx context.Context, userIDStr string, assetID int64, req *modelsv1.UpdateMediaAssetRequest) (int, error) {
	ownerID, err := uuid.Parse(userIDStr)
	if err != nil {
		return 400, fmt.Errorf("invalid user")
	}

	asset, err := repositories.GetMediaAssetByID(ctx, assetID, ownerID)
	if err != nil {
		return 500, err
	}
	if asset == nil {
		return 404, fmt.Errorf("asset not found")
	}

	title := req.Title
	if title == "" {
		title = asset.Title
	}
	status := req.Status
	if status == "" {
		status = asset.Status
	}

	if err := repositories.UpdateMediaAsset(ctx, assetID, ownerID, title, status); err != nil {
		return 500, err
	}
	return 200, nil
}

func DeleteMediaAsset(ctx context.Context, userIDStr string, assetID int64) (int, error) {
	ownerID, err := uuid.Parse(userIDStr)
	if err != nil {
		return 400, fmt.Errorf("invalid user")
	}

	asset, err := repositories.GetMediaAssetByID(ctx, assetID, ownerID)
	if err != nil {
		return 500, err
	}
	if asset == nil {
		return 404, fmt.Errorf("asset not found")
	}

	urls, err := repositories.GetAllBlobURLsForAsset(ctx, assetID)
	if err != nil {
		return 500, err
	}

	if err := repositories.DeleteMediaAsset(ctx, assetID, ownerID); err != nil {
		return 500, err
	}

	if store := blobstore.Get(); store != nil && len(urls) > 0 {
		_ = store.Delete(ctx, urls)
	}

	return 200, nil
}

func AddVersion(ctx context.Context, userIDStr string, assetID int64, req *modelsv1.AddVersionRequest, upload *modelsv1.UploadResponse) (*modelsv1.MediaVersionResponse, int, error) {
	ownerID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, 400, fmt.Errorf("invalid user")
	}

	asset, err := repositories.GetMediaAssetByID(ctx, assetID, ownerID)
	if err != nil {
		return nil, 500, err
	}
	if asset == nil {
		return nil, 404, fmt.Errorf("asset not found")
	}

	nextNum, err := repositories.NextVersionNumber(ctx, assetID)
	if err != nil {
		return nil, 500, err
	}

	label := req.Label
	if label == "" {
		label = fmt.Sprintf("version_%d", nextNum)
	}

	version := &entities.MediaVersion{
		MediaAssetID:  assetID,
		VersionNumber: nextNum,
		Label:         label,
		BlobURL:       upload.URL,
		BlobPathKey:   upload.Pathname,
		FileName:      upload.Pathname,
		ContentType:   upload.ContentType,
		FileSize:      upload.FileSize,
		Metadata:      []byte(`{}`),
		Notes:         req.Notes,
	}
	if err := repositories.CreateMediaVersion(ctx, version); err != nil {
		return nil, 500, err
	}

	if err := repositories.SetCurrentVersion(ctx, assetID, ownerID, version.ID); err != nil {
		return nil, 500, err
	}

	return mapVersionToResponse(version), 200, nil
}

func DeleteVersion(ctx context.Context, userIDStr string, assetID int64, versionID int64) (int, error) {
	ownerID, err := uuid.Parse(userIDStr)
	if err != nil {
		return 400, fmt.Errorf("invalid user")
	}

	asset, err := repositories.GetMediaAssetByID(ctx, assetID, ownerID)
	if err != nil {
		return 500, err
	}
	if asset == nil {
		return 404, fmt.Errorf("asset not found")
	}

	version, err := repositories.GetMediaVersionByID(ctx, versionID, assetID)
	if err != nil {
		return 500, err
	}
	if version == nil {
		return 404, fmt.Errorf("version not found")
	}

	if err := repositories.DeleteMediaVersion(ctx, versionID, assetID); err != nil {
		return 500, err
	}

	if store := blobstore.Get(); store != nil {
		_ = store.Delete(ctx, []string{version.BlobURL})
	}

	if asset.CurrentVersionID != nil && *asset.CurrentVersionID == versionID {
		versions, _ := repositories.ListVersionsByAssetID(ctx, assetID)
		if len(versions) > 0 {
			last := versions[len(versions)-1]
			_ = repositories.SetCurrentVersion(ctx, assetID, ownerID, last.ID)
		}
	}

	return 200, nil
}

func SetCurrentVersion(ctx context.Context, userIDStr string, assetID int64, req *modelsv1.SetCurrentVersionRequest) (int, error) {
	ownerID, err := uuid.Parse(userIDStr)
	if err != nil {
		return 400, fmt.Errorf("invalid user")
	}

	asset, err := repositories.GetMediaAssetByID(ctx, assetID, ownerID)
	if err != nil {
		return 500, err
	}
	if asset == nil {
		return 404, fmt.Errorf("asset not found")
	}

	version, err := repositories.GetMediaVersionByID(ctx, req.VersionID, assetID)
	if err != nil {
		return 500, err
	}
	if version == nil {
		return 404, fmt.Errorf("version not found")
	}

	if err := repositories.SetCurrentVersion(ctx, assetID, ownerID, req.VersionID); err != nil {
		return 500, err
	}
	return 200, nil
}

func PublishMediaAsset(ctx context.Context, userIDStr string, assetID int64, req *modelsv1.PublishMediaAssetRequest) (*modelsv1.CreateScheduledPostResponse, int, error) {
	ownerID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, 400, fmt.Errorf("invalid user")
	}

	asset, err := repositories.GetMediaAssetByID(ctx, assetID, ownerID)
	if err != nil {
		return nil, 500, err
	}
	if asset == nil {
		return nil, 404, fmt.Errorf("asset not found")
	}
	if asset.CurrentVersionID == nil {
		return nil, 400, fmt.Errorf("no version selected for publishing")
	}

	version, err := repositories.GetMediaVersionByID(ctx, *asset.CurrentVersionID, assetID)
	if err != nil {
		return nil, 500, err
	}
	if version == nil {
		return nil, 400, fmt.Errorf("current version not found")
	}

	postType := "image"
	mediaKind := "image"
	if strings.HasPrefix(version.ContentType, "video/") {
		postType = "video"
		mediaKind = "video"
	}

	assetIDCopy := assetID
	schedReq := &modelsv1.CreateScheduledPostRequest{
		ChannelIDs: req.ChannelIDs,
		Platforms:  []string{"instagram"},
		PublishNow: true,
		PostType:   postType,
		Caption:    req.Caption,
		Media: modelsv1.ScheduledMediaPayload{
			Items: []modelsv1.ScheduledMediaItem{
				{Kind: mediaKind, URL: version.BlobURL, MediaAssetID: &assetIDCopy},
			},
		},
	}

	resp, code, err := CreateScheduledPost(ctx, userIDStr, schedReq)
	if err != nil {
		return resp, code, err
	}

	// Status is set by CreateScheduledPost on full publish success via mediapublish.MarkLinkedMediaAssetsPublished.
	return resp, code, nil
}

// --- helpers ---

func mapAssetToResponse(a *entities.MediaAsset, versions []entities.MediaVersion) *modelsv1.MediaAssetResponse {
	resp := &modelsv1.MediaAssetResponse{
		ID:               a.ID,
		Title:            a.Title,
		AssetType:        a.AssetType,
		Status:           a.Status,
		CurrentVersionID: a.CurrentVersionID,
		CreatedAt:        modelsv1.FormatTime(a.CreatedAt),
		UpdatedAt:        modelsv1.FormatTime(a.UpdatedAt),
	}
	if versions != nil {
		resp.Versions = make([]modelsv1.MediaVersionResponse, 0, len(versions))
		for i := range versions {
			resp.Versions = append(resp.Versions, *mapVersionToResponse(&versions[i]))
		}
	}
	return resp
}

func mapVersionToResponse(v *entities.MediaVersion) *modelsv1.MediaVersionResponse {
	var meta any
	_ = json.Unmarshal(v.Metadata, &meta)
	return &modelsv1.MediaVersionResponse{
		ID:            v.ID,
		VersionNumber: v.VersionNumber,
		Label:         v.Label,
		BlobURL:       v.BlobURL,
		FileName:      v.FileName,
		ContentType:   v.ContentType,
		FileSize:      v.FileSize,
		Metadata:      meta,
		Notes:         v.Notes,
		CreatedAt:     modelsv1.FormatTime(v.CreatedAt),
	}
}

func extensionFromContentType(ct string) string {
	switch ct {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "video/mp4":
		return ".mp4"
	case "video/quicktime":
		return ".mov"
	case "video/webm":
		return ".webm"
	default:
		return ""
	}
}

