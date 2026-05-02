package education_test

import (
	"strings"
	"testing"

	"gitflic.ru/lms/internal/domain/person/profile/education"
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
				edu, err := education.New(tt.in)

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
				name:    "пустой текст является некорректным",
				in:      "",
				wantErr: education.ErrInvalid,
			},
			{
				name:    "текст без непробельных символов является некорректным",
				in:      "      ",
				wantErr: education.ErrInvalid,
			},
			{
				name:    "текст без непробельных символов (с управляющими пробельными последовательностями) является некорректным",
				in:      " \t \t \t ",
				wantErr: education.ErrInvalid,
			},
			{
				name:    "текст с количеством символов > лимита некорректный",
				in:      strings.Repeat("A", education.ValueCharsLimit+1),
				wantErr: education.ErrInvalid,
			},
		}

		for _, tt := range tc {
			t.Run(tt.name, func(t *testing.T) {
				// Arrange
				_, err := education.New(tt.in)

				// Assert
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			})
		}
	})
}
