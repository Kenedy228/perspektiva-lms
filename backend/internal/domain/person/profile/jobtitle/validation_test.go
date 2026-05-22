package jobtitle

import (
	"strings"
	"testing"
)

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
			name: "значение не превышает лимита по символам",
			args: args{
				value: strings.Repeat("a", ValueCharsLimit),
			},
			wantErr: false,
		},
		{
			name: "значение превышает лимит по символам",
			args: args{
				value: strings.Repeat("a", ValueCharsLimit+1),
			},
			wantErr: true,
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
			name: "значение не превышает лимита по символам",
			args: args{
				value: strings.Repeat("a", ValueCharsLimit),
			},
			wantErr: false,
		},
		{
			name: "значение превышает лимит по символам",
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
			wantErr: true,
		},
		{
			name: "значение из пробелов",
			args: args{
				value: strings.Repeat(" ", 20),
			},
			wantErr: false,
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
			if err := validateValue(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("validateValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
