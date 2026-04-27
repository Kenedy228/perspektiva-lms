package element_test

import (
	"strings"
	"testing"

	"gitflic.ru/lms/internal/domain/element"
	"gitflic.ru/lms/internal/domain/shared/s3validator"
	"github.com/stretchr/testify/assert"
)

func TestNewSlidesContent_Success(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		sizeBytes int64
		wantErr   error
	}{
		{
			name:      "успешное создание slides content",
			key:       "slides/course-1/intro.pptx",
			sizeBytes: 1024,
			wantErr:   nil,
		},
		{
			name:      "успешное создание slides content с uppercase расширением",
			key:       "slides/course-1/intro.PPTX",
			sizeBytes: 2048,
			wantErr:   nil,
		},
		{
			name:      "успешное создание если размер ровно 100 мб",
			key:       "slides/course-1/intro.pptx",
			sizeBytes: 100 * 1024 * 1024,
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			got, err := element.NewSlidesContent(tt.key, tt.sizeBytes)

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, got.Key(), tt.key)
			assert.Equal(t, got.SizeBytes(), tt.sizeBytes)
		})
	}
}

func TestNewSlidesContent_Fail(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		sizeBytes int64
		wantErr   error
	}{
		{
			name:      "ошибка если ключ пустой",
			key:       "",
			sizeBytes: 1024,
			wantErr:   s3validator.ErrEmptyS3Key,
		},
		{
			name:      "ошибка если ключ из пробелов",
			key:       "   ",
			sizeBytes: 1024,
			wantErr:   s3validator.ErrEmptyS3Key,
		},
		{
			name:      "ошибка если ключ слишком длинный",
			key:       strings.Repeat("a", 1025) + ".pptx",
			sizeBytes: 1024,
			wantErr:   s3validator.ErrTooLongS3Key,
		},
		{
			name:      "ошибка если ключ невалидный по формату",
			key:       "/slides/course-1/intro.pptx",
			sizeBytes: 1024,
			wantErr:   s3validator.ErrInvalidS3Key,
		},
		{
			name:      "ошибка если ключ содержит unsafe path segment",
			key:       "slides/../intro.pptx",
			sizeBytes: 1024,
			wantErr:   s3validator.ErrUnsafeS3KeyPath,
		},
		{
			name:      "ошибка если расширение не pptx",
			key:       "slides/course-1/intro.pdf",
			sizeBytes: 1024,
			wantErr:   element.ErrInvalid,
		},
		{
			name:      "ошибка если расширение отсутствует",
			key:       "slides/course-1/intro",
			sizeBytes: 1024,
			wantErr:   element.ErrInvalid,
		},
		{
			name:      "ошибка если имя файла пустое",
			key:       "slides/course-1/.pptx",
			sizeBytes: 1024,
			wantErr:   element.ErrInvalid,
		},
		{
			name:      "ошибка если размер файла нулевой",
			key:       "slides/course-1/intro.pptx",
			sizeBytes: 0,
			wantErr:   element.ErrInvalid,
		},
		{
			name:      "ошибка если размер файла отрицательный",
			key:       "slides/course-1/intro.pptx",
			sizeBytes: -1,
			wantErr:   element.ErrInvalid,
		},
		{
			name:      "ошибка если размер больше 100 мб",
			key:       "slides/course-1/intro.pptx",
			sizeBytes: 100*1024*1024 + 1,
			wantErr:   element.ErrInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			_, err := element.NewSlidesContent(tt.key, tt.sizeBytes)

			//Assert
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
