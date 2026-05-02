package jobtitle_test

import (
	"strings"
	"testing"

	"gitflic.ru/lms/internal/domain/person/profile/jobtitle"
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
				jt, err := jobtitle.New(tt.in)

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
				name:    "пустой текст является некорректным",
				in:      "",
				wantErr: jobtitle.ErrInvalid,
			},
			{
				name:    "текст без непробельных символов является некорректным",
				in:      "      ",
				wantErr: jobtitle.ErrInvalid,
			},
			{
				name:    "текст без непробельных символов (с управляющими пробельными последовательностями) является некорректным",
				in:      " \t \t \t ",
				wantErr: jobtitle.ErrInvalid,
			},
			{
				name:    "текст с количеством символов > лимита некорректный",
				in:      strings.Repeat("A", jobtitle.JobTitleCharsLimit+1),
				wantErr: jobtitle.ErrInvalid,
			},
		}

		for _, tt := range tc {
			t.Run(tt.name, func(t *testing.T) {
				// Arrange
				_, err := jobtitle.New(tt.in)

				// Assert
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			})
		}
	})
}
