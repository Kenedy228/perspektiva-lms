package option

import (
	"reflect"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    Option
		wantErr bool
	}{
		{
			name: "пустое значение",
			args: args{
				value: "",
			},
			want:    Option{},
			wantErr: true,
		},
		{
			name: "значение из пробелов",
			args: args{
				value: strings.Repeat(" ", ValueCharsLimit),
			},
			want:    Option{},
			wantErr: true,
		},
		{
			name: "значение превышает лимит по символам",
			args: args{
				value: strings.Repeat("a", ValueCharsLimit+1),
			},
			want:    Option{},
			wantErr: true,
		},
		{
			name: "значение не превышает лимит по символам",
			args: args{
				value: strings.Repeat("a", ValueCharsLimit),
			},
			want: Option{
				value: strings.Repeat("a", ValueCharsLimit),
			},
			wantErr: false,
		},
		{
			name: "значение нормализуется",
			args: args{
				value: "    val     ue   ",
			},
			want: Option{
				value: "val ue",
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

func TestOption_IsZero(t *testing.T) {
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
				value: "value",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := Option{
				value: tt.fields.value,
			}
			if got := o.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOption_Value(t *testing.T) {
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
				value: "value",
			},
			want: "value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := Option{
				value: tt.fields.value,
			}
			if got := o.Value(); got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
