package login

import "testing"

func Test_normalizeValue(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "удалить пробелы по краям",
			args: args{
				value: " login123 ",
			},
			want: "login123",
		},
		{
			name: "удалить все пробелы",
			args: args{
				value: "login 1 2 3",
			},
			want: "login123",
		},
		{
			name: "удалить пробелы по краям и внутри текста",
			args: args{
				value: " login 12 3 ",
			},
			want: "login123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeValue(tt.args.value); got != tt.want {
				t.Errorf("normalizeValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
