package name

import (
	"reflect"
	"strings"
	"testing"
)

func TestName_IsZero(t *testing.T) {
	type fields struct {
		value string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "пустое имя",
			fields: fields{
				value: "",
			},
			want: true,
		},
		{
			name: "непустое имя",
			fields: fields{
				value: "наименование",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Name{
				value: tt.fields.value,
			}
			if got := n.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestName_Value(t *testing.T) {
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
				value: "наименование",
			},
			want: "наименование",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Name{
				value: tt.fields.value,
			}
			if got := n.Value(); got != tt.want {
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
		want    Name
		wantErr bool
	}{
		{
			name: "пустое значение",
			args: args{
				value: "",
			},
			want: Name{
				value: "",
			},
			wantErr: true,
		},
		{
			name: "превышает лимит символов",
			args: args{
				value: strings.Repeat("a", ValueCharsLimit+1),
			},
			want: Name{
				value: "",
			},
			wantErr: true,
		},
		{
			name: "не превышает лимит",
			args: args{
				value: strings.Repeat("a", ValueCharsLimit),
			},
			want: Name{
				value: strings.Repeat("a", ValueCharsLimit),
			},
			wantErr: false,
		},
		{
			name: "сохраняет нормализованное значение",
			args: args{
				value: "   ООО     Ромашка    ",
			},
			want: Name{
				value: "ООО Ромашка",
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
