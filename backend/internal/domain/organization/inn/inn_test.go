package inn

import (
	"reflect"
	"strings"
	"testing"
)

func TestINN_IsZero(t *testing.T) {
	type fields struct {
		value string
		t     Type
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "нет value",
			fields: fields{
				value: "",
				t:     TypeNaturalPerson,
			},
			want: true,
		},
		{
			name: "нет типа",
			fields: fields{
				value: "123",
				t:     "",
			},
			want: true,
		},
		{
			name: "нет типа и value",
			fields: fields{
				value: "",
				t:     "",
			},
			want: true,
		},
		{
			name: "есть тип и value",
			fields: fields{
				value: "123",
				t:     TypeIndividualEntrepreneur,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inn := INN{
				value: tt.fields.value,
				t:     tt.fields.t,
			}
			if got := inn.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestINN_Type(t *testing.T) {
	type fields struct {
		value string
		t     Type
	}
	tests := []struct {
		name   string
		fields fields
		want   Type
	}{
		{
			name: "возвращает как есть",
			fields: fields{
				value: "",
				t:     TypeLegalEntity,
			},
			want: TypeLegalEntity,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inn := INN{
				value: tt.fields.value,
				t:     tt.fields.t,
			}
			if got := inn.Type(); got != tt.want {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestINN_Value(t *testing.T) {
	type fields struct {
		value string
		t     Type
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "возвращает как есть",
			fields: fields{
				value: "123",
				t:     TypeLegalEntity,
			},
			want: "123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inn := INN{
				value: tt.fields.value,
				t:     tt.fields.t,
			}
			if got := inn.Value(); got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		value string
		t     Type
	}
	tests := []struct {
		name    string
		args    args
		want    INN
		wantErr bool
	}{
		{
			name: "пустая строка",
			args: args{
				value: "",
				t:     TypeLegalEntity,
			},
			want: INN{
				value: "",
				t:     "",
			},
			wantErr: true,
		},
		{
			name: "только пробелы",
			args: args{
				value: strings.Repeat(" ", 10),
				t:     TypeLegalEntity,
			},
			want: INN{
				value: "",
				t:     "",
			},
			wantErr: true,
		},
		{
			name: "содержит нецифры",
			args: args{
				value: "123adsb-cds",
				t:     TypeLegalEntity,
			},
			want: INN{
				value: "",
				t:     "",
			},
			wantErr: true,
		},
		{
			name: "запрещенное значение для ЮЛ",
			args: args{
				value: strings.Repeat("0", 10),
				t:     TypeLegalEntity,
			},
			want: INN{
				value: "",
				t:     "",
			},
			wantErr: true,
		},
		{
			name: "валидный ИНН ЮЛ",
			args: args{
				value: "1030000000",
				t:     TypeLegalEntity,
			},
			want: INN{
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
			want: INN{
				value: "",
				t:     "",
			},
			wantErr: true,
		},
		{
			name: "валидный ИНН для ФЛ",
			args: args{
				value: "500100732259",
				t:     TypeNaturalPerson,
			},
			want: INN{
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
			want: INN{
				value: "",
				t:     "",
			},
			wantErr: true,
		},
		{
			name: "неверная вторая контрольная цифра ФЛ",
			args: args{
				value: "500100732258", // правильная 9, заменили на 8
				t:     TypeNaturalPerson,
			},
			want: INN{
				value: "",
				t:     "",
			},
			wantErr: true,
		},
		{
			name: "валидный ИНН для ИП",
			args: args{
				value: "500100732259",
				t:     TypeIndividualEntrepreneur,
			},
			want: INN{
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
			want: INN{
				value: "",
				t:     "",
			},
			wantErr: true,
		},
		{
			name: "неверная вторая контрольная цифра ИП",
			args: args{
				value: "500100732258", // правильная 9, заменили на 8
				t:     TypeIndividualEntrepreneur,
			},
			want: INN{
				value: "",
				t:     "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.value, tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
		})
	}
}
