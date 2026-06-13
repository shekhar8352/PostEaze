package googledrive

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const driveAPIBase = "https://www.googleapis.com/drive/v3"

type Client struct {
	HTTPClient *http.Client
}

func NewClient() *Client {
	return &Client{HTTPClient: &http.Client{Timeout: 0}} // no timeout for large downloads
}

type DriveFile struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	MimeType     string `json:"mimeType"`
	Size         string `json:"size"`
	ModifiedTime string `json:"modifiedTime"`
	IconLink     string `json:"iconLink"`
	ThumbnailLink string `json:"thumbnailLink"`
	Parents      []string `json:"parents"`
}

type FileList struct {
	Files         []DriveFile `json:"files"`
	NextPageToken string      `json:"nextPageToken"`
}

type Revision struct {
	ID               string `json:"id"`
	ModifiedTime     string `json:"modifiedTime"`
	KeepForever      bool   `json:"keepForever"`
	OriginalFilename string `json:"originalFilename"`
	Size             string `json:"size"`
	MimeType         string `json:"mimeType"`
}

type RevisionList struct {
	Revisions []Revision `json:"revisions"`
}

type DownloadResult struct {
	Body          io.ReadCloser
	ContentType   string
	ContentLength int64
	ContentRange  string
	StatusCode    int
}

func (c *Client) ListFiles(accessToken, folderID, pageToken, query string) (*FileList, error) {
	params := url.Values{}
	params.Set("pageSize", "50")
	params.Set("fields", "nextPageToken,files(id,name,mimeType,size,modifiedTime,iconLink,thumbnailLink,parents)")
	params.Set("orderBy", "folder,name")
	if pageToken != "" {
		params.Set("pageToken", pageToken)
	}

	q := "trashed = false"
	if folderID != "" && folderID != "root" {
		q += fmt.Sprintf(" and '%s' in parents", folderID)
	} else if folderID == "" || folderID == "root" {
		q += " and 'root' in parents"
	}
	if query != "" {
		q += " and (name contains '" + escapeQuery(query) + "' or fullText contains '" + escapeQuery(query) + "')"
	}
	params.Set("q", q)

	u := driveAPIBase + "/files?" + params.Encode()
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("drive list files (%d): %s", resp.StatusCode, string(body))
	}

	var out FileList
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetFile(accessToken, fileID string) (*DriveFile, error) {
	u := fmt.Sprintf("%s/files/%s?fields=id,name,mimeType,size,modifiedTime", driveAPIBase, fileID)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("drive get file (%d): %s", resp.StatusCode, string(body))
	}

	var out DriveFile
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ListRevisions(accessToken, fileID string) (*RevisionList, error) {
	u := fmt.Sprintf("%s/files/%s/revisions?fields=revisions(id,modifiedTime,keepForever,originalFilename,size,mimeType)", driveAPIBase, fileID)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("drive list revisions (%d): %s", resp.StatusCode, string(body))
	}

	var out RevisionList
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) Download(accessToken, fileID, revisionID, rangeHeader string) (*DownloadResult, error) {
	var u string
	if revisionID != "" {
		u = fmt.Sprintf("%s/files/%s/revisions/%s?alt=media", driveAPIBase, fileID, revisionID)
	} else {
		u = fmt.Sprintf("%s/files/%s?alt=media", driveAPIBase, fileID)
	}

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	if rangeHeader != "" {
		req.Header.Set("Range", rangeHeader)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("drive download (%d): %s", resp.StatusCode, string(body))
	}

	ct := resp.Header.Get("Content-Type")
	if ct == "" {
		ct = "application/octet-stream"
	}
	var cl int64
	if v := resp.Header.Get("Content-Length"); v != "" {
		fmt.Sscanf(v, "%d", &cl)
	}

	return &DownloadResult{
		Body:          resp.Body,
		ContentType:   ct,
		ContentLength: cl,
		ContentRange:  resp.Header.Get("Content-Range"),
		StatusCode:    resp.StatusCode,
	}, nil
}

func escapeQuery(s string) string {
	return strings.ReplaceAll(s, "'", "\\'")
}

func ParseSize(s string) int64 {
	var n int64
	fmt.Sscanf(s, "%d", &n)
	return n
}
