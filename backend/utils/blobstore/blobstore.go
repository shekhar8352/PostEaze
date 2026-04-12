package blobstore

import (
	"context"
	"io"
)

type BlobResult struct {
	URL         string `json:"url"`
	Pathname    string `json:"pathname"`
	ContentType string `json:"contentType"`
	Size        int64  `json:"size"`
}

type BlobStore interface {
	Upload(ctx context.Context, pathname string, contentType string, body io.Reader) (*BlobResult, error)
	Delete(ctx context.Context, urls []string) error
}

var store BlobStore

func Init(s BlobStore) {
	store = s
}

func Get() BlobStore {
	return store
}
