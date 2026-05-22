//go:build legacy
// +build legacy

package inn_test

import (
	"testing"

	inn2 "gitflic.ru/lms/backend/internal/domain/organization/inn"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("успех", func(t *testing.T) {
		tc := []struct {
			name     string
			code     string
			kind     inn2.Type
			wantCode string
		}{
			{
				name:     "валидный ИНН юридического лица (10 цифр)",
				code:     "7728168971", // Альфа-Банк (как пример корректной КС)
				kind:     inn2.TypeOrganization,
				wantCode: "7728168971",
			},
			{
				name:     "валидный ИНН юридического лица с пробелами",
				code:     " 7728 1689 71 ",
				kind:     inn2.TypeOrganization,
				wantCode: "7728168971",
			},
			{
				name:     "валидный ИНН индивидуального предпринимателя (12 цифр)",
				code:     "500100732259", // случайный валидный 12-значный
				kind:     inn2.TypeIP,
				wantCode: "500100732259",
			},
			{
				name:     "валидный ИНН физического лица (12 цифр)",
				code:     "500100732259",
				kind:     inn2.TypePhysical,
				wantCode: "500100732259",
			},
			{
				name:     "валидный ИНН юрлица, где остаток контрольной суммы равен 10 -> 0",
				code:     "1030000000",
				kind:     inn2.TypeOrganization,
				wantCode: "1030000000",
			},
		}

		for _, tt := range tc {
			t.Run(tt.name, func(t *testing.T) {
				// Act
				i, err := inn2.New(tt.code, tt.kind)

				// Assert
				assert.NoError(t, err)
				assert.Equal(t, tt.wantCode, i.Code())
				assert.Equal(t, tt.kind, i.Type())
			})
		}
	})

	t.Run("ошибка", func(t *testing.T) {
		tc := []struct {
			name    string
			code    string
			kind    inn2.Type
			wantErr error
		}{
			{
				name:    "пустая строка",
				code:    "",
				kind:    inn2.TypeOrganization,
				wantErr: inn2.ErrInvalid,
			},
			{
				name:    "только пробелы",
				code:    "          ",
				kind:    inn2.TypeOrganization,
				wantErr: inn2.ErrInvalid,
			},
			{
				name:    "содержит буквы",
				code:    "772816897A",
				kind:    inn2.TypeOrganization,
				wantErr: inn2.ErrInvalid,
			},
			{
				name:    "содержит спецсимволы",
				code:    "7728-16897",
				kind:    inn2.TypeOrganization,
				wantErr: inn2.ErrInvalid,
			},
			{
				name:    "неверная длина для юрлица (9 цифр)",
				code:    "772816897",
				kind:    inn2.TypeOrganization,
				wantErr: inn2.ErrInvalid,
			},
			{
				name:    "неверная длина для юрлица (11 цифр)",
				code:    "77281689711",
				kind:    inn2.TypeOrganization,
				wantErr: inn2.ErrInvalid,
			},
			{
				name:    "неверная длина для ИП (11 цифр)",
				code:    "50010073225",
				kind:    inn2.TypeIP,
				wantErr: inn2.ErrInvalid,
			},
			{
				name:    "неверная длина для физлица (10 цифр)",
				code:    "7728168971",
				kind:    inn2.TypePhysical,
				wantErr: inn2.ErrInvalid,
			},
			{
				name:    "юрлицо: неверная контрольная сумма",
				code:    "7728168972", // Правильная 1, заменили на 2
				kind:    inn2.TypeOrganization,
				wantErr: inn2.ErrInvalid,
			},
			{
				name:    "ИП: неверная первая контрольная цифра",
				code:    "500100732249", // Правильная 5, заменили на 4
				kind:    inn2.TypeIP,
				wantErr: inn2.ErrInvalid,
			},
			{
				name:    "ИП: неверная вторая контрольная цифра",
				code:    "500100732258", // Правильная 9, заменили на 8
				kind:    inn2.TypeIP,
				wantErr: inn2.ErrInvalid,
			},
			{
				name:    "физлицо: неверная контрольная сумма",
				code:    "500100732258",
				kind:    inn2.TypePhysical,
				wantErr: inn2.ErrInvalid,
			},
			{
				name:    "неизвестный тип ИНН",
				code:    "7728168971",
				kind:    inn2.Type("unknown"),
				wantErr: inn2.ErrInvalid,
			},
		}

		for _, tt := range tc {
			t.Run(tt.name, func(t *testing.T) {
				// Act
				_, err := inn2.New(tt.code, tt.kind)

				// Assert
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			})
		}
	})
}
