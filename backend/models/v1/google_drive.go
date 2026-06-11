package modelsv1

type ConnectGoogleDriveRequest struct {
	Code        string `json:"code" binding:"required"`
	RedirectURI string `json:"redirect_uri"`
}

type GoogleDriveStatusResponse struct {
	Connected bool   `json:"connected"`
	Email     string `json:"email,omitempty"`
	Scopes    string `json:"scopes,omitempty"`
}

type DriveFileItem struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	MimeType      string `json:"mime_type"`
	Size          int64  `json:"size"`
	IsFolder      bool   `json:"is_folder"`
	ModifiedTime  string `json:"modified_time,omitempty"`
	ThumbnailLink string `json:"thumbnail_link,omitempty"`
}

type DriveFileListResponse struct {
	Files         []DriveFileItem `json:"files"`
	NextPageToken string          `json:"next_page_token,omitempty"`
}

type DriveRevisionItem struct {
	ID               string `json:"id"`
	ModifiedTime     string `json:"modified_time"`
	KeepForever      bool   `json:"keep_forever"`
	OriginalFilename string `json:"original_filename,omitempty"`
	Size             int64  `json:"size"`
	MimeType         string `json:"mime_type,omitempty"`
}

type DriveRevisionListResponse struct {
	Revisions []DriveRevisionItem `json:"revisions"`
}

type ImportGoogleDriveRequest struct {
	FileID     string `json:"file_id" binding:"required"`
	RevisionID string `json:"revision_id"`
	Title      string `json:"title"`
	Label      string `json:"label"`
}

type ImportDriveRevisionRequest struct {
	RevisionID string `json:"revision_id" binding:"required"`
	Label      string `json:"label"`
}

type CreateYouTubeChannelRequest struct {
	Code        string `json:"code" binding:"required"`
	RedirectURI string `json:"redirect_uri"`
}

type CreateYouTubeChannelResponse struct {
	ChannelID   int64  `json:"channel_id"`
	ChannelName string `json:"channel_name"`
}
