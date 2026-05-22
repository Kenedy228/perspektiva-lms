package minio

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	courseports "gitflic.ru/lms/backend/internal/application/ports/course"
)

type ObjectStorage struct {
	cfg    Config
	client *http.Client
}

func NewObjectStorage(cfg Config, client *http.Client) (*ObjectStorage, error) {
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("validate minio config: %w", err)
	}
	if client == nil {
		client = http.DefaultClient
	}
	return &ObjectStorage{cfg: cfg, client: client}, nil
}

func (s *ObjectStorage) PutObject(ctx context.Context, key string, body io.Reader, size int64, contentType string) (courseports.StoredObject, error) {
	if strings.TrimSpace(key) == "" {
		return courseports.StoredObject{}, errors.New("object key is required")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, s.objectURL(key), body)
	if err != nil {
		return courseports.StoredObject{}, err
	}
	req.ContentLength = size
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return courseports.StoredObject{}, fmt.Errorf("put minio object: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return courseports.StoredObject{}, fmt.Errorf("put minio object: unexpected status %d", resp.StatusCode)
	}
	return courseports.StoredObject{Key: key, ContentURL: s.publicObjectURL(key)}, nil
}

func (s *ObjectStorage) GetDownloadURL(_ context.Context, key string, _ time.Duration) (string, error) {
	if strings.TrimSpace(key) == "" {
		return "", errors.New("object key is required")
	}
	return s.publicObjectURL(key), nil
}

func (s *ObjectStorage) DeleteObject(ctx context.Context, key string) error {
	if strings.TrimSpace(key) == "" {
		return errors.New("object key is required")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, s.objectURL(key), nil)
	if err != nil {
		return err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("delete minio object: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("delete minio object: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func (s *ObjectStorage) objectURL(key string) string {
	return strings.TrimRight(s.cfg.Endpoint, "/") + "/" + url.PathEscape(s.cfg.Bucket) + "/" + url.PathEscape(key)
}

func (s *ObjectStorage) publicObjectURL(key string) string {
	base := s.cfg.PublicURL
	if base == "" {
		base = s.cfg.Endpoint
	}
	return strings.TrimRight(base, "/") + "/" + url.PathEscape(s.cfg.Bucket) + "/" + url.PathEscape(key)
}
