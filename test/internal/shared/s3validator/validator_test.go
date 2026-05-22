//go:build legacy
// +build legacy

package s3validator_test

import (
	"strings"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/shared/s3validator"
	"github.com/stretchr/testify/assert"
)

func TestValidateS3Key_Success(t *testing.T) {
	tc := []struct {
		name string
		key  string
	}{
		{
			name: "валидный простой ключ",
			key:  "slides/course-1/intro.pptx",
		},
		{
			name: "валидный ключ с unicode",
			key:  "курсы/модуль-1/презентация.pptx",
		},
		{
			name: "валидный ключ с точками и дефисами",
			key:  "courses/module-1/v1.0/final-presentation.pdf",
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Act
			err := s3validator.ValidateS3Key(tt.key)

			//Assert
			assert.NoError(t, err)
		})
	}
}

func TestValidateS3Key_ErrEmptyS3Key(t *testing.T) {
	tc := []struct {
		name string
		key  string
	}{
		{
			name: "пустая строка",
			key:  "",
		},
		{
			name: "строка из пробелов",
			key:  "   ",
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Act
			err := s3validator.ValidateS3Key(tt.key)

			//Assert
			assert.ErrorIs(t, err, s3validator.ErrEmptyS3Key)
		})
	}
}

func TestValidateS3Key_ErrTooLongS3Key(t *testing.T) {
	tc := []struct {
		name string
		key  string
	}{
		{
			name: "слишком длинный ключ ascii",
			key:  strings.Repeat("a", 1e5),
		},
		{
			name: "слишком длинный unicode ключ в байтах",
			key:  strings.Repeat("й", 1e5),
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Act
			err := s3validator.ValidateS3Key(tt.key)

			//Assert
			assert.ErrorIs(t, err, s3validator.ErrTooLongS3Key)
		})
	}
}

func TestValidateS3Key_ErrInvalidS3Key(t *testing.T) {
	tc := []struct {
		name string
		key  string
	}{
		{
			name: "начинается со слеша",
			key:  "/slides/course-1/intro.pptx",
		},
		{
			name: "заканчивается слешем",
			key:  "slides/course-1/",
		},
		{
			name: "двойной слеш",
			key:  "slides//course-1/intro.pptx",
		},
		{
			name: "пробел внутри ключа",
			key:  "slides/course 1/intro.pptx",
		},
		{
			name: "tab внутри ключа",
			key:  "slides/\tintro.pptx",
		},
		{
			name: "перевод строки внутри ключа",
			key:  "slides/\nintro.pptx",
		},
		{
			name: "carriage return внутри ключа",
			key:  "slides/\rintro.pptx",
		},
		{
			name: "backslash в ключе",
			key:  `slides\course-1\intro.pptx`,
		},
		{
			name: "управляющий символ",
			key:  "slides/" + string(rune(0x01)) + "intro.pptx",
		},
		{
			name: "невалидный utf8 ключ",
			key:  string([]byte{0xff, 0xfe, 0xfd}),
		},
		{
			name: "содержит двойной backslash",
			key:  `slides\\course-1\\intro.pptx`,
		},
		{
			name: "состоит из slash",
			key:  "/",
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Act
			err := s3validator.ValidateS3Key(tt.key)

			//Assert
			assert.ErrorIs(t, err, s3validator.ErrInvalidS3Key)
		})
	}
}

func TestValidateS3Key_ErrUnsafeS3KeyPath(t *testing.T) {
	tc := []struct {
		name string
		key  string
	}{
		{
			name: "сегмент точка",
			key:  "slides/./intro.pptx",
		},
		{
			name: "сегмент две точки",
			key:  "slides/../intro.pptx",
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Act
			err := s3validator.ValidateS3Key(tt.key)

			//Assert
			assert.ErrorIs(t, err, s3validator.ErrUnsafeS3KeyPath)
		})
	}
}
