//go:build legacy
// +build legacy

package jobtitle_test

import (
	"strings"
	"testing"

	jobtitle2 "gitflic.ru/lms/backend/internal/domain/person/profile/jobtitle"
	"github.com/stretchr/testify/assert"
)

func TestNew_Success(t *testing.T) {
	t.Run("успех", func(t *testing.T) {
		tc := []struct {
			name string
			in   string
			want string
		}{
			{
				name: "при создании удаляет незначащие пробелы в начале и конце значения",
				in:   " инженер первой категории  ",
				want: "инженер первой категории",
			},
			{
				name: "при создании удаляет пробелы внутри текста",
				in:   "   инженер\t\t\tпервой\rкатегории   ",
				want: "инженер первой категории",
			},
		}

		for _, tt := range tc {
			t.Run(tt.name, func(t *testing.T) {
				// Arrange
				jt, err := jobtitle2.New(tt.in)

				// Assert
				assert.NoError(t, err)
				assert.Equal(t, tt.want, jt.Title())
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
				in:      strings.Repeat("A", jobtitle2.JobTitleCharsLimit+1),
				wantErr: jobtitle2.ErrInvalid,
			},
		}

		for _, tt := range tc {
			t.Run(tt.name, func(t *testing.T) {
				// Arrange
				_, err := jobtitle2.New(tt.in)

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
				jt, err := jobtitle2.New(tt.in)
				assert.NoError(t, err)
				assert.Equal(t, "", jt.Title())
			})
		}
	})
}
