//go:build legacy
// +build legacy

package snils_test

import (
	"testing"

	snils2 "gitflic.ru/lms/backend/internal/domain/person/profile/snils"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("успех", func(t *testing.T) {
		tc := []struct {
			name          string
			input         string
			wantPlain     string
			wantFormatted string
		}{
			{
				name:          "валидный форматированный ввод",
				input:         "112-233-445 95",
				wantPlain:     "11223344595",
				wantFormatted: "112-233-445 95",
			},
			{
				name:          "валидный ввод без форматирования",
				input:         "11223344595",
				wantPlain:     "11223344595",
				wantFormatted: "112-233-445 95",
			},
			{
				name:          "старый снилс ниже порога проверки контрольной суммы",
				input:         "00100199800",
				wantPlain:     "00100199800",
				wantFormatted: "001-001-998 00",
			},
		}

		for _, tt := range tc {
			t.Run(tt.name, func(t *testing.T) {
				// Arrange
				s, err := snils2.New(tt.input)

				// Assert
				assert.NoError(t, err)
				assert.Equal(t, tt.wantPlain, s.Value())
				assert.Equal(t, tt.wantFormatted, s.Formatted())
			})
		}
	})

	t.Run("ошибка", func(t *testing.T) {
		tc := []struct {
			name    string
			input   string
			wantErr error
		}{
			{
				name:    "невалидная контрольная сумма",
				input:   "112-233-445 00",
				wantErr: snils2.ErrInvalid,
			},
			{
				name:    "недопустимый символ в отформатированном значении",
				input:   "112-233-445 9A",
				wantErr: snils2.ErrInvalid,
			},
			{
				name:    "слишком короткий ввод",
				input:   "112-233-445 9",
				wantErr: snils2.ErrInvalid,
			},
			{
				name:    " пустая строка",
				input:   "",
				wantErr: snils2.ErrInvalid,
			},
		}

		for _, tt := range tc {
			t.Run(tt.name, func(t *testing.T) {
				// Arrange
				_, err := snils2.New(tt.input)

				// Assert
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			})
		}
	})
}
