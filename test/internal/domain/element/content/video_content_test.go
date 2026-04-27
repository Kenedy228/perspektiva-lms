package content_test

import (
	"strings"
	"testing"

	"gitflic.ru/lms/internal/domain/element/content"
	"gitflic.ru/lms/internal/domain/shared/s3validator"
	"github.com/stretchr/testify/assert"
)

func TestNewVideoContent(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		sizeBytes int64
		wantErr   error
	}{
		{
			name:      "валидное mp4 видео",
			key:       "videos/course-1/intro.mp4",
			sizeBytes: 1024,
			wantErr:   nil,
		},
		{
			name:      "валидное webm видео",
			key:       "videos/course-1/intro.webm",
			sizeBytes: 2048,
			wantErr:   nil,
		},
		{
			name:      "валидное видео с unicode в пути",
			key:       "видео/курс-1/введение.mp4",
			sizeBytes: 4096,
			wantErr:   nil,
		},
		{
			name:      "размер ровно 500 мб",
			key:       "videos/course-1/intro.mp4",
			sizeBytes: 500 * 1024 * 1024,
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			got, err := content.NewVideoContent(tt.key, tt.sizeBytes)

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, got.Key(), tt.key)
			assert.Equal(t, got.SizeBytes(), tt.sizeBytes)
		})
	}
}

func TestNewVideoContent_Fail(t *testing.T) {
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
			key:       strings.Repeat("a", 1025) + ".mp4",
			sizeBytes: 1024,
			wantErr:   s3validator.ErrTooLongS3Key,
		},
		{
			name:      "невалидный s3 ключ с ведущим слешем",
			key:       "/videos/course-1/intro.mp4",
			sizeBytes: 1024,
			wantErr:   s3validator.ErrInvalidS3Key,
		},
		{
			name:      "небезопасный path segment",
			key:       "videos/../intro.mp4",
			sizeBytes: 1024,
			wantErr:   s3validator.ErrUnsafeS3KeyPath,
		},
		{
			name:      "неподдерживаемое расширение avi",
			key:       "videos/course-1/intro.avi",
			sizeBytes: 1024,
			wantErr:   content.ErrInvalidFormat,
		},
		{
			name:      "нет расширения",
			key:       "videos/course-1/intro",
			sizeBytes: 1024,
			wantErr:   content.ErrInvalidFormat,
		},
		{
			name:      "пустое имя файла с расширением",
			key:       "videos/course-1/.mp4",
			sizeBytes: 1024,
			wantErr:   content.ErrInvalidFormat,
		},
		{
			name:      "нулевой размер файла",
			key:       "videos/course-1/intro.mp4",
			sizeBytes: 0,
			wantErr:   content.ErrEmptyFileSize,
		},
		{
			name:      "отрицательный размер файла",
			key:       "videos/course-1/intro.mp4",
			sizeBytes: -1,
			wantErr:   content.ErrEmptyFileSize,
		},
		{
			name:      "размер больше 500 мб",
			key:       "videos/course-1/intro.mp4",
			sizeBytes: 500*1024*1024 + 1,
			wantErr:   content.ErrTooLargeFile,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			_, err := content.NewVideoContent(tt.key, tt.sizeBytes)

			//Assert
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
