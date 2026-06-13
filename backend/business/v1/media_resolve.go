package businessv1

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/utils/mediasign"
)

const (
	igImageMaxBytes    = 8 << 20   // 8 MB
	igFeedVideoMax     = 100 << 20 // 100 MB
	igReelVideoMax     = 1 << 30   // 1 GB
)

// ResolveVersionMediaURL returns a URL platforms can fetch for this version.
func ResolveVersionMediaURL(version *entities.MediaVersion) (string, error) {
	if version == nil {
		return "", fmt.Errorf("version is nil")
	}
	switch version.StorageProvider {
	case "", entities.StorageProviderBlob:
		if version.BlobURL == "" {
			return "", fmt.Errorf("blob url missing")
		}
		return version.BlobURL, nil
	case entities.StorageProviderGoogleDrive:
		return mediasign.SignedStreamURL(version.ID, 6*time.Hour)
	default:
		return "", fmt.Errorf("unsupported storage provider %q", version.StorageProvider)
	}
}

func ValidateMediaForInstagram(version *entities.MediaVersion, postType string) error {
	if version == nil {
		return fmt.Errorf("version required")
	}
	isVideo := strings.HasPrefix(version.ContentType, "video/")
	isImage := strings.HasPrefix(version.ContentType, "image/")

	switch postType {
	case "image":
		if !isImage {
			return fmt.Errorf("image post requires an image file")
		}
		if version.FileSize > igImageMaxBytes {
			return fmt.Errorf("instagram images must be <= 8MB (got %d bytes)", version.FileSize)
		}
	case "video":
		if !isVideo {
			return fmt.Errorf("video post requires a video file")
		}
		if version.FileSize > igFeedVideoMax {
			return fmt.Errorf("instagram feed video must be <= 100MB (got %d bytes); use Reels for larger files", version.FileSize)
		}
	case "reel":
		if !isVideo {
			return fmt.Errorf("reel requires a video file")
		}
		if version.FileSize > igReelVideoMax {
			return fmt.Errorf("instagram reel video must be <= 1GB (got %d bytes)", version.FileSize)
		}
	default:
		return nil
	}
	return nil
}

func ValidateMediaForYouTube(version *entities.MediaVersion) error {
	if version == nil {
		return fmt.Errorf("version required")
	}
	if !strings.HasPrefix(version.ContentType, "video/") {
		return fmt.Errorf("youtube requires a video file")
	}
	return nil
}

// CopyDriveVersionToBlobIfSmall streams a drive-backed version to blob when under maxUploadSize.
func CopyDriveVersionToBlobIfSmall(ctx context.Context, ownerID uuid.UUID, version *entities.MediaVersion) (*entities.MediaVersion, error) {
	if version.StorageProvider != entities.StorageProviderGoogleDrive {
		return version, nil
	}
	if version.FileSize > maxUploadSize {
		return version, nil
	}
	// For small drive-backed files, copy to blob at publish time if needed
	token, err := driveTokenSvc.GetValidAccessToken(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	driveFileID := ""
	if version.DriveFileID != nil {
		driveFileID = *version.DriveFileID
	}
	revID := ""
	if version.DriveRevisionID != nil {
		revID = *version.DriveRevisionID
	}
	v, err := buildVersionFromDrive(ctx, ownerID, token, version.MediaAssetID, version.VersionNumber, version.Label+"_blob", version.FileName, version.ContentType, version.FileSize, driveFileID, revID)
	if err != nil {
		return nil, err
	}
	v.ID = version.ID
	v.Metadata = version.Metadata
	v.Notes = version.Notes
	return v, nil
}
