package businessv1

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/entities/repositories"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/provider/google"
	"github.com/shekhar8352/PostEaze/provider/googledrive"
	"github.com/shekhar8352/PostEaze/services/google_drive_service"
	"github.com/shekhar8352/PostEaze/utils/blobstore"
	"github.com/shekhar8352/PostEaze/utils/encryption"
)

var (
	driveClient     = googledrive.NewClient()
	googleOAuth     = google.NewOAuthProvider()
	driveTokenSvc   = google_drive_service.NewTokenService()
)

func ConnectGoogleDrive(ctx context.Context, userIDStr string, req *modelsv1.ConnectGoogleDriveRequest) (*modelsv1.GoogleDriveStatusResponse, int, error) {
	ownerID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, 400, fmt.Errorf("invalid user")
	}

	redirectURI := req.RedirectURI
	if redirectURI == "" {
		redirectURI = os.Getenv("GOOGLE_OAUTH_REDIRECT_URI")
	}
	if redirectURI == "" {
		return nil, 500, fmt.Errorf("GOOGLE_OAUTH_REDIRECT_URI not configured")
	}

	tokens, err := googleOAuth.ExchangeCode(req.Code, redirectURI)
	if err != nil {
		return nil, 400, fmt.Errorf("oauth exchange: %w", err)
	}

	userInfo, _ := googleOAuth.GetUserInfo(tokens.AccessToken)
	encAccess, err := encryption.Encrypt(tokens.AccessToken)
	if err != nil {
		return nil, 500, err
	}
	var encRefresh []byte
	if tokens.RefreshToken != "" {
		encRefresh, err = encryption.Encrypt(tokens.RefreshToken)
		if err != nil {
			return nil, 500, err
		}
	}
	var expiresAt *time.Time
	if tokens.ExpiresIn > 0 {
		t := time.Now().Add(time.Duration(tokens.ExpiresIn) * time.Second)
		expiresAt = &t
	}

	var email *string
	if userInfo != nil && userInfo.Email != "" {
		email = &userInfo.Email
	}
	scopes := tokens.Scope
	ui := &entities.UserIntegration{
		UserID:               ownerID,
		Provider:             entities.UserIntegrationProviderGoogleDrive,
		ProviderAccountEmail: email,
		AccessToken:          encAccess,
		RefreshToken:         encRefresh,
		Scopes:               &scopes,
		ExpiresAt:            expiresAt,
	}
	if err := repositories.UpsertUserIntegration(ctx, ui); err != nil {
		return nil, 500, err
	}

	resp := &modelsv1.GoogleDriveStatusResponse{Connected: true, Scopes: scopes}
	if email != nil {
		resp.Email = *email
	}
	return resp, 200, nil
}

func GetGoogleDriveStatus(ctx context.Context, userIDStr string) (*modelsv1.GoogleDriveStatusResponse, int, error) {
	ownerID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, 400, fmt.Errorf("invalid user")
	}
	ui, err := repositories.GetUserIntegration(ctx, ownerID, entities.UserIntegrationProviderGoogleDrive)
	if err != nil {
		return nil, 500, err
	}
	if ui == nil {
		return &modelsv1.GoogleDriveStatusResponse{Connected: false}, 200, nil
	}
	resp := &modelsv1.GoogleDriveStatusResponse{Connected: true}
	if ui.ProviderAccountEmail != nil {
		resp.Email = *ui.ProviderAccountEmail
	}
	if ui.Scopes != nil {
		resp.Scopes = *ui.Scopes
	}
	return resp, 200, nil
}

func DisconnectGoogleDrive(ctx context.Context, userIDStr string) (int, error) {
	ownerID, err := uuid.Parse(userIDStr)
	if err != nil {
		return 400, fmt.Errorf("invalid user")
	}
	if err := repositories.RevokeUserIntegration(ctx, ownerID, entities.UserIntegrationProviderGoogleDrive); err != nil {
		return 404, fmt.Errorf("not connected")
	}
	return 200, nil
}

