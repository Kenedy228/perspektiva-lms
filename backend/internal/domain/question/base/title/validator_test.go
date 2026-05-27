package title

import (
	"strings"
	"testing"
)

func Test_validateRequired(t *testing.T) {
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
			name: "значение из пробелов",
			args: args{
				value: strings.Repeat(" ", ValueCharsLimit),
			},
			wantErr: false,
		},
		{
			name: "корректное значение",
			args: args{
				value: strings.Repeat("a", ValueCharsLimit),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateRequired(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("validateRequired() error = %v, wantErr %v", err, tt.wantErr)
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
			name: "пустая строка",
			args: args{
				value: "",
			},
			wantErr: true,
		},
		{
			name: "строка превышает количество символов",
			args: args{
				value: strings.Repeat("a", ValueCharsLimit+1),
			},
			wantErr: true,
		},
		{
			name: "строка не превышает количество символов",
			args: args{
				value: strings.Repeat("a", ValueCharsLimit),
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
			name: "не превышает лимит по символам",
			args: args{
				value: strings.Repeat("a", ValueCharsLimit),
			},
			wantErr: false,
		},
		{
			name: "превышает лимит по символам",
			args: args{
				value: strings.Repeat("a", ValueCharsLimit+1),
			},
			wantErr: true,
		},
		{
			name: "пустая строка",
			args: args{
				value: "",
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
