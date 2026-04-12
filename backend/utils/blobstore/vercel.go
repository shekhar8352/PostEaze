package blobstore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const vercelBlobAPI = "https://blob.vercel-storage.com"

type VercelBlobStore struct {
	token      string
	httpClient *http.Client
}

func NewVercelBlobStore(token string) *VercelBlobStore {
	return &VercelBlobStore{
		token:      token,
		httpClient: &http.Client{},
	}
}

type vercelPutResponse struct {
	URL                string `json:"url"`
	Pathname           string `json:"pathname"`
	ContentType        string `json:"contentType"`
	ContentDisposition string `json:"contentDisposition"`
}

func (v *VercelBlobStore) Upload(ctx context.Context, pathname string, contentType string, body io.Reader) (*BlobResult, error) {
	url := fmt.Sprintf("%s/%s", vercelBlobAPI, pathname)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, body)
	if err != nil {
		return nil, fmt.Errorf("blobstore: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+v.token)
	req.Header.Set("x-api-version", "7")
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("x-content-type", contentType)

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("blobstore: upload request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("blobstore: upload failed (%d): %s", resp.StatusCode, string(respBody))
	}

	var result vercelPutResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("blobstore: decode response: %w", err)
	}

	return &BlobResult{
		URL:         result.URL,
		Pathname:    result.Pathname,
		ContentType: result.ContentType,
	}, nil
}

type vercelDeleteRequest struct {
	URLs []string `json:"urls"`
}

func (v *VercelBlobStore) Delete(ctx context.Context, urls []string) error {
	if len(urls) == 0 {
		return nil
	}

	payload, err := json.Marshal(vercelDeleteRequest{URLs: urls})
	if err != nil {
		return fmt.Errorf("blobstore: marshal delete: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, vercelBlobAPI+"/delete", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("blobstore: create delete request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+v.token)
	req.Header.Set("x-api-version", "7")
	req.Header.Set("Content-Type", "application/json")

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("blobstore: delete request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("blobstore: delete failed (%d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// AllowedContentTypes returns the set of MIME types accepted for upload.
func AllowedContentTypes() map[string]bool {
	return map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
		"video/mp4":  true,
		"video/quicktime": true,
		"video/webm": true,
	}
}

func IsAllowedContentType(ct string) bool {
	ct = strings.Split(ct, ";")[0]
	ct = strings.TrimSpace(ct)
	return AllowedContentTypes()[ct]
}
