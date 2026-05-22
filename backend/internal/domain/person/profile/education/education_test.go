package education

import (
	"reflect"
	"strings"
	"testing"
)

func TestEducation_IsZero(t *testing.T) {
	type fields struct {
		value string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "пустое значение",
			fields: fields{
				value: "",
			},
			want: true,
		},
		{
			name: "непустое значение",
			fields: fields{
				value: strings.Repeat("a", 20),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Education{
				value: tt.fields.value,
			}
			if got := e.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEducation_Value(t *testing.T) {
	type fields struct {
		value string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "возвращает значение как есть",
			fields: fields{
				value: strings.Repeat("a", 20),
			},
			want: strings.Repeat("a", 20),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Education{
				value: tt.fields.value,
			}
			if got := e.Value(); got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    Education
		wantErr bool
	}{
		{
			name: "значение не превышает лимита по символам",
			args: args{
				value: strings.Repeat("a", ValueCharsLimit),
			},
			want: Education{
				value: strings.Repeat("a", ValueCharsLimit),
			},
			wantErr: false,
		},
		{
			name: "значение превышает лимит по символам",
			args: args{
				value: strings.Repeat("a", ValueCharsLimit+1),
			},
			want: Education{
				value: "",
			},
			wantErr: true,
		},
		{
			name: "пустое значение",
			args: args{
				value: "",
			},
			want: Education{
				value: "",
			},
			wantErr: false,
		},
		{
			name: "значение из пробелов",
			args: args{
				value: strings.Repeat(" ", 20),
			},
			want: Education{
				value: "",
			},
			wantErr: false,
		},
		{
			name: "непустое значение",
			args: args{
				value: strings.Repeat("a", 20),
			},
			want: Education{
				value: strings.Repeat("a", 20),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.value)
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
