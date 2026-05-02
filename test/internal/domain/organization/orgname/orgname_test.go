package orgname_test

import (
	"strings"
	"testing"

	"gitflic.ru/lms/internal/domain/organization/orgname"
	"github.com/stretchr/testify/assert"
)

func TestNew_Success(t *testing.T) {
	tc := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "при создании удаляет незначащие пробелы в начале и конце значения",
			in:   " ООО Ромашка  ",
			want: "ООО Ромашка",
		},
		{
			name: "при создании не удаляет пробелы внутри текста",
			in:   "   ООО  Ромашка   ",
			want: "ООО  Ромашка",
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			name, err := orgname.New(tt.in)

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, name.Value(), tt.want)
		})
	}
}

func TestNew_Fail(t *testing.T) {
	tc := []struct {
		name string
		in   string
	}{
		{
			name: "пустой текст является некорректным",
			in:   "",
		},
		{
			name: "текст без непробельных символов является некорректным",
			in:   "      ",
		},
		{
			name: "текст без непробельных символов (с управляющими пробельными последовательностями) является некорректным",
			in:   " \t \t \t ",
		},
		{
			name: "текст с количеством символов > лимита некорректный",
			in:   strings.Repeat("A", 1e5),
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			_, err := orgname.New(tt.in)

			//Assert
			assert.ErrorIs(t, err, orgname.ErrInvalid)
		})
	}
}
