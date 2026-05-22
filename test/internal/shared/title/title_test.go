//go:build legacy
// +build legacy

package title_test

import (
	"strings"
	"testing"

	title2 "gitflic.ru/lms/backend/internal/domain/shared/title"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("успех", func(t *testing.T) {
		tc := []struct {
			name      string
			value     string
			wantValue string
		}{
			{
				name:      "корректный кейс",
				value:     "заголовок вопроса",
				wantValue: "заголовок вопроса",
			},
			{
				name:      "при создании объекта удаляет пробелы по краям",
				value:     "   заголовок вопроса   ",
				wantValue: "заголовок вопроса",
			},
		}

		for _, tt := range tc {
			t.Run(tt.name, func(t *testing.T) {
				// Arrange
				got, err := title2.New(tt.value)

				// Assert
				assert.NoError(t, err)
				assert.Equal(t, tt.wantValue, got.Value())
			})
		}
	})

	t.Run("ошибка", func(t *testing.T) {
		tc := []struct {
			name    string
			value   string
			wantErr error
		}{
			{
				name:    "пустой заголовок",
				value:   "",
				wantErr: title2.ErrInvalid,
			},
			{
				name:    "заголовок из пробелов",
				value:   "",
				wantErr: title2.ErrInvalid,
			},
			{
				name:    "длина заголовка больше лимита",
				value:   strings.Repeat("a", title2.ValueCharsLimit+1),
				wantErr: title2.ErrInvalid,
			},
		}

		for _, tt := range tc {
			t.Run(tt.name, func(t *testing.T) {
				// Arrange
				_, err := title2.New(tt.value)

				// Assert
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			})
		}
	})
}
