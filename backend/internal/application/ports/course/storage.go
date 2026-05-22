package course

import (
	"context"
	"io"
	"time"
)

type StoredObject struct {
	Key        string
	ContentURL string
}

type ObjectStorage interface {
	PutObject(ctx context.Context, key string, body io.Reader, size int64, contentType string) (StoredObject, error)
	GetDownloadURL(ctx context.Context, key string, expires time.Duration) (string, error)
	DeleteObject(ctx context.Context, key string) error
}
