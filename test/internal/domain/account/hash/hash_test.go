package hash_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/account/passhash"
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
			name: "валидный хеш bcrypt версии 2a",
			in:   "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
			want: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
		},
		{
			name: "валидный хеш bcrypt версии 2b",
			in:   "$2b$12$D24p4h.P6P4.82.jR.X1U.1Q6Qx4/G0iB2.JzH8H4w2rP/T5k0eZ2",
			want: "$2b$12$D24p4h.P6P4.82.jR.X1U.1Q6Qx4/G0iB2.JzH8H4w2rP/T5k0eZ2",
		},
		{
			name: "валидный хеш bcrypt версии 2y",
			in:   "$2y$08$9G1v.G7l1/0Xo.4oM5iX/OeD1/K7V8k9wB2rP/T5k0eZ2X1U.1Q6Q",
			want: "$2y$08$9G1v.G7l1/0Xo.4oM5iX/OeD1/K7V8k9wB2rP/T5k0eZ2X1U.1Q6Q",
		},
		{
			name: "удаляет пробелы по краям, но не трогает саму строку (в т.ч. регистр)",
			in:   "  $2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy  ",
			want: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			h, err := passhash.New(tt.in)

			//Assert
			require.NoError(t, err)
			assert.Equal(t, tt.want, h.Hash())
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
			in:   "       ",
		},
		{
			name: "обычный текст (не хеш)",
			in:   "my_super_password_123",
		},
		{
			name: "MD5/SHA1 хеш (без префикса bcrypt)",
			in:   "e10adc3949ba59abbe56e057f20f883e", // MD5 от "123456"
		},
		{
			name: "bcrypt, но неверная версия префикса",
			in:   "$2c$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // $2c$ не существует
		},
		{
			name: "bcrypt, но неверный формат Cost (параметр сложности)",
			in:   "$2a$1$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // "1" вместо двузначного числа
		},
		{
			name: "bcrypt, но слишком короткий",
			// Удалили часть символов в конце
			in: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p",
		},
		{
			name: "bcrypt, но слишком длинный",
			// Добавили символы в конец
			in: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWyAAAA",
		},
		{
			name: "bcrypt, но содержит недопустимые символы в Base64 части (например % или &)",
			in:   "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl%p92ldGxad68LJZdL17lhWy", // Вставили %
		},
		{
			name: "строка, которая после TrimSpace становится пустой",
			in:   " \t \n ",
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			_, err := passhash.New(tt.in)

			// Assert
			require.Error(t, err)
			assert.ErrorIs(t, err, passhash.ErrInvalid)
		})
	}
}
