package passhash

import (
	"reflect"
	"strings"
	"testing"
)

func TestHash_IsZero(t *testing.T) {
	type fields struct {
		hash string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "пустое значение хеша",
			fields: fields{
				hash: "",
			},
			want: true,
		},
		{
			name: "непустое значение хеша",
			fields: fields{
				hash: "hash",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := Hash{
				hash: tt.fields.hash,
			}
			if got := ph.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHash_Value(t *testing.T) {
	type fields struct {
		hash string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "значение возвращается как есть",
			fields: fields{
				hash: "hash",
			},
			want: "hash",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := Hash{
				hash: tt.fields.hash,
			}
			if got := ph.Value(); got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		hash string
	}
	tests := []struct {
		name    string
		args    args
		want    Hash
		wantErr bool
	}{
		{
			name: "пустой хеш",
			args: args{
				hash: "",
			},
			want: Hash{
				hash: "",
			},
			wantErr: true,
		},
		{
			name: "хеш превышает лимит по весу",
			args: args{
				hash: strings.Repeat("a", MaxHashBytes+1),
			},
			want: Hash{
				hash: "",
			},
			wantErr: true,
		},
		{
			name: "валидный хеш",
			args: args{
				hash: strings.Repeat("a", MaxHashBytes),
			},
			want: Hash{
				hash: strings.Repeat("a", MaxHashBytes),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.hash)
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
