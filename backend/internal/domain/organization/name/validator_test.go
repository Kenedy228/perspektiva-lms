package name

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
			name: "непустое значение",
			args: args{
				value: strings.Repeat("a", 20),
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
			name: "пустое значение",
			args: args{
				value: "",
			},
			wantErr: true,
		},
		{
			name: "превышает лимит символов",
			args: args{
				value: strings.Repeat("a", ValueCharsLimit+1),
			},
			wantErr: true,
		},
		{
			name: "не превышает лимит",
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
			name: "превышает лимит символов",
			args: args{
				value: strings.Repeat("a", ValueCharsLimit+1),
			},
			wantErr: true,
		},

		{
			name: "пустое значение",
			args: args{
				value: "",
			},
			wantErr: false,
		},
		{
			name: "не превышает лимит",
			args: args{
				value: strings.Repeat("a", ValueCharsLimit),
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
