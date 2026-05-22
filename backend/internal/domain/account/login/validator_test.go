package login

import (
	"strings"
	"testing"
)

func Test_validateAllowedCharacters(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "только латиница",
			args: args{
				value: "abcd",
			},
			wantErr: false,
		},
		{
			name: "только цифры",
			args: args{
				value: "123",
			},
			wantErr: false,
		},
		{
			name: "только дефисы",
			args: args{
				value: "-----",
			},
			wantErr: false,
		},
		{
			name: "латиница и цифры",
			args: args{
				value: "abcd123",
			},
			wantErr: false,
		},
		{
			name: "латиница и дефисы",
			args: args{
				value: "abcd---",
			},
			wantErr: false,
		},
		{
			name: "цифры и дефисы",
			args: args{
				value: "123---",
			},
			wantErr: false,
		},
		{
			name: "не латиница, не цифры или не дефисы",
			args: args{
				value: "привет!;*",
			},
			wantErr: true,
		},
		{
			name: "валидные и невалидные символы",
			args: args{
				value: "прdb1",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateAllowedCharacters(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("validateAllowedCharacters() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateRequiredValue(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "пустое значение возвращает ошибку",
			args: args{
				"",
			},
			wantErr: true,
		},
		{
			name: "непустое значение не возвращает ошибку",
			args: args{
				value: "a",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateRequiredValue(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("validateRequiredValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateValue(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "пустое значение",
			args: args{
				value: "",
			},
			wantErr: true,
		},
		{
			name: "превышает количество символов",
			args: args{
				value: strings.Repeat("a", MaxValueCharsCount+1),
			},
			wantErr: true,
		},
		{
			name: "меньше нужного количества символов",
			args: args{
				value: strings.Repeat("a", MinValueCharsCount-1),
			},
			wantErr: true,
		},
		{
			name: "содержит запрещенные символы",
			args: args{
				value: strings.Repeat(".", 5),
			},
			wantErr: true,
		},
		{
			name: "корректное значение",
			args: args{
				value: "abcd123-",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateValue(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("validateValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateValueCharsLimit(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "меньше минимального количества символов символов",
			args: args{
				value: strings.Repeat("a", MinValueCharsCount-1),
			},
			wantErr: true,
		},
		{
			name: "пустое значение",
			args: args{
				value: "",
			},
			wantErr: true,
		},
		{
			name: "больше максимального количества символов",
			args: args{
				value: strings.Repeat("a", MaxValueCharsCount+1),
			},
			wantErr: true,
		},
		{
			name: "корректное количество символов",
			args: args{
				value: strings.Repeat("a", MinValueCharsCount+1),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateValueCharsLimit(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("validateValueCharsLimit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
