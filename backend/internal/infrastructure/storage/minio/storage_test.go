package minio

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestObjectStoragePutAndDownloadURL(t *testing.T) {
	var method string
	var path string
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		method = r.Method
		path = r.URL.Path
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("")),
			Header:     make(http.Header),
		}, nil
	})}

	storage, err := NewObjectStorage(Config{
		Endpoint:  "http://minio.local",
		Bucket:    "courses",
		PublicURL: "http://cdn.local",
	}, client)
	if err != nil {
		t.Fatalf("create storage: %v", err)
	}

	obj, err := storage.PutObject(context.Background(), "lesson.pdf", strings.NewReader("pdf"), 3, "application/pdf")
	if err != nil {
		t.Fatalf("put object: %v", err)
	}
	if method != http.MethodPut {
		t.Fatalf("expected PUT, got %s", method)
	}
	if path != "/courses/lesson.pdf" {
		t.Fatalf("expected object path, got %s", path)
	}
	if obj.ContentURL != "http://cdn.local/courses/lesson.pdf" {
		t.Fatalf("unexpected content url: %s", obj.ContentURL)
	}

	url, err := storage.GetDownloadURL(context.Background(), "lesson.pdf", 0)
	if err != nil {
		t.Fatalf("get url: %v", err)
	}
	if url != obj.ContentURL {
		t.Fatalf("expected %s, got %s", obj.ContentURL, url)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}
