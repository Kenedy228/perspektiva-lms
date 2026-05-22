package inn

import (
	"strings"
	"testing"
)

func Test_validateValue(t *testing.T) {
	type args struct {
		value string
		t     Type
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "пустая строка",
			args: args{
				value: "",
				t:     TypeLegalEntity,
			},
			wantErr: true,
		},
		{
			name: "только пробелы",
			args: args{
				value: strings.Repeat(" ", 10),
				t:     TypeLegalEntity,
			},
			wantErr: true,
		},
		{
			name: "содержит нецифры",
			args: args{
				value: "123adsb-cds",
				t:     TypeLegalEntity,
			},
			wantErr: true,
		},
		{
			name: "запрещенное значение для ЮЛ",
			args: args{
				value: strings.Repeat("0", 10),
				t:     TypeLegalEntity,
			},
			wantErr: true,
		},
		{
			name: "валидный ИНН ЮЛ",
			args: args{
				value: "1030000000",
				t:     TypeLegalEntity,
			},
			wantErr: false,
		},
		{
			name: "невалидный ИНН ЮЛ",
			args: args{
				value: "7728168972", // Правильная 1, заменили на 2
				t:     TypeLegalEntity,
			},
			wantErr: true,
		},
		{
			name: "валидный ИНН для ФЛ",
			args: args{
				value: "500100732259",
				t:     TypeNaturalPerson,
			},
			wantErr: false,
		},
		{
			name: "неверная вторая контрольная цифра ФЛ",
			args: args{
				value: "500100732249", // правильная 5, заменили на 4
				t:     TypeNaturalPerson,
			},
			wantErr: true,
		},
		{
			name: "неверная вторая контрольная цифра ФЛ",
			args: args{
				value: "500100732258", // правильная 9, заменили на 8
				t:     TypeNaturalPerson,
			},
			wantErr: true,
		},
		{
			name: "валидный ИНН для ИП",
			args: args{
				value: "500100732259",
				t:     TypeIndividualEntrepreneur,
			},
			wantErr: false,
		},
		{
			name: "неверная вторая контрольная цифра ИП",
			args: args{
				value: "500100732249", // правильная 5, заменили на 4
				t:     TypeIndividualEntrepreneur,
			},
			wantErr: true,
		},
		{
			name: "неверная вторая контрольная цифра ИП",
			args: args{
				value: "500100732258", // правильная 9, заменили на 8
				t:     TypeIndividualEntrepreneur,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateValue(tt.args.value, tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("validateValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