func ListGoogleDriveFiles(ctx context.Context, userIDStr, folderID, pageToken, query string) (*modelsv1.DriveFileListResponse, int, error) {
	ownerID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, 400, fmt.Errorf("invalid user")
	}
	token, err := driveTokenSvc.GetValidAccessToken(ctx, ownerID)
	if err != nil {
		return nil, 400, err
	}
	if folderID == "" {
		folderID = "root"
	}
	list, err := driveClient.ListFiles(token, folderID, pageToken, query)
	if err != nil {
		return nil, 500, err
	}
	items := make([]modelsv1.DriveFileItem, 0, len(list.Files))
	for _, f := range list.Files {
		isFolder := f.MimeType == "application/vnd.google-apps.folder"
		items = append(items, modelsv1.DriveFileItem{
			ID:            f.ID,
			Name:          f.Name,
			MimeType:      f.MimeType,
			Size:          googledrive.ParseSize(f.Size),
			IsFolder:      isFolder,
			ModifiedTime:  f.ModifiedTime,
			ThumbnailLink: f.ThumbnailLink,
		})
	}
	return &modelsv1.DriveFileListResponse{Files: items, NextPageToken: list.NextPageToken}, 200, nil
}

func ListGoogleDriveRevisions(ctx context.Context, userIDStr, fileID string) (*modelsv1.DriveRevisionListResponse, int, error) {
	ownerID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, 400, fmt.Errorf("invalid user")
	}
	token, err := driveTokenSvc.GetValidAccessToken(ctx, ownerID)
	if err != nil {
		return nil, 400, err
	}
	list, err := driveClient.ListRevisions(token, fileID)
	if err != nil {
		return nil, 500, err
	}
	items := make([]modelsv1.DriveRevisionItem, 0, len(list.Revisions))
	for _, r := range list.Revisions {
		items = append(items, modelsv1.DriveRevisionItem{
			ID:               r.ID,
			ModifiedTime:     r.ModifiedTime,
			KeepForever:      r.KeepForever,
			OriginalFilename: r.OriginalFilename,
			Size:             googledrive.ParseSize(r.Size),
			MimeType:         r.MimeType,
		})
	}
	return &modelsv1.DriveRevisionListResponse{Revisions: items}, 200, nil
}

func ImportFromGoogleDrive(ctx context.Context, userIDStr string, req *modelsv1.ImportGoogleDriveRequest) (*modelsv1.MediaAssetResponse, int, error) {
	ownerID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, 400, fmt.Errorf("invalid user")
	}
	token, err := driveTokenSvc.GetValidAccessToken(ctx, ownerID)
	if err != nil {
		return nil, 400, err
	}

	file, err := driveClient.GetFile(token, req.FileID)
	if err != nil {
		return nil, 400, err
	}
	if file.MimeType == "application/vnd.google-apps.folder" {
		return nil, 400, fmt.Errorf("cannot import a folder")
	}

	size := googledrive.ParseSize(file.Size)
	ct := normalizeDriveMime(file.MimeType)
	if !blobstore.IsAllowedContentType(ct) {
		return nil, 400, fmt.Errorf("unsupported file type %q", file.MimeType)
	}

	title := req.Title
	if title == "" {
		title = file.Name
	}
	assetType := "photo"
	if strings.HasPrefix(ct, "video/") {
		assetType = "video"
	}

	driveFileID := req.FileID
	asset := &entities.MediaAsset{
		OwnerUserID: ownerID,
		Title:       title,
		AssetType:   assetType,
		Status:      string(entities.MediaAssetStatusDraft),
		DriveFileID: &driveFileID,
	}
	if err := repositories.CreateMediaAsset(ctx, asset); err != nil {
		return nil, 500, err
	}

	version, err := buildVersionFromDrive(ctx, ownerID, token, asset.ID, 1, req.Label, file.Name, ct, size, driveFileID, req.RevisionID)
	if err != nil {
		return nil, 500, err
	}
	if err := repositories.CreateMediaVersion(ctx, version); err != nil {
		return nil, 500, err
	}
	if err := repositories.SetCurrentVersion(ctx, asset.ID, ownerID, version.ID); err != nil {
		return nil, 500, err
	}
	asset.CurrentVersionID = &version.ID

	return mapAssetToResponse(asset, []entities.MediaVersion{*version}), 200, nil
}

