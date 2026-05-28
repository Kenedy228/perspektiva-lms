package lecturematerial

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/shared/file"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		size     int64
		wantErr  bool
	}{
		{name: "valid mp4", fileName: "lesson.mp4", size: 100, wantErr: false},
		{name: "valid webm", fileName: "lesson.webm", size: 100, wantErr: false},
		{name: "valid pdf", fileName: "lecture.pdf", size: 100, wantErr: false},
		{name: "unsupported txt", fileName: "notes.txt", size: 100, wantErr: true},
		{name: "unsupported docx", fileName: "notes.docx", size: 100, wantErr: true},
		{name: "unsupported avi", fileName: "video.avi", size: 100, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := file.New(tt.fileName, tt.size)
			if err != nil {
				t.Fatal(err)
			}
			_, err = New(f)
			if (err != nil) != tt.wantErr {
				t.Fatalf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContentType(t *testing.T) {
	f, _ := file.New("lesson.mp4", 100)
	c, err := New(f)
	if err != nil {
		t.Fatal(err)
	}
	if ct := c.ContentType(); ct.String() != "lecture_material" {
		t.Fatalf("unexpected content type: %q", ct)
	}
}

func TestIsInteractive(t *testing.T) {
	f, _ := file.New("lesson.mp4", 100)
	c, err := New(f)
	if err != nil {
		t.Fatal(err)
	}
	if c.IsInteractive() {
		t.Fatal("expected lecture material content to be non-interactive")
	}
}

func TestClone(t *testing.T) {
	f, _ := file.New("lesson.mp4", 100)
	c, err := New(f)
	if err != nil {
		t.Fatal(err)
	}
	cloned := c.Clone()
	switch ct := cloned.(type) {
	case Content:
		if ct.File().Name() != "lesson.mp4" {
			t.Fatal("cloned content has different file")
		}
	default:
		t.Fatalf("unexpected cloned type: %T", cloned)
	}
}

func TestIsSupported(t *testing.T) {
	tests := []struct {
		name      string
		fileName  string
		supported bool
	}{
		{name: "mp4", fileName: "video.mp4", supported: true},
		{name: "webm", fileName: "video.webm", supported: true},
		{name: "pdf", fileName: "doc.pdf", supported: true},
		{name: "txt", fileName: "text.txt", supported: false},
		{name: "docx", fileName: "doc.docx", supported: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, _ := file.New(tt.fileName, 100)
			if got := IsSupported(f); got != tt.supported {
				t.Fatalf("IsSupported() = %v, want %v", got, tt.supported)
			}
		})
	}
}
