package inn_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/organization/inn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_Success(t *testing.T) {
	tc := []struct {
		name   string
		inCode string
		inType inn.Type
		want   string
	}{
		{
			name:   "валидный ИНН юридического лица (10 цифр)",
			inCode: "7728168971", // Альфа-Банк (как пример корректной КС)
			inType: inn.TypeOrganization,
			want:   "7728168971",
		},
		{
			name:   "валидный ИНН юридического лица с пробелами",
			inCode: " 7728 1689 71 ",
			inType: inn.TypeOrganization,
			want:   "7728168971",
		},
		{
			name:   "валидный ИНН индивидуального предпринимателя (12 цифр)",
			inCode: "500100732259", // случайный валидный 12-значный
			inType: inn.TypeIP,
			want:   "500100732259",
		},
		{
			name:   "валидный ИНН физического лица (12 цифр)",
			inCode: "500100732259",
			inType: inn.TypePhysical,
			want:   "500100732259",
		},
		{
			name:   "валидный ИНН юрлица, где остаток контрольной суммы равен 10 -> 0",
			inCode: "1030000000",
			inType: inn.TypeOrganization,
			want:   "1030000000",
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			i, err := inn.New(tt.inCode, tt.inType)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, tt.want, i.Code())
			assert.Equal(t, tt.inType, i.Type())
		})
	}
}

func TestNew_Fail(t *testing.T) {
	tc := []struct {
		name   string
		inCode string
		inType inn.Type
	}{
		{
			name:   "пустая строка",
			inCode: "",
			inType: inn.TypeOrganization,
		},
		{
			name:   "только пробелы",
			inCode: "          ",
			inType: inn.TypeOrganization,
		},
		{
			name:   "содержит буквы",
			inCode: "772816897A",
			inType: inn.TypeOrganization,
		},
		{
			name:   "содержит спецсимволы",
			inCode: "7728-16897",
			inType: inn.TypeOrganization,
		},
		{
			name:   "неверная длина для юрлица (9 цифр)",
			inCode: "772816897",
			inType: inn.TypeOrganization,
		},
		{
			name:   "неверная длина для юрлица (11 цифр)",
			inCode: "77281689711",
			inType: inn.TypeOrganization,
		},
		{
			name:   "неверная длина для ИП (11 цифр)",
			inCode: "50010073225",
			inType: inn.TypeIP,
		},
		{
			name:   "неверная длина для физлица (10 цифр)",
			inCode: "7728168971",
			inType: inn.TypePhysical,
		},
		{
			name:   "юрлицо: неверная контрольная сумма",
			inCode: "7728168972", // Правильная 1, заменили на 2
			inType: inn.TypeOrganization,
		},
		{
			name:   "ИП: неверная первая контрольная цифра",
			inCode: "500100732249", // Правильная 5, заменили на 4
			inType: inn.TypeIP,
		},
		{
			name:   "ИП: неверная вторая контрольная цифра",
			inCode: "500100732258", // Правильная 9, заменили на 8
			inType: inn.TypeIP,
		},
		{
			name:   "физлицо: неверная контрольная сумма",
			inCode: "500100732258",
			inType: inn.TypePhysical,
		},
		{
			name:   "неизвестный тип ИНН",
			inCode: "7728168971",
			inType: inn.Type("unknown"),
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			_, err := inn.New(tt.inCode, tt.inType)

			// Assert
			require.Error(t, err)
			assert.ErrorIs(t, err, inn.ErrInvalid)
		})
	}
}
