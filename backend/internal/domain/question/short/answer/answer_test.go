package answer

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnswer_Clone(t *testing.T) {
	type fields struct {
		value string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "возвращает значение как есть",
			fields: fields{
				value: "value",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Answer{
				value: tt.fields.value,
			}
			got := a.Clone()

			assert.NotNil(t, got)
		})
	}
}

func TestAnswer_IsEmpty(t *testing.T) {
	type fields struct {
		value string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "пустой вариант ответа",
			fields: fields{
				value: "",
			},
			want: true,
		},
		{
			name: "непустой вариант ответа",
			fields: fields{
				value: "value",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Answer{
				value: tt.fields.value,
			}
			if got := a.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnswer_Value(t *testing.T) {
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
			a := Answer{
				value: tt.fields.value,
			}
			if got := a.Value(); got != tt.want {
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
		name string
		args args
		want Answer
	}{
		{
			name: "возвращает новый объект",
			args: args{
				value: "value",
			},
			want: Answer{
				value: "value",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
