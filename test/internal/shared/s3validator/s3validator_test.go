package s3validator_test

import (
	"strings"
	"testing"

	"gitflic.ru/lms/internal/domain/shared/s3validator"
	"github.com/stretchr/testify/assert"
)

func TestValidateS3Key(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want error
	}{
		{
			name: "валидный простой ключ",
			key:  "slides/course-1/intro.pptx",
			want: nil,
		},
		{
			name: "валидный ключ с unicode",
			key:  "курсы/модуль-1/презентация.pptx",
			want: nil,
		},
		{
			name: "валидный ключ с точками и дефисами",
			key:  "courses/module-1/v1.0/final-presentation.pdf",
			want: nil,
		},
		{
			name: "пустая строка",
			key:  "",
			want: s3validator.ErrEmptyS3Key,
		},
		{
			name: "строка из пробелов",
			key:  "   ",
			want: s3validator.ErrEmptyS3Key,
		},
		{
			name: "слишком длинный ключ",
			key:  strings.Repeat("a", 1025),
			want: s3validator.ErrTooLongS3Key,
		},
		{
			name: "начинается со слеша",
			key:  "/slides/course-1/intro.pptx",
			want: s3validator.ErrInvalidS3Key,
		},
		{
			name: "заканчивается слешем",
			key:  "slides/course-1/",
			want: s3validator.ErrInvalidS3Key,
		},
		{
			name: "двойной слеш",
			key:  "slides//course-1/intro.pptx",
			want: s3validator.ErrInvalidS3Key,
		},
		{
			name: "сегмент точка",
			key:  "slides/./intro.pptx",
			want: s3validator.ErrUnsafeS3KeyPath,
		},
		{
			name: "сегмент две точки",
			key:  "slides/../intro.pptx",
			want: s3validator.ErrUnsafeS3KeyPath,
		},
		{
			name: "пробел внутри ключа",
			key:  "slides/course 1/intro.pptx",
			want: s3validator.ErrInvalidS3Key,
		},
		{
			name: "tab внутри ключа",
			key:  "slides/\tintro.pptx",
			want: s3validator.ErrInvalidS3Key,
		},
		{
			name: "перевод строки внутри ключа",
			key:  "slides/\nintro.pptx",
			want: s3validator.ErrInvalidS3Key,
		},
		{
			name: "carriage return внутри ключа",
			key:  "slides/\rintro.pptx",
			want: s3validator.ErrInvalidS3Key,
		},
		{
			name: "backslash в ключе",
			key:  `slides\course-1\intro.pptx`,
			want: s3validator.ErrInvalidS3Key,
		},
		{
			name: "управляющий символ",
			key:  "slides/" + string(rune(0x01)) + "intro.pptx",
			want: s3validator.ErrInvalidS3Key,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Act
			err := s3validator.ValidateS3Key(tt.key)

			//Assert
			assert.ErrorIs(t, err, tt.want)
		})
	}
}
