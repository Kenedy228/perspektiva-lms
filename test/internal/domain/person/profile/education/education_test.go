//go:build legacy
// +build legacy

package education_test

import (
	"strings"
	"testing"

	education2 "gitflic.ru/lms/backend/internal/domain/person/profile/education"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("успех", func(t *testing.T) {
		tc := []struct {
			name string
			in   string
			want string
		}{
			{
				name: "при создании удаляет незначащие пробелы в начале и конце значения",
				in:   " высшее технической образование  ",
				want: "высшее технической образование",
			},
			{
				name: "при создании удаляет пробелы внутри текста",
				in:   "   высшее\t\t\tтехнической\rобразование  ",
				want: "высшее технической образование",
			},
		}

		for _, tt := range tc {
			t.Run(tt.name, func(t *testing.T) {
				// Arrange
				edu, err := education2.New(tt.in)

				// Assert
				assert.NoError(t, err)
				assert.Equal(t, tt.want, edu.Value())
			})
		}
	})

	t.Run("ошибка", func(t *testing.T) {
		tc := []struct {
			name    string
			in      string
			wantErr error
		}{
			{
				name:    "текст с количеством символов > лимита некорректный",
				in:      strings.Repeat("A", education2.ValueCharsLimit+1),
				wantErr: education2.ErrInvalid,
			},
		}

		for _, tt := range tc {
			t.Run(tt.name, func(t *testing.T) {
				// Arrange
				_, err := education2.New(tt.in)

				// Assert
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			})
		}
	})

	t.Run("успех: пустое значение допустимо", func(t *testing.T) {
		tc := []struct {
			name string
			in   string
		}{
			{name: "пустая строка", in: ""},
			{name: "только пробелы", in: "      "},
			{name: "управляющие пробельные последовательности", in: " \t \t \t "},
		}

		for _, tt := range tc {
			t.Run(tt.name, func(t *testing.T) {
				edu, err := education2.New(tt.in)
				assert.NoError(t, err)
				assert.Equal(t, "", edu.Value())
			})
		}
	})
}
