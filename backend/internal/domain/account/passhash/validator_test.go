package passhash

import (
	"strings"
	"testing"
)

func Test_validateHash(t *testing.T) {
	type args struct {
		hash string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "пустой хеш",
			args: args{
				hash: "",
			},
			wantErr: true,
		},
		{
			name: "хеш из пробелов",
			args: args{
				hash: strings.Repeat(" ", 10),
			},
			wantErr: false,
		},
		{
			name: "непустой хеш",
			args: args{
				hash: "hash",
			},
			wantErr: false,
		},
		{
			name: "хеш весит больше лимита",
			args: args{
				hash: strings.Repeat("a", MaxHashBytes+1),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateHash(tt.args.hash); (err != nil) != tt.wantErr {
				t.Errorf("validateHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateHashSize(t *testing.T) {
	type args struct {
		hash string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "вес хеша превышает максимальный лимит по памяти",
			args: args{
				hash: strings.Repeat("a", MaxHashBytes+1),
			},
			wantErr: true,
		},
		{
			name: "вес хеша не превышает максимальный лимит по памяти",
			args: args{
				hash: strings.Repeat("a", MaxHashBytes-1),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateHashSize(tt.args.hash); (err != nil) != tt.wantErr {
				t.Errorf("validateHashSize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateRequiredHash(t *testing.T) {
	type args struct {
		hash string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "пустой хеш",
			args: args{
				hash: "",
			},
			wantErr: true,
		},
		{
			name: "непустой хеш",
			args: args{
				hash: "hash",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateRequiredHash(tt.args.hash); (err != nil) != tt.wantErr {
				t.Errorf("validateRequiredHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
