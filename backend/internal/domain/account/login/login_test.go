package login

import (
	"reflect"
	"strings"
	"testing"
)

func TestLogin_IsZero(t *testing.T) {
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
				value: "admin123",
			},
			want: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Login{
				value: tt.fields.value,
			}
			if got := l.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogin_Value(t *testing.T) {
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
				value: "admin123",
			},
			want: "admin123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Login{
				value: tt.fields.value,
			}
			if got := l.Value(); got != tt.want {
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
		want    Login
		wantErr bool
	}{
		{
			name: "валидный логин",
			args: args{
				value: "admin123-",
			},
			want: Login{
				value: "admin123-",
			},
			wantErr: false,
		},
		{
			name: "удаление пробелов по краям",
			args: args{
				value: " admin123- ",
			},
			want: Login{
				value: "admin123-",
			},
		},
		{
			name: "удаление пробелов внутри строки и по краям",
			args: args{
				value: " admin 123 - ",
			},
			want: Login{
				value: "admin123-",
			},
		},
		{
			name: "пустое значение",
			args: args{
				value: "",
			},
			want: Login{
				value: "",
			},
			wantErr: true,
		},
		{
			name: "символов меньше минимального количества",
			args: args{
				value: strings.Repeat("a", MinValueCharsCount-1),
			},
			want: Login{
				value: "",
			},
			wantErr: true,
		},
		{
			name: "символов больше максимального количества",
			args: args{
				value: strings.Repeat("a", MaxValueCharsCount+1),
			},
			want: Login{
				value: "",
			},
			wantErr: true,
		},
		{
			name: "содержит запрещенные символы",
			args: args{
				value: "привет",
			},
			want: Login{
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
