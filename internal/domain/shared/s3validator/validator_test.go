package s3validator_test

import (
	"strings"
	"testing"

	"gitflic.ru/lms/internal/domain/shared/s3validator"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
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
			name: "валидный ключ с подчеркиванием",
			key:  "courses/module_1/final_presentation.pdf",
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
			name: "слишком длинный ключ ascii",
			key:  strings.Repeat("a", 1025),
			want: s3validator.ErrTooLongS3Key,
		},
		{
			name: "слишком длинный unicode ключ в байтах",
			key:  strings.Repeat("й", 600),
			want: s3validator.ErrTooLongS3Key,
		},
		{
			name: "невалидный utf8 ключ",
			key:  string([]byte{0xff, 0xfe, 0xfd}),
			want: s3validator.ErrInvalidS3Key,
		},
		{
			name: "содержит двойной слеш",
			key:  "slides//course-1/intro.pptx",
			want: s3validator.ErrInvalidS3Key,
		},
		{
			name: "содержит двойной backslash",
			key:  `slides\\course-1\\intro.pptx`,
			want: s3validator.ErrInvalidS3Key,
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
			name: "пробел внутри сегмента",
			key:  "slides/course 1/intro.pptx",
			want: s3validator.ErrInvalidS3Key,
		},
		{
			name: "tab внутри сегмента",
			key:  "slides/\tintro.pptx",
			want: s3validator.ErrInvalidS3Key,
		},
		{
			name: "newline внутри сегмента",
			key:  "slides/\nintro.pptx",
			want: s3validator.ErrInvalidS3Key,
		},
		{
			name: "carriage return внутри сегмента",
			key:  "slides/\rintro.pptx",
			want: s3validator.ErrInvalidS3Key,
		},
		{
			name: "контрольный символ внутри сегмента",
			key:  "slides/" + string(rune(0x01)) + "intro.pptx",
			want: s3validator.ErrInvalidS3Key,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s3validator.ValidateS3Key(tt.key)

			assert.ErrorIs(t, err, tt.want)
		})
	}
}
