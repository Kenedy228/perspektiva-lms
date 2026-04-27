package element

import (
	"testing"

	"gitflic.ru/lms/internal/domain/shared/s3validator"
	"github.com/stretchr/testify/assert"
)

func TestValidateSlidesFile(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		sizeBytes int64
		want      error
	}{
		{
			name:      "валидный pptx",
			key:       "slides/course-1/intro.pptx",
			sizeBytes: 1024,
			want:      nil,
		},
		{
			name:      "валидный pptx uppercase extension",
			key:       "slides/course-1/intro.PPTX",
			sizeBytes: 2048,
			want:      nil,
		},
		{
			name:      "пустой ключ",
			key:       "",
			sizeBytes: 1024,
			want:      s3validator.ErrEmptyS3Key,
		},
		{
			name:      "невалидный s3 ключ",
			key:       "/slides/course-1/intro.pptx",
			sizeBytes: 1024,
			want:      s3validator.ErrInvalidS3Key,
		},
		{
			name:      "небезопасный path segment",
			key:       "slides/../intro.pptx",
			sizeBytes: 1024,
			want:      s3validator.ErrUnsafeS3KeyPath,
		},
		{
			name:      "неверное расширение pdf",
			key:       "slides/course-1/intro.pdf",
			sizeBytes: 1024,
			want:      ErrInvalid,
		},
		{
			name:      "без расширения",
			key:       "slides/course-1/intro",
			sizeBytes: 1024,
			want:      ErrInvalid,
		},
		{
			name:      "пустое имя файла",
			key:       "slides/course-1/.pptx",
			sizeBytes: 1024,
			want:      ErrInvalid,
		},
		{
			name:      "нулевой размер",
			key:       "slides/course-1/intro.pptx",
			sizeBytes: 0,
			want:      ErrInvalid,
		},
		{
			name:      "отрицательный размер",
			key:       "slides/course-1/intro.pptx",
			sizeBytes: -1,
			want:      ErrInvalid,
		},
		{
			name:      "ровно 100 мб",
			key:       "slides/course-1/intro.pptx",
			sizeBytes: 100 * 1024 * 1024,
			want:      nil,
		},
		{
			name:      "больше 100 мб",
			key:       "slides/course-1/intro.pptx",
			sizeBytes: 100*1024*1024 + 1,
			want:      ErrInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Act
			err := validateSlidersFile(tt.key, tt.sizeBytes)

			//Assert
			assert.ErrorIs(t, err, tt.want)
		})
	}
}