func ImportDriveRevision(ctx context.Context, userIDStr string, assetID int64, req *modelsv1.ImportDriveRevisionRequest) (*modelsv1.MediaVersionResponse, int, error) {
	ownerID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, 400, fmt.Errorf("invalid user")
	}
	token, err := driveTokenSvc.GetValidAccessToken(ctx, ownerID)
	if err != nil {
		return nil, 400, err
	}

	asset, err := repositories.GetMediaAssetByID(ctx, assetID, ownerID)
	if err != nil {
		return nil, 500, err
	}
	if asset == nil {
		return nil, 404, fmt.Errorf("asset not found")
	}
	if asset.DriveFileID == nil || *asset.DriveFileID == "" {
		return nil, 400, fmt.Errorf("asset is not linked to a Google Drive file")
	}

	file, err := driveClient.GetFile(token, *asset.DriveFileID)
	if err != nil {
		return nil, 400, err
	}
	size := googledrive.ParseSize(file.Size)
	ct := normalizeDriveMime(file.MimeType)

	nextNum, err := repositories.NextVersionNumber(ctx, assetID)
	if err != nil {
		return nil, 500, err
	}
	label := req.Label
	if label == "" {
		label = fmt.Sprintf("drive_rev_%s", req.RevisionID)
	}

	version, err := buildVersionFromDrive(ctx, ownerID, token, assetID, nextNum, label, file.Name, ct, size, *asset.DriveFileID, req.RevisionID)
	if err != nil {
		return nil, 500, err
	}
	if err := repositories.CreateMediaVersion(ctx, version); err != nil {
		return nil, 500, err
	}
	if err := repositories.SetCurrentVersion(ctx, assetID, ownerID, version.ID); err != nil {
		return nil, 500, err
	}
	return mapVersionToResponse(version), 200, nil
}

func buildVersionFromDrive(ctx context.Context, ownerID uuid.UUID, accessToken string, assetID int64, versionNum int, label, fileName, contentType string, fileSize int64, driveFileID, revisionID string) (*entities.MediaVersion, error) {
	if label == "" {
		label = "raw"
	}
	var revPtr *string
	if revisionID != "" {
		revPtr = &revisionID
	}
	drivePtr := &driveFileID

	if fileSize <= maxUploadSize {
		dl, err := driveClient.Download(accessToken, driveFileID, revisionID, "")
		if err != nil {
			return nil, err
		}
		defer dl.Body.Close()

		store := blobstore.Get()
		if store == nil {
			return nil, fmt.Errorf("blob storage not configured")
		}
		ext := extensionFromContentType(contentType)
		pathname := fmt.Sprintf("%s/%s%s", ownerID.String(), uuid.New().String(), ext)
		result, err := store.Upload(ctx, pathname, contentType, dl.Body)
		if err != nil {
			return nil, fmt.Errorf("blob upload: %w", err)
		}
		return &entities.MediaVersion{
			MediaAssetID:    assetID,
			VersionNumber:   versionNum,
			Label:           label,
			StorageProvider: entities.StorageProviderBlob,
			BlobURL:         result.URL,
			BlobPathKey:     result.Pathname,
			DriveFileID:     drivePtr,
			DriveRevisionID: revPtr,
			FileName:        fileName,
			ContentType:     contentType,
			FileSize:        fileSize,
			Metadata:        []byte(`{}`),
		}, nil
	}

	if !strings.HasPrefix(contentType, "video/") {
		return nil, fmt.Errorf("files over 50MB must be video; got %s (%d bytes)", contentType, fileSize)
	}

	return &entities.MediaVersion{
		MediaAssetID:    assetID,
		VersionNumber:   versionNum,
		Label:           label,
		StorageProvider: entities.StorageProviderGoogleDrive,
		BlobURL:         "",
		BlobPathKey:     "",
		DriveFileID:     drivePtr,
		DriveRevisionID: revPtr,
		FileName:        fileName,
		ContentType:     contentType,
		FileSize:        fileSize,
		Metadata:        []byte(`{}`),
	}, nil
}

func normalizeDriveMime(m string) string {
	switch m {
	case "image/jpg":
		return "image/jpeg"
	default:
		return m
	}
}

// StreamDriveVersion proxies a drive-backed version to the response writer.
func StreamDriveVersion(ctx context.Context, ownerID uuid.UUID, v *entities.MediaVersion, rangeHeader string) (body io.ReadCloser, contentType string, contentLength int64, contentRange string, statusCode int, err error) {
	if v.DriveFileID == nil {
		return nil, "", 0, "", 0, fmt.Errorf("missing drive file id")
	}
	token, err := driveTokenSvc.GetValidAccessToken(ctx, ownerID)
	if err != nil {
		return nil, "", 0, "", 0, err
	}
	revID := ""
	if v.DriveRevisionID != nil {
		revID = *v.DriveRevisionID
	}
	dl, err := driveClient.Download(token, *v.DriveFileID, revID, rangeHeader)
	if err != nil {
		return nil, "", 0, "", 0, err
	}
	return dl.Body, dl.ContentType, dl.ContentLength, dl.ContentRange, dl.StatusCode, nil
}
