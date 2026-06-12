package youtube

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	apiBase       = "https://www.googleapis.com/youtube/v3"
	uploadBase    = "https://www.googleapis.com/upload/youtube/v3"
	defaultChunk  = 16 * 1024 * 1024 // 16MB
)

type Client struct {
	HTTPClient *http.Client
}

func NewClient() *Client {
	return &Client{HTTPClient: &http.Client{Timeout: 0}}
}

type Channel struct {
	ID      string         `json:"id"`
	Snippet ChannelSnippet `json:"snippet"`
}

type ChannelSnippet struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	CustomURL   string `json:"customUrl"`
	Thumbnails  map[string]struct {
		URL string `json:"url"`
	} `json:"thumbnails"`
}

type ChannelListResponse struct {
	Items []Channel `json:"items"`
}

type VideoMetadata struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Privacy     string `json:"privacy"` // public | private | unlisted
}

type RangeReader interface {
	ReadRange(offset, length int64) (io.ReadCloser, int64, error)
	TotalSize() int64
}

func (c *Client) GetMyChannel(accessToken string) (*Channel, error) {
	u := apiBase + "/channels?part=snippet&mine=true"
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
		return nil, fmt.Errorf("youtube channels (%d): %s", resp.StatusCode, string(body))
	}

	var out ChannelListResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	if len(out.Items) == 0 {
		return nil, fmt.Errorf("no YouTube channel found for this account")
	}
	return &out.Items[0], nil
}

func (c *Client) ResumableUpload(accessToken string, meta VideoMetadata, reader RangeReader) (string, error) {
	if meta.Privacy == "" {
		meta.Privacy = "private"
	}
	videoBody := map[string]any{
		"snippet": map[string]string{
			"title":       meta.Title,
			"description": meta.Description,
		},
		"status": map[string]string{
			"privacyStatus": meta.Privacy,
		},
	}
	bodyBytes, _ := json.Marshal(videoBody)

	initURL := uploadBase + "/videos?uploadType=resumable&part=snippet,status"
	req, err := http.NewRequest(http.MethodPost, initURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("X-Upload-Content-Type", "video/*")
	req.Header.Set("X-Upload-Content-Length", strconv.FormatInt(reader.TotalSize(), 10))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("youtube resumable init (%d)", resp.StatusCode)
	}

	sessionURL := resp.Header.Get("Location")
	if sessionURL == "" {
		return "", fmt.Errorf("youtube resumable init: missing Location header")
	}

	total := reader.TotalSize()
	var offset int64
	for offset < total {
		chunkSize := int64(defaultChunk)
		if offset+chunkSize > total {
			chunkSize = total - offset
		}

		body, n, err := reader.ReadRange(offset, chunkSize)
		if err != nil {
			return "", fmt.Errorf("read chunk at %d: %w", offset, err)
		}
		if n <= 0 {
			break
		}

		end := offset + n - 1
		contentRange := fmt.Sprintf("bytes %d-%d/%d", offset, end, total)

		putReq, err := http.NewRequest(http.MethodPut, sessionURL, body)
		if err != nil {
			body.Close()
			return "", err
		}
		putReq.Header.Set("Authorization", "Bearer "+accessToken)
		putReq.Header.Set("Content-Length", strconv.FormatInt(n, 10))
		putReq.Header.Set("Content-Type", "video/*")
		putReq.ContentLength = n
		if end < total-1 {
			putReq.Header.Set("Content-Range", contentRange)
		}

		putResp, err := c.HTTPClient.Do(putReq)
		body.Close()
		if err != nil {
			return "", err
		}

		putBody, _ := io.ReadAll(putResp.Body)
		putResp.Body.Close()

		switch putResp.StatusCode {
		case http.StatusOK, http.StatusCreated:
			var result struct {
				ID string `json:"id"`
			}
			if err := json.Unmarshal(putBody, &result); err != nil {
				return "", err
			}
			return result.ID, nil
		case http.StatusPermanentRedirect: // 308 resume
			if cr := putResp.Header.Get("Range"); cr != "" {
				if parts := strings.Split(cr, "-"); len(parts) == 2 {
					if next, err := strconv.ParseInt(parts[1], 10, 64); err == nil {
						offset = next + 1
						continue
					}
				}
			}
			offset = end + 1
		default:
			return "", fmt.Errorf("youtube upload chunk (%d): %s", putResp.StatusCode, string(putBody))
		}
	}
	return "", fmt.Errorf("youtube upload incomplete")
}

// HTTPRangeReader streams from a remote URL with Range support.
type HTTPRangeReader struct {
	URL       string
	Total     int64
	Client    *http.Client
}

func (r *HTTPRangeReader) TotalSize() int64 { return r.Total }

func (r *HTTPRangeReader) ReadRange(offset, length int64) (io.ReadCloser, int64, error) {
	client := r.Client
	if client == nil {
		client = &http.Client{Timeout: 0}
	}
	end := offset + length - 1
	req, err := http.NewRequest(http.MethodGet, r.URL, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", offset, end))

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, 0, fmt.Errorf("http range read (%d): %s", resp.StatusCode, string(body))
	}
	n := resp.ContentLength
	if n < 0 {
		n = length
	}
	return resp.Body, n, nil
}

type driveRangeReader struct {
	accessToken string
	fileID      string
	revisionID  string
	total       int64
	download    func(accessToken, fileID, revisionID, rangeHeader string) (io.ReadCloser, int64, error)
}

func NewDriveRangeReader(accessToken, fileID, revisionID string, total int64, downloadFn func(string, string, string, string) (io.ReadCloser, int64, error)) *driveRangeReader {
	return &driveRangeReader{
		accessToken: accessToken,
		fileID:      fileID,
		revisionID:  revisionID,
		total:       total,
		download:    downloadFn,
	}
}

func (r *driveRangeReader) TotalSize() int64 { return r.total }

func (r *driveRangeReader) ReadRange(offset, length int64) (io.ReadCloser, int64, error) {
	end := offset + length - 1
	rangeHeader := fmt.Sprintf("bytes=%d-%d", offset, end)
	return r.download(r.accessToken, r.fileID, r.revisionID, rangeHeader)
}

// Ensure driveRangeReader implements RangeReader
var _ RangeReader = (*driveRangeReader)(nil)
var _ RangeReader = (*HTTPRangeReader)(nil)

// Retry helper with backoff for transient errors
func retryWithBackoff(attempts int, fn func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		if err = fn(); err == nil {
			return nil
		}
		if i < attempts-1 {
			time.Sleep(time.Duration(i+1) * 2 * time.Second)
		}
	}
	return err
}
