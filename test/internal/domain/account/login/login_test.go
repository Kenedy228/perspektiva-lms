package login_test

import (
	"strings"
	"testing"

	"gitflic.ru/lms/internal/domain/account/login"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_Success(t *testing.T) {
	tc := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "простой валидный логин",
			in:   "student2026",
			want: "student2026",
		},
		{
			name: "логин с точками и дефисами",
			in:   "ivan.ivanov-99",
			want: "ivan.ivanov-99",
		},
		{
			name: "логин со знаком подчеркивания",
			in:   "admin_user",
			want: "admin_user",
		},
		{
			name: "приведение к нижнему регистру",
			in:   "DaNil_Zh",
			want: "danil_zh",
		},
		{
			name: "удаление пробелов по краям",
			in:   "   petrov.a   ",
			want: "petrov.a",
		},
		{
			name: "удаление пробелов внутри строки",
			in:   "ivan  ivanov",
			want: "ivanivanov",
		},
		{
			name: "минимально допустимая длина (4 символа)",
			in:   "dany",
			want: "dany",
		},
		{
			name: "максимально допустимая длина (30 символов)",
			in:   strings.Repeat("a", 30),
			want: strings.Repeat("a", 30),
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			l, err := login.New(tt.in)

			//Assert
			require.NoError(t, err)
			assert.Equal(t, tt.want, l.Value())
		})
	}
}

func TestNew_Fail(t *testing.T) {
	tc := []struct {
		name string
		in   string
	}{
		{
			name: "пустая строка",
			in:   "",
		},
		{
			name: "только пробелы",
			in:   "      ",
		},
		{
			name: "меньше минимальной длины (3 символа)",
			in:   "dan",
		},
		{
			name: "больше максимальной длины (31 символ)",
			in:   strings.Repeat("a", 31),
		},
		{
			name: "русские буквы (кириллица)",
			in:   "данил_журавлев",
		},
		{
			name: "недопустимые спецсимволы (!@#)",
			in:   "student@123",
		},
		{
			name: "недопустимые символы (знак равенства)",
			in:   "admin=true",
		},
		{
			name: "недопустимые символы (слеши)",
			in:   "user/name",
		},
		{
			name: "эмодзи",
			in:   "admin👨‍💻",
		},
		{
			name: "длина становится меньше минимальной после нормализации",
			// "d a n" -> "dan" (3 символа, что меньше 4)
			in: "d a n",
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			_, err := login.New(tt.in)

			//Assert
			require.Error(t, err)
			assert.ErrorIs(t, err, login.ErrInvalid)
		})
	}
}
