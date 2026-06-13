package publishing

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/google/uuid"
	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/entities/repositories"
	"github.com/shekhar8352/PostEaze/provider/googledrive"
	"github.com/shekhar8352/PostEaze/provider/youtube"
	"github.com/shekhar8352/PostEaze/services/google_drive_service"
	"github.com/shekhar8352/PostEaze/services/youtube_service"
)

type YouTubePublisher struct {
	YTClient    *youtube.Client
	YTService   *youtube_service.YouTubeService
	DriveClient *googledrive.Client
	DriveTokens *google_drive_service.TokenService
}

func NewYouTubePublisher() *YouTubePublisher {
	return &YouTubePublisher{
		YTClient:    youtube.NewClient(),
		YTService:   youtube_service.NewYouTubeService(),
		DriveClient: googledrive.NewClient(),
		DriveTokens: google_drive_service.NewTokenService(),
	}
}

func (p *YouTubePublisher) Platform() string { return PlatformYouTube }

type YouTubePublishInput struct {
	ChannelID int64
	OwnerID   uuid.UUID
	Title     string
	Caption   string
	Version   *entities.MediaVersion
}

func (p *YouTubePublisher) Publish(ctx context.Context, in YouTubePublishInput) (videoID string, err error) {
	if in.Version == nil {
		return "", fmt.Errorf("media version required")
	}
	if !strings.HasPrefix(in.Version.ContentType, "video/") {
		return "", fmt.Errorf("youtube requires video content")
	}

	accessToken, err := p.YTService.GetValidAccessToken(ctx, in.ChannelID)
	if err != nil {
		return "", err
	}

	meta := youtube.VideoMetadata{
		Title:       in.Title,
		Description: in.Caption,
		Privacy:     "private",
	}
	if meta.Title == "" {
		meta.Title = in.Version.FileName
	}

	var reader youtube.RangeReader
	switch in.Version.StorageProvider {
	case entities.StorageProviderGoogleDrive:
		driveToken, err := p.DriveTokens.GetValidAccessToken(ctx, in.OwnerID)
		if err != nil {
			return "", err
		}
		fileID := ""
		if in.Version.DriveFileID != nil {
			fileID = *in.Version.DriveFileID
		}
		revID := ""
		if in.Version.DriveRevisionID != nil {
			revID = *in.Version.DriveRevisionID
		}
		reader = youtube.NewDriveRangeReader(driveToken, fileID, revID, in.Version.FileSize, func(token, fid, rid, rangeHeader string) (io.ReadCloser, int64, error) {
			dl, err := p.DriveClient.Download(token, fid, rid, rangeHeader)
			if err != nil {
				return nil, 0, err
			}
			n := dl.ContentLength
			if n <= 0 {
				n = in.Version.FileSize
			}
			return dl.Body, n, nil
		})
	default:
		if in.Version.BlobURL == "" {
			return "", fmt.Errorf("missing blob url for upload")
		}
		reader = &youtube.HTTPRangeReader{
			URL:   in.Version.BlobURL,
			Total: in.Version.FileSize,
		}
	}

	return p.YTClient.ResumableUpload(accessToken, meta, reader)
}

func LoadVersionForYouTube(ctx context.Context, ownerID uuid.UUID, versionID int64) (*entities.MediaVersion, error) {
	// versionID only - need asset id; scan via join
	v, oid, err := repositories.GetMediaVersionWithOwner(ctx, versionID)
	if err != nil {
		return nil, err
	}
	if v == nil || oid != ownerID {
		return nil, fmt.Errorf("media version not found")
	}
	return v, nil
}
