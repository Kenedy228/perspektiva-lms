package snils

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
		want    SNILS
		wantErr bool
	}{
		{
			name: "корректный отформатированный СНИЛС",
			args: args{
				value: "112-233-445 95",
			},
			want: SNILS{
				value: "11223344595",
			},
			wantErr: false,
		},
		{
			name: "корректный неотформатированный СНИЛС",
			args: args{
				value: "11223344595",
			},
			want: SNILS{
				value: "11223344595",
			},
			wantErr: false,
		},
		{
			name: "некорректный СНИЛС",
			args: args{
				value: "11223344594",
			},
			want: SNILS{
				value: "",
			},
			wantErr: true,
		},
		{
			name: "СНИЛС с недопустимыми символами",
			args: args{
				value: "abc123",
			},
			want: SNILS{
				value: "",
			},
			wantErr: true,
		},
		{
			name: "СНИЛС по длине меньше требуемой",
			args: args{
				value: "123",
			},
			want: SNILS{
				value: "",
			},
			wantErr: true,
		},
		{
			name: "СНИЛС по длине больше требуемой",
			args: args{
				value: strings.Repeat("1", 1e5),
			},
			want: SNILS{
				value: "",
			},
			wantErr: true,
		},
		{
			name: "запрещенное значение СНИЛС",
			args: args{
				value: "00000000000",
			},
			want: SNILS{
				value: "",
			},
			wantErr: true,
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

func TestSNILS_IsZero(t *testing.T) {
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
				value: "11111111111",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SNILS{
				value: tt.fields.value,
			}
			if got := s.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSNILS_Value(t *testing.T) {
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
				value: "1234567899",
			},
			want: "1234567899",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SNILS{
				value: tt.fields.value,
			}
			if got := s.Value(); got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
