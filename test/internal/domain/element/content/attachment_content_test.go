//go:build legacy
// +build legacy

package content_test

import (
	"strings"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/shared/s3validator"
	"gitflic.ru/lms/backend/internal/domain/element"
	"gitflic.ru/lms/backend/internal/domain/element/content"
	"github.com/stretchr/testify/assert"
)

func TestNewAttachmentContent(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		sizeBytes int64
		wantErr   error
	}{
		{
			name:      "валидный pdf",
			key:       "attachments/course-1/handout.pdf",
			sizeBytes: 1024,
			wantErr:   nil,
		},
		{
			name:      "валидный docx",
			key:       "attachments/course-1/notes.docx",
			sizeBytes: 2048,
			wantErr:   nil,
		},
		{
			name:      "валидный xlsx",
			key:       "attachments/course-1/data.xlsx",
			sizeBytes: 4096,
			wantErr:   nil,
		},
		{
			name:      "валидный pptx",
			key:       "attachments/course-1/slides.pptx",
			sizeBytes: 8192,
			wantErr:   nil,
		},
		{
			name:      "валидный mp4",
			key:       "attachments/course-1/video.mp4",
			sizeBytes: 10 * 1024 * 1024,
			wantErr:   nil,
		},
		{
			name:      "валидный webm",
			key:       "attachments/course-1/video.webm",
			sizeBytes: 10 * 1024 * 1024,
			wantErr:   nil,
		},
		{
			name:      "размер ровно 700 мб",
			key:       "attachments/course-1/file.pdf",
			sizeBytes: 700 * 1024 * 1024,
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := content.NewAttachmentContent(tt.key, tt.sizeBytes)

			assert.NoError(t, err)
			assert.Equal(t, tt.key, got.Key())
			assert.Equal(t, tt.sizeBytes, got.SizeBytes())
			assert.False(t, got.IsInteractive())
			assert.Equal(t, element.TypeAttachment, got.Type())
		})
	}
}

func TestNewAttachmentContent_Fail(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		sizeBytes int64
		wantErr   error
	}{
		{
			name:      "пустой ключ",
			key:       "",
			sizeBytes: 1024,
			wantErr:   s3validator.ErrEmptyS3Key,
		},
		{
			name:      "ключ из пробелов",
			key:       "   ",
			sizeBytes: 1024,
			wantErr:   s3validator.ErrEmptyS3Key,
		},
		{
			name:      "слишком длинный ключ",
			key:       strings.Repeat("a", 1025) + ".pdf",
			sizeBytes: 1024,
			wantErr:   s3validator.ErrTooLongS3Key,
		},
		{
			name:      "невалидный s3 ключ",
			key:       "/attachments/course-1/file.pdf",
			sizeBytes: 1024,
			wantErr:   s3validator.ErrInvalidS3Key,
		},
		{
			name:      "небезопасный path segment",
			key:       "attachments/../file.pdf",
			sizeBytes: 1024,
			wantErr:   s3validator.ErrUnsafeS3KeyPath,
		},
		{
			name:      "неподдерживаемое расширение txt",
			key:       "attachments/course-1/readme.txt",
			sizeBytes: 1024,
			wantErr:   content.ErrInvalidFormat,
		},
		{
			name:      "нет расширения",
			key:       "attachments/course-1/file",
			sizeBytes: 1024,
			wantErr:   content.ErrInvalidFormat,
		},
		{
			name:      "пустое имя файла с расширением",
			key:       "attachments/course-1/.pdf",
			sizeBytes: 1024,
			wantErr:   content.ErrInvalidFormat,
		},
		{
			name:      "нулевой размер",
			key:       "attachments/course-1/file.pdf",
			sizeBytes: 0,
			wantErr:   content.ErrEmptyFileSize,
		},
		{
			name:      "отрицательный размер",
			key:       "attachments/course-1/file.pdf",
			sizeBytes: -1,
			wantErr:   content.ErrEmptyFileSize,
		},
		{
			name:      "размер больше 700 мб",
			key:       "attachments/course-1/file.pdf",
			sizeBytes: 700*1024*1024 + 1,
			wantErr:   content.ErrTooLargeFile,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			_, err := content.NewAttachmentContent(tt.key, tt.sizeBytes)

			//Assert
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
