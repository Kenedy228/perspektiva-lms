//go:build legacy
// +build legacy

package name_test

import (
	"strings"
	"testing"

	name2 "gitflic.ru/lms/backend/internal/domain/organization/name"
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
				name:      "при создании удаляет незначащие пробелы в начале и конце значения",
				value:     " ООО Ромашка  ",
				wantValue: "ООО Ромашка",
			},
			{
				name:      "при создании не удаляет пробелы внутри текста",
				value:     "   ООО  Ромашка   ",
				wantValue: "ООО  Ромашка",
			},
		}

		for _, tt := range tc {
			t.Run(tt.name, func(t *testing.T) {
				// Arrange
				n, err := name2.New(tt.value)

				// Assert
				assert.NoError(t, err)
				assert.Equal(t, tt.wantValue, n.Value())
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
				name:    "пустой текст является некорректным",
				value:   "",
				wantErr: name2.ErrInvalid,
			},
			{
				name:    "текст без непробельных символов является некорректным",
				value:   "      ",
				wantErr: name2.ErrInvalid,
			},
			{
				name:    "текст без непробельных символов (с управляющими пробельными последовательностями) является некорректным",
				value:   " \t \t \t ",
				wantErr: name2.ErrInvalid,
			},
			{
				name:    "текст с количеством символов > лимита некорректный",
				value:   strings.Repeat("A", name2.ValueCharsLimit+1),
				wantErr: name2.ErrInvalid,
			},
		}

		for _, tt := range tc {
			t.Run(tt.name, func(t *testing.T) {
				// Arrange
				_, err := name2.New(tt.value)

				// Assert
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			})
		}
	})
}
