package minio

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	courseports "gitflic.ru/lms/backend/internal/application/ports/course"
	miniogo "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type ObjectStorage struct {
	cfg    Config
	client *miniogo.Client
}

func NewObjectStorage(cfg Config) (*ObjectStorage, error) {
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("validate minio config: %w", err)
	}
	// minio-go expects host:port without URL scheme; strip it if present.
	endpoint := cfg.Endpoint
	useSSL := cfg.UseSSL
	if strings.HasPrefix(endpoint, "https://") {
		endpoint = strings.TrimPrefix(endpoint, "https://")
		useSSL = true
	} else if strings.HasPrefix(endpoint, "http://") {
		endpoint = strings.TrimPrefix(endpoint, "http://")
	}
	client, err := miniogo.New(endpoint, &miniogo.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("create minio client: %w", err)
	}
	return &ObjectStorage{cfg: cfg, client: client}, nil
}

func (s *ObjectStorage) PutObject(ctx context.Context, key string, body io.Reader, size int64, contentType string) (courseports.StoredObject, error) {
	if strings.TrimSpace(key) == "" {
		return courseports.StoredObject{}, fmt.Errorf("object key is required")
	}
	_, err := s.client.PutObject(ctx, s.cfg.Bucket, key, body, size, miniogo.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return courseports.StoredObject{}, fmt.Errorf("put minio object: %w", err)
	}
	return courseports.StoredObject{Key: key, ContentURL: s.publicObjectURL(key)}, nil
}

func (s *ObjectStorage) GetDownloadURL(ctx context.Context, key string, expires time.Duration) (string, error) {
	if strings.TrimSpace(key) == "" {
		return "", fmt.Errorf("object key is required")
	}
	if expires <= 0 {
		expires = 15 * time.Minute
	}
	u, err := s.client.PresignedGetObject(ctx, s.cfg.Bucket, key, expires, url.Values{})
	if err != nil {
		return "", fmt.Errorf("presign minio object: %w", err)
	}
	return u.String(), nil
}

func (s *ObjectStorage) DeleteObject(ctx context.Context, key string) error {
	if strings.TrimSpace(key) == "" {
		return fmt.Errorf("object key is required")
	}
	err := s.client.RemoveObject(ctx, s.cfg.Bucket, key, miniogo.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("delete minio object: %w", err)
	}
	return nil
}

func (s *ObjectStorage) publicObjectURL(key string) string {
	base := s.cfg.PublicURL
	if base == "" {
		scheme := "http"
		if s.cfg.UseSSL {
			scheme = "https"
		}
		base = scheme + "://" + s.cfg.Endpoint
	}
	return strings.TrimRight(base, "/") + "/" + s.cfg.Bucket + "/" + key
}
