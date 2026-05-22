package snils

import (
	"strings"
	"testing"
)

func Test_normalize(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "только цифры",
			args: args{
				s: strings.Repeat("1", 30),
			},
			want: strings.Repeat("1", 30),
		},
		{
			name: "цифры и буквы",
			args: args{
				s: "123abcd",
			},
			want: "123abcd",
		},
		{
			name: "цифры и пробелы",
			args: args{
				s: "123 321   4555",
			},
			want: "1233214555",
		},
		{
			name: "цифры и дефисы",
			args: args{
				s: "123-321-455",
			},
			want: "123321455",
		},
		{
			name: "цифры, дефисы и пробелы",
			args: args{
				s: " 123-321 455--",
			},
			want: "123321455",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalize(tt.args.s); got != tt.want {
				t.Errorf("normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}
