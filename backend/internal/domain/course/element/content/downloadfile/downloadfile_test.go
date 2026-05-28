package downloadfile

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
		{name: "valid docx", fileName: "doc.docx", size: 100, wantErr: false},
		{name: "valid doc", fileName: "doc.doc", size: 100, wantErr: false},
		{name: "valid xlsx", fileName: "sheet.xlsx", size: 100, wantErr: false},
		{name: "valid xls", fileName: "sheet.xls", size: 100, wantErr: false},
		{name: "valid pdf", fileName: "file.pdf", size: 100, wantErr: false},
		{name: "valid txt", fileName: "text.txt", size: 100, wantErr: false},
		{name: "valid png", fileName: "img.png", size: 100, wantErr: false},
		{name: "valid avi", fileName: "video.avi", size: 100, wantErr: false},
		{name: "valid mov", fileName: "video.mov", size: 100, wantErr: false},
		{name: "unsupported mp3", fileName: "audio.mp3", size: 100, wantErr: true},
		{name: "unsupported exe", fileName: "app.exe", size: 100, wantErr: true},
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
	f, _ := file.New("doc.docx", 100)
	c, err := New(f)
	if err != nil {
		t.Fatal(err)
	}
	if ct := c.ContentType(); ct.String() != "download_file" {
		t.Fatalf("unexpected content type: %q", ct)
	}
}

func TestIsInteractive(t *testing.T) {
	f, _ := file.New("doc.docx", 100)
	c, err := New(f)
	if err != nil {
		t.Fatal(err)
	}
	if c.IsInteractive() {
		t.Fatal("expected download file content to be non-interactive")
	}
}

func TestClone(t *testing.T) {
	f, _ := file.New("doc.docx", 100)
	c, err := New(f)
	if err != nil {
		t.Fatal(err)
	}
	cloned := c.Clone()
	switch ct := cloned.(type) {
	case Content:
		if ct.File().Name() != "doc.docx" {
			t.Fatal("cloned content has different file")
		}
	default:
		t.Fatalf("unexpected cloned type: %T", cloned)
	}
}
