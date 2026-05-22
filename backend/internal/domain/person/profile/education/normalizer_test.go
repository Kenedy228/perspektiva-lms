package education

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
			name: "удаляет незначащие пробелы",
			args: args{
				value: " высшее ",
			},
			want: "высшее",
		},
		{
			name: "удаляет дубликаты пробелов внутри",
			args: args{
				value: " высшее  юридическое  \tобразование",
			},
			want: "высшее юридическое образование",
		},
		{
			name: "не меняет валидную строку",
			args: args{
				value: "образование",
			},
			want: "образование",
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
