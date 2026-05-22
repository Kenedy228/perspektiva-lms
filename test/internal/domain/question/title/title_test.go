//go:build legacy
// +build legacy

package title_test

import (
	"strings"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/question/title"
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
				titl, err := title.New(tt.value)

				// Assert
				assert.NoError(t, err)
				assert.Equal(t, tt.wantValue, titl.Value())
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
				wantErr: title.ErrInvalid,
			},
			{
				name:    "заголовок из пробелов",
				value:   "",
				wantErr: title.ErrInvalid,
			},
			{
				name:    "длина заголовка больше лимита",
				value:   strings.Repeat("a", title.ValueCharsLimit+1),
				wantErr: title.ErrInvalid,
			},
		}

		for _, tt := range tc {
			t.Run(tt.name, func(t *testing.T) {
				// Arrange
				_, err := title.New(tt.value)

				// Assert
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			})
		}
	})
}
