package minio

import (
	"context"
	"strings"
	"testing"
	"time"
)

// PutObject, DeleteObject и GetDownloadURL (presign) требуют реального MinIO —
// покрываются интеграционным тестом postgres_integration_test.go.

func TestNewObjectStorage_InvalidConfig(t *testing.T) {
	cases := []struct {
		name string
		cfg  Config
	}{
		{name: "missing endpoint", cfg: Config{AccessKey: "key", SecretKey: "secret", Bucket: "b"}},
		{name: "missing bucket", cfg: Config{Endpoint: "localhost:9000", AccessKey: "key", SecretKey: "secret"}},
		{name: "missing access key", cfg: Config{Endpoint: "localhost:9000", SecretKey: "secret", Bucket: "b"}},
		{name: "missing secret key", cfg: Config{Endpoint: "localhost:9000", AccessKey: "key", Bucket: "b"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewObjectStorage(tc.cfg)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestGetDownloadURL_EmptyKey(t *testing.T) {
	s, err := NewObjectStorage(Config{
		Endpoint:  "localhost:9000",
		AccessKey: "minioadmin",
		SecretKey: "minioadmin",
		Bucket:    "lms",
	})
	if err != nil {
		t.Fatalf("create storage: %v", err)
	}
	_, err = s.GetDownloadURL(context.Background(), "   ", 15*time.Minute)
	if err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestPublicObjectURL_WithPublicURL(t *testing.T) {
	s, err := NewObjectStorage(Config{
		Endpoint:  "localhost:9000",
		AccessKey: "minioadmin",
		SecretKey: "minioadmin",
		Bucket:    "lms",
		PublicURL: "https://cdn.example.com",
	})
	if err != nil {
		t.Fatalf("create storage: %v", err)
	}
	got := s.publicObjectURL("lessons/video.mp4")
	want := "https://cdn.example.com/lms/lessons/video.mp4"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestPublicObjectURL_WithoutPublicURL_HTTP(t *testing.T) {
	s, err := NewObjectStorage(Config{
		Endpoint:  "localhost:9000",
		AccessKey: "minioadmin",
		SecretKey: "minioadmin",
		Bucket:    "lms",
		UseSSL:    false,
	})
	if err != nil {
		t.Fatalf("create storage: %v", err)
	}
	got := s.publicObjectURL("doc.pdf")
	if !strings.HasPrefix(got, "http://") {
		t.Fatalf("expected http scheme, got: %s", got)
	}
	if !strings.Contains(got, "lms/doc.pdf") {
		t.Fatalf("expected bucket/key in url, got: %s", got)
	}
}
